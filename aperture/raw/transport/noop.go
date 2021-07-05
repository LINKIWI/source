package transport

// Noop is a Transport interface-compliant implementation that noops on all operations.
// It is intended to be used as a placeholder implementation when no real transport is available.
type Noop struct{}

// NewNoop creates a new Noop transport.
func NewNoop() Transport {
	return &Noop{}
}

// Send noops.
func (t *Noop) Send(data []byte) (int, error) {
	return 0, nil
}

// Close noops.
func (t *Noop) Close() error {
	return nil
}
