package transport

import (
	"sync"

	"lib.kevinlin.info/aperture/internal/errors"
)

// Buffered is an abstraction over a Transport that internally buffers the transmission of
// individual payloads with a configurable batch size.
type Buffered struct {
	buf   chan []byte
	mutex sync.Mutex
	Transport
}

// NewBuffered creates a new buffered transport over an existing Transport backend.
func NewBuffered(backendFactory Factory, size int) (Transport, error) {
	if size <= 0 {
		return nil, errors.New(
			"transport",
			"buffered transport batch size must be greater than 0",
		)
	}

	tport, err := backendFactory()
	if err != nil {
		return nil, err
	}

	return &Buffered{
		Transport: tport,
		buf:       make(chan []byte, size),
	}, nil
}

// Send asynchronously queues the payload for transmission. If the queue size has reached the batch
// size, all queued payloads are flushed.
func (t *Buffered) Send(data []byte) (int, error) {
	t.buf <- data

	if len(t.buf) == cap(t.buf) {
		t.flush()
	}

	return len(data), nil
}

// Close flushes all remaining queued payloads and closes the backing transport.
func (t *Buffered) Close() error {
	t.flush()

	return t.Transport.Close()
}

// flush sends all queued payloads, ignoring all errors.
func (t *Buffered) flush() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	elements := len(t.buf)

	for i := 0; i < elements; i++ {
		t.Transport.Send(<-t.buf)
	}
}
