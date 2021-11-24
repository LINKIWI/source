package server

import (
	"net"
	"time"
)

// TimeoutConn is a net.Conn that wraps a connection with R/W socket timeouts.
type TimeoutConn struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	net.Conn
}

// Read adds a socket read deadline and defers to the underlying connection.
func (c *TimeoutConn) Read(b []byte) (n int, err error) {
	if c.ReadTimeout > 0 {
		if err := c.Conn.SetReadDeadline(time.Now().Add(c.ReadTimeout)); err != nil {
			return 0, err
		}
	}

	return c.Conn.Read(b)
}

// Write adds a socket write deadline and defers to the underlying connection.
func (c *TimeoutConn) Write(b []byte) (n int, err error) {
	if c.WriteTimeout > 0 {
		if err := c.Conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout)); err != nil {
			return 0, err
		}
	}

	return c.Conn.Write(b)
}
