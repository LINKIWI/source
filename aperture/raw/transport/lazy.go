package transport

import (
	"sync"
)

// Lazy is an abstraction over a Transport that initializes the underlying transport only when the
// client attempts to make its first transmission, rather than at construction time.
type Lazy struct {
	backendFactory Factory
	mutex          sync.Mutex
	Transport
}

// NewLazy creates a Transport that lazily initializes the underlying network transport client.
func NewLazy(backendFactory Factory) (Transport, error) {
	return &Lazy{backendFactory: backendFactory}, nil
}

// Send idempotently initializes the underlying transport, followed by using the transport to send
// the requested data.
func (t *Lazy) Send(data []byte) (int, error) {
	if err := t.init(); err != nil {
		return 0, err
	}

	return t.Transport.Send(data)
}

// Close defers to the underlying transport implementation, if initialized, to close the transport.
func (t *Lazy) Close() error {
	if t.Transport == nil {
		return nil
	}

	return t.Transport.Close()
}

// init statefully but idempotently initializes the underlying transport. This method noops if the
// transport is already initialized; otherwise, it initializes and sets the transport for future
// reuse on subsequent transmissions. The transport is not set in case of initialization errors.
func (t *Lazy) init() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Transport is already initialized; abort immediately.
	if t.Transport != nil {
		return nil
	}

	tport, err := t.backendFactory()
	if err != nil {
		return err
	}

	t.Transport = tport

	return nil
}
