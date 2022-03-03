package config

import (
	"net"
	"time"
)

// timeoutConn is an abstraction over a net.Conn that sets timeouts for every socket operation.
type timeoutConn struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	net.Conn
}

// newTimeoutConn creates a new timeoutConn with the specified timeout parameters.
func newTimeoutConn(readTimeout time.Duration, writeTimeout time.Duration, conn net.Conn) net.Conn {
	return &timeoutConn{
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		Conn:         conn,
	}
}

// Read sets the socket read deadline appropriately and defers to the underlying net.Conn.
func (c *timeoutConn) Read(b []byte) (int, error) {
	if c.readTimeout > 0 {
		if err := c.Conn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			return 0, err
		}
	}

	return c.Conn.Read(b)
}

// Write sets the socket write deadline appropriately and defers to the underlying net.Conn.
func (c *timeoutConn) Write(b []byte) (int, error) {
	if c.writeTimeout > 0 {
		if err := c.Conn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
			return 0, err
		}
	}

	return c.Conn.Write(b)
}
