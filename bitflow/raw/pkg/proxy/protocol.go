package proxy

import (
	"fmt"
	"net"
)

const (
	// Supported HAProxy protocols.
	protoTCP4    = "TCP4"
	protoTCP6    = "TCP6"
	protoUnknown = "UNKNOWN"
)

// protocolHeader is an implementation of the HAProxy PROXY protocol header, derived from a source
// and destination net.Conn connection.
// Reference: https://www.haproxy.org/download/1.8/doc/proxy-protocol.txt
type protocolHeader struct {
	src net.Conn
	dst net.Conn
}

// String returns the full PROXY header byte sequence per the protocol specification.
func (h *protocolHeader) String() string {
	var proto string

	clientNet, clientIP, clientPort, err := parseAddress(h.src.RemoteAddr())
	if err != nil {
		proto = protoUnknown
		clientIP = "0.0.0.0"
		clientPort = 0
	}

	_, targetIP, targetPort, err := parseAddress(h.dst.RemoteAddr())
	if err != nil {
		proto = protoUnknown
		targetIP = "0.0.0.0"
		targetPort = 0
	}

	switch clientNet {
	case "tcp4":
		proto = protoTCP4
	case "tcp6":
		proto = protoTCP6
	default:
		proto = protoUnknown
	}

	return fmt.Sprintf(
		"PROXY %s %s %s %d %d\r\n",
		proto,
		clientIP,
		targetIP,
		clientPort,
		targetPort,
	)
}
