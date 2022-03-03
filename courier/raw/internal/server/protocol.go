package server

import (
	"bytes"
	"net"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"

	"courier/internal/metrics"
	"courier/internal/util"
)

const (
	metricProxyHeaderParseSuccess = "server.connection.proxy_header.success"
	metricProxyHeaderParseError   = "server.connection.proxy_header.error"
)

// String sequences known to version 1 of the HAProxy PROXY protocol.
// Reference: https://www.haproxy.org/download/1.8/doc/proxy-protocol.txt
var (
	headerPreamble     = []byte("PROXY")
	headerProtoTCP4    = []byte("TCP4")
	headerProtoTCP6    = []byte("TCP6")
	headerProtoUnknown = []byte("UNKNOWN")
	headerDelimiter    = []byte(" ")
	headerTrailer      = []byte("\r\n")
)

// proxyListener wraps a net.Listener and automatically provides connections that are PROXY protocol
// aware.
type proxyListener struct {
	name          string
	headerTimeout time.Duration
	net.Listener
}

// proxyConn wraps a net.Conn to parse and consume PROXY protocol headers before serving any reads.
// It additionally overrides RemoteAddr to transparently provide the proxy-supplied source address.
type proxyConn struct {
	proto           string
	sourceIP        net.IP
	sourcePort      int
	destinationIP   net.IP
	destinationPort int

	name          string
	headerTimeout time.Duration
	headerParser  sync.Once
	bufferedConn  *bufferedConn

	net.Conn
}

// newProxyListener creates a PROXY protocol-aware net.Listener.
func newProxyListener(name string, headerTimeout time.Duration, ln net.Listener) net.Listener {
	return &proxyListener{
		name:          name,
		headerTimeout: headerTimeout,
		Listener:      ln,
	}
}

// Accept defers to the underlying net.Listener and returns a PROXY protocol-aware wrapped net.Conn
// on success.
func (l *proxyListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return conn, &util.Error{
			Namespace:    "server",
			Message:      "failed to accept on listener",
			StackedError: err,
		}
	}

	return &proxyConn{
		name:          l.name,
		headerTimeout: l.headerTimeout,
		bufferedConn:  newBufferedConn(conn),
		Conn:          conn,
	}, nil
}

// Read attempts to read the PROXY protocol header if it hasn't already been read, followed by
// deferring to the underlying connection for the requested read.
func (c *proxyConn) Read(b []byte) (n int, err error) {
	c.headerParser.Do(c.initProxyMetadata)

	return c.bufferedConn.Read(b)
}

// RemoteAddr attempts to read the PROXY protocol header if it hasn't already been read, and returns
// a net.TCPAddr with the proxy-supplied source address.
func (c *proxyConn) RemoteAddr() net.Addr {
	c.headerParser.Do(c.initProxyMetadata)

	return &net.TCPAddr{
		IP:   c.sourceIP,
		Port: c.sourcePort,
	}
}

// initProxyMetadata initializes proxy metadata by reading the PROXY protocol header. Due to strict
// usage enforcement of the proxy header, this method will also close the client connection if the
// header cannot be parsed or is otherwise invalid.
func (c *proxyConn) initProxyMetadata() {
	tags := map[string]interface{}{"name": c.name}

	if err := c.parseProxyHeader(); err != nil {
		metrics.Client.Incr(metricProxyHeaderParseError, tags)
		zap.L().Warn("error parsing proxy protocol header", zap.Error(err))

		c.Close()
		return
	}

	metrics.Client.Incr(metricProxyHeaderParseSuccess, tags)
}

// parseProxyHeader statefully reads from the connection and attempts to parse the PROXY protocol
// header. It returns a non-nil error on any parse failures.
func (c *proxyConn) parseProxyHeader() error {
	if c.headerTimeout > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(c.headerTimeout))
		defer c.Conn.SetReadDeadline(time.Time{})
	}

	bufPreamble, err := c.bufferedConn.ReadUntil(headerDelimiter, len(headerPreamble)+1)
	if err != nil {
		return &util.Error{
			Namespace:    "server",
			Message:      "error reading proxy header preamble",
			StackedError: err,
		}
	} else if !bytes.Equal(bufPreamble, headerPreamble) {
		return &util.Error{
			Namespace: "server",
			Message:   "proxy header preamble is invalid",
			Tags:      map[string]interface{}{"preamble": string(bufPreamble)},
		}
	}

	bufProto, err := c.bufferedConn.ReadUntil(headerDelimiter, len(headerProtoUnknown)+1)
	if err != nil {
		return err
	}

	switch string(bufProto) {
	case string(headerProtoTCP4):
		c.proto = "TCP4"
	case string(headerProtoTCP6):
		c.proto = "TCP6"
	case string(headerProtoUnknown):
		c.proto = "UNKNOWN"
	default:
		return &util.Error{
			Namespace: "server",
			Message:   "unknown proxy protocol",
			Tags: map[string]interface{}{
				"protocol": string(bufProto),
			},
		}
	}

	bufSrcIP, err := c.bufferedConn.ReadUntil(headerDelimiter, 40)
	if err != nil {
		return &util.Error{
			Namespace:    "server",
			Message:      "error reading proxy header source IP",
			StackedError: err,
		}
	} else if c.sourceIP = net.ParseIP(string(bufSrcIP)); c.sourceIP == nil {
		return &util.Error{
			Namespace: "server",
			Message:   "error parsing source IP",
			Tags:      map[string]interface{}{"ip": string(bufSrcIP)},
		}
	}

	bufDstIP, err := c.bufferedConn.ReadUntil(headerDelimiter, 40)
	if err != nil {
		return &util.Error{
			Namespace:    "server",
			Message:      "error reading proxy header destination IP",
			StackedError: err,
		}
	} else if c.destinationIP = net.ParseIP(string(bufDstIP)); c.destinationIP == nil {
		return &util.Error{
			Namespace: "server",
			Message:   "error parsing destination IP",
			Tags:      map[string]interface{}{"ip": string(bufDstIP)},
		}
	}

	bufSrcPort, err := c.bufferedConn.ReadUntil(headerDelimiter, 6)
	if err != nil {
		return &util.Error{
			Namespace:    "server",
			Message:      "error reading proxy header source port",
			StackedError: err,
		}
	} else if c.sourcePort, err = strconv.Atoi(string(bufSrcPort)); err != nil {
		return &util.Error{
			Namespace:    "server",
			Message:      "error parsing source port",
			Tags:         map[string]interface{}{"port": string(bufSrcPort)},
			StackedError: err,
		}
	} else if c.sourcePort < 0 || c.sourcePort > 65535 {
		return &util.Error{
			Namespace: "server",
			Message:   "invalid source port",
			Tags:      map[string]interface{}{"port": c.sourcePort},
		}
	}

	bufDstPort, err := c.bufferedConn.ReadUntil(headerTrailer, 7)
	if err != nil {
		return &util.Error{
			Namespace:    "server",
			Message:      "error reading proxy header destination port",
			StackedError: err,
		}
	} else if c.destinationPort, err = strconv.Atoi(string(bufDstPort)); err != nil {
		return &util.Error{
			Namespace:    "server",
			Message:      "error parsing destination port",
			Tags:         map[string]interface{}{"port": string(bufDstPort)},
			StackedError: err,
		}
	} else if c.destinationPort < 0 || c.destinationPort > 65535 {
		return &util.Error{
			Namespace: "server",
			Message:   "invalid destination port",
			Tags:      map[string]interface{}{"port": c.destinationPort},
		}
	}

	return nil
}
