package transport

// Factory is a type alias for a lazy provider of a Transport backend.
type Factory func() (Transport, error)

// Transport provides the interface for any underlying unidirectional network transport for shipping
// arbitrary data.
type Transport interface {
	// Send a payload with the transport, returning the number of bytes written.
	Send(data []byte) (int, error)

	// Close the transport.
	Close() error
}
