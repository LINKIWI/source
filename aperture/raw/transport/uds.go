package transport

import (
	"net"

	"lib.kevinlin.info/aperture/internal/errors"
)

// UDS is a Transport for transmitting datagrams to a Unix datagram socket.
type UDS struct {
	conn *net.UnixConn
}

// NewUDS creates a Unix datagram socket transport for a Unix socket at the target disk path.
func NewUDS(path string) (Transport, error) {
	conn, err := net.DialUnix("unixgram", nil, &net.UnixAddr{Name: path, Net: "unixgram"})
	if err != nil {
		return nil, errors.Stack(
			"transport",
			"failed to connect to Unix socket",
			err,
			errors.Tag{Key: "socket", Value: path},
		)
	}

	return &UDS{conn}, nil
}

// Send transmits the payload over the Unix datagram socket.
func (t *UDS) Send(data []byte) (int, error) {
	return t.conn.Write(data)
}

// Close closes the underlying Unix datagram socket.
func (t *UDS) Close() error {
	return t.conn.Close()
}
