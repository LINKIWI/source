package transport

import (
	"net"
	"time"

	"golang.org/x/net/proxy"

	"lib.kevinlin.info/aperture/internal/errors"
)

// TCP is a Transport for transmitting bytes to a TCP address.
type TCP struct {
	conn net.Conn
	cfg  TCPSocketConfig
}

// TCPSocketConfig provides configuration parameters for the underlying TCP socket.
type TCPSocketConfig struct {
	ConnectTimeout    time.Duration
	SendTimeout       time.Duration
	KeepAliveInterval time.Duration
}

// NewDirectTCP establishes a TCP transport client against the target address directly.
func NewDirectTCP(addr string, cfg TCPSocketConfig) (Transport, error) {
	dialer := &net.Dialer{
		Timeout:   cfg.ConnectTimeout,
		KeepAlive: cfg.KeepAliveInterval,
	}

	return newTCP(addr, dialer, cfg)
}

// NewProxiedTCP establishes a TCP transport client against the target address via an intermediate
// SOCKS5 proxy server.
func NewProxiedTCP(addr string, proxyNet string, proxyAddr string, cfg TCPSocketConfig) (Transport, error) {
	forward := &net.Dialer{
		Timeout:   cfg.ConnectTimeout,
		KeepAlive: cfg.KeepAliveInterval,
	}

	dialer, err := proxy.SOCKS5(proxyNet, proxyAddr, nil, forward)
	if err != nil {
		return nil, err
	}

	return newTCP(addr, dialer, cfg)
}

// newTCP establishes a TCP transport client using a proxy.Dialer and set of socket options.
func newTCP(addr string, dialer proxy.Dialer, cfg TCPSocketConfig) (Transport, error) {
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Stack(
			"transport",
			"failed to connect to TCP socket",
			err,
			errors.Tag{Key: "addr", Value: addr},
		)
	}

	return &TCP{
		conn: conn,
		cfg:  cfg,
	}, nil
}

// Send initiates a reliable transmission of an arbitrary payload, respecting timeout parameters.
func (t *TCP) Send(data []byte) (int, error) {
	if t.cfg.SendTimeout > 0 {
		t.conn.SetWriteDeadline(time.Now().Add(t.cfg.SendTimeout))
	}

	return t.conn.Write(data)
}

// Close closes the underlying TCP socket.
func (t *TCP) Close() error {
	return t.conn.Close()
}
