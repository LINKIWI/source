package server

import (
	"bufio"
	"bytes"
	"net"

	"courier/internal/util"
)

// bufferedConn is an abstraction over a net.Conn that internally buffers reads from the connection.
type bufferedConn struct {
	reader *bufio.Reader
	net.Conn
}

// newBufferedConn wraps an existing net.Conn with internal buffering.
func newBufferedConn(conn net.Conn) *bufferedConn {
	return &bufferedConn{
		reader: bufio.NewReader(conn),
		Conn:   conn,
	}
}

// Peek allows peeking into the first n bytes of the connection without actually consuming them.
func (c *bufferedConn) Peek(n int) ([]byte, error) {
	return c.reader.Peek(n)
}

// Read reads from the buffered connection, consuming those bytes in the process.
func (c *bufferedConn) Read(b []byte) (int, error) {
	return c.reader.Read(b)
}

// ReadUntil is a utility for reading from the connection until a particular byte sequence is
// encountered. The read is limited to at most maxSeek bytes; failure to find the termination
// sequence within the first maxSeek bytes of the reader will result in an error. Only the size of
// the payload up to and including the termination sequence is actually read (i.e. consumed) from
// the underlying buffered connection. The returned slice omits the termination byte sequence.
func (c *bufferedConn) ReadUntil(sequence []byte, maxSeek int) ([]byte, error) {
	if maxSeek < len(sequence) {
		return nil, &util.Error{
			Namespace: "server",
			Message:   "sequence is shorter than maximum allowed seek",
		}
	}

	buf := make([]byte, 0, maxSeek)
	scan := make([]byte, 1)

	for i := 0; i < maxSeek; i++ {
		_, err := c.Read(scan)
		if err != nil {
			return nil, &util.Error{
				Namespace:    "server",
				Message:      "failed to read from buffered connection",
				StackedError: err,
			}
		}

		buf = append(buf, scan...)

		if i >= len(sequence)-1 && bytes.Equal(buf[len(buf)-len(sequence):], sequence) {
			return buf[:len(buf)-len(sequence)], nil
		}
	}

	return nil, &util.Error{
		Namespace: "server",
		Message:   "exhausted maximum allowed seek without finding sequence",
		Tags: map[string]interface{}{
			"sequence": sequence,
			"seek":     maxSeek,
		},
	}
}
