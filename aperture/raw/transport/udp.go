package transport

import (
	"net"

	"lib.kevinlin.info/aperture/internal/errors"
)

// UDP is a Transport for transmitting datagrams to a UDP address.
type UDP struct {
	conn *net.UDPConn
}

// NewUDP creates a UDP transport for a target address.
func NewUDP(addr string) (Transport, error) {
	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, errors.Stack(
			"transport",
			"UDP address resolution failed",
			err,
			errors.Tag{Key: "addr", Value: addr},
		)
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, errors.Stack(
			"transport",
			"failed to create UDP socket",
			err,
			errors.Tag{Key: "addr", Value: addr},
		)
	}

	return &UDP{conn}, nil
}

// Send initiates a fire-and-forget transmission of an arbitrary payload.
func (t *UDP) Send(data []byte) (int, error) {
	return t.conn.Write(data)
}

// Close closes the underlying UDP socket.
func (t *UDP) Close() error {
	return t.conn.Close()
}
