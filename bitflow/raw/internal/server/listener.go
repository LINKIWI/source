package server

import (
	"fmt"
	"net"

	"bitflow/internal/config"
)

// restrictedListener is a net.Listener that only accepts TCP connections from allowed remote IP
// addresses. Connections from disallowed source IPs are closed immediately after acceptance.
type restrictedListener struct {
	sources []config.CIDR
	net.Listener
}

// Accept defers to the underlying net.Listener and additionally verifies that the accepted
// connection originates from a permitted IP block. Disallowed connections are closed immediately,
// and Accept returns an error.
func (l *restrictedListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, fmt.Errorf("server: error accepting listener connection: err=%v", err)
	}

	addr, ok := conn.RemoteAddr().(*net.TCPAddr)
	if !ok {
		conn.Close()
		return nil, fmt.Errorf(
			"server: socket type not supported for restricted listener: type=%T",
			conn.RemoteAddr(),
		)
	}

	for _, source := range l.sources {
		if source.Contains(addr.IP) {
			return conn, nil
		}
	}

	conn.Close()
	return conn, fmt.Errorf(
		"server: rejected connection from unauthorized source: addr=%v sources=%v",
		conn.RemoteAddr(),
		l.sources,
	)
}
