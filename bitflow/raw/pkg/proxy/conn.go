package proxy

import (
	"net"
	"strconv"
	"time"
)

// timeoutConn is a net.Conn that wraps a connection with R/W socket timeouts.
type timeoutConn struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	net.Conn
}

// Read adds a socket read deadline and defers to the underlying connection.
func (c *timeoutConn) Read(b []byte) (n int, err error) {
	if c.readTimeout > 0 {
		if err := c.Conn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			return 0, err
		}
	}

	return c.Conn.Read(b)
}

// Write adds a socket write deadline and defers to the underlying connection.
func (c *timeoutConn) Write(b []byte) (n int, err error) {
	if c.writeTimeout > 0 {
		if err := c.Conn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
			return 0, err
		}
	}

	return c.Conn.Write(b)
}

// parseAddress is a convenience function for parsing a net.Addr into a network, host, and port.
// It additionally qualifies the network for TCP addresses to tcp4 or tcp6.
func parseAddress(addr net.Addr) (string, string, int, error) {
	network := addr.Network()

	if tcpAddr, ok := addr.(*net.TCPAddr); ok {
		switch {
		case tcpAddr.IP.To4() != nil:
			network = "tcp4"
		case tcpAddr.IP.To16() != nil:
			network = "tcp6"
		}
	}

	host, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		return network, "", 0, err
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return network, "", 0, err
	}

	return network, host, portNum, nil
}
