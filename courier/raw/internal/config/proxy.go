package config

import (
	"bufio"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"

	"courier/internal/util"
)

// init registers all known forward proxy dialer implementations.
func init() {
	proxy.RegisterDialerType("http", newHTTPProxy)
	proxy.RegisterDialerType("http+unix", newHTTPProxy)
	proxy.RegisterDialerType("socks5+unix", newSocks5Proxy)
	proxy.RegisterDialerType("socks5h+unix", newSocks5Proxy)
}

// httpProxy is an HTTP CONNECT forward proxy.
type httpProxy struct {
	proxyNet  string
	proxyAddr string
	forward   proxy.Dialer
}

// newHTTPProxy creates a new httpProxy. It is suitable for use with proxy.RegisterDialerType to
// register a Dial implementation for an HTTP proxy server.
func newHTTPProxy(proxyURL *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
	var proxyNet string
	var proxyAddr string

	switch proxyURL.Scheme {
	case "http":
		proxyNet = "tcp"
		proxyAddr = proxyURL.Host
	case "unix":
		proxyNet = "unix"
		proxyAddr = proxyURL.Host + proxyURL.Path
	}

	return &httpProxy{
		proxyNet:  proxyNet,
		proxyAddr: proxyAddr,
		forward:   forward,
	}, nil
}

// Dial dials the proxy server and attempts to negotiate a connection to the target upstream.
func (p *httpProxy) Dial(network string, addr string) (net.Conn, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
	default:
		return nil, &util.Error{
			Namespace: "config",
			Message:   "upstream address network not implemented",
			Tags:      map[string]interface{}{"network": network},
		}
	}

	// Use the forward dialer to dial the HTTP proxy server
	conn, err := p.forward.Dial(p.proxyNet, p.proxyAddr)
	if err != nil {
		return nil, &util.Error{
			Namespace: "config",
			Message:   "failed to dial proxy host with forward dialer",
			Tags: map[string]interface{}{
				"network": p.proxyNet,
				"address": p.proxyAddr,
			},
			StackedError: err,
		}
	}

	// Ask the proxy server to CONNECT to the original target upstream address
	upstreamURL := &url.URL{Host: addr}
	proxyReq, err := http.NewRequest(http.MethodConnect, upstreamURL.String(), nil)
	if err != nil {
		conn.Close()
		return nil, &util.Error{
			Namespace:    "config",
			Message:      "failed to create proxy CONNECT request",
			Tags:         map[string]interface{}{"url": upstreamURL},
			StackedError: err,
		}
	}

	if err := proxyReq.Write(conn); err != nil {
		conn.Close()
		return nil, &util.Error{
			Namespace:    "config",
			Message:      "failed to write to proxy server connection",
			StackedError: err,
		}
	}

	proxyResp, err := http.ReadResponse(bufio.NewReader(conn), proxyReq)
	if err != nil {
		conn.Close()
		return nil, &util.Error{
			Namespace:    "config",
			Message:      "failed to read HTTP response from proxy server",
			StackedError: err,
		}
	}

	defer proxyResp.Body.Close()

	if proxyResp.StatusCode != http.StatusOK {
		conn.Close()
		return nil, &util.Error{
			Namespace: "config",
			Message:   "HTTP proxy server returned non-OK status for CONNECT",
			Tags:      map[string]interface{}{"status": proxyResp.StatusCode},
		}
	}

	// Pass the connection along to the consumer at this point: the proxy server has accepted
	// the CONNECT request and further reads and writes from the connection henceforth are
	// tunneled to the original target upstream server.
	return conn, nil
}

// socks5Proxy is a SOCKS5 forward proxy.
type socks5Proxy struct {
	proxyAddr string
	forward   proxy.Dialer
}

// newSocks5Proxy creates a new socks5Proxy. It is suitable for use with proxy.RegisterDialerType to
// register a Dial implementation for a SOCKS5 proxy server over a Unix domain socket.
// Note that the golang.org/x/net/proxy standard library already registers a scheme for "socks5" and
// "socks5h" over TCP by default.
func newSocks5Proxy(proxyURL *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
	return &socks5Proxy{
		proxyAddr: proxyURL.Path,
		forward:   forward,
	}, nil
}

// Dial dials the proxy server over a Unix domain socket and uses the SOCKS5 dialer to dial the
// target upstream address.
func (s *socks5Proxy) Dial(network string, addr string) (net.Conn, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
	default:
		return nil, &util.Error{
			Namespace: "config",
			Message:   "upstream address network not implemented",
			Tags:      map[string]interface{}{"network": network},
		}
	}

	dialer, err := proxy.SOCKS5("unix", s.proxyAddr, nil, s.forward)
	if err != nil {
		return nil, &util.Error{
			Namespace:    "config",
			Message:      "failed to create SOCKS5 proxy dialer",
			StackedError: err,
		}
	}

	conn, err := dialer.Dial(network, addr)
	if err != nil {
		return nil, &util.Error{
			Namespace: "config",
			Message:   "failed to dial target upstream host with SOCKS5 proxy dialer",
			Tags: map[string]interface{}{
				"network": network,
				"address": addr,
			},
			StackedError: err,
		}
	}

	return conn, nil
}
