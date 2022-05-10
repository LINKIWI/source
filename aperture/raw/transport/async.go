package transport

import (
	"context"

	"lib.kevinlin.info/aperture/internal/errors"
)

// Async is an abstraction over a Transport that internally queues transmission of individual
// payloads for asynchronous delivery by a background flush thread.
type Async struct {
	queue  chan []byte
	closed bool
	cancel context.CancelFunc
	Transport
}

// NewAsync creates a new asynchronous transport over an existing Transport backend.
func NewAsync(backendFactory Factory, size int) (Transport, error) {
	if size <= 0 {
		return nil, errors.New(
			"transport",
			"async transport queue size must be greater than 0",
		)
	}

	tport, err := backendFactory()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	async := &Async{
		queue:     make(chan []byte, size),
		cancel:    cancel,
		Transport: tport,
	}

	go async.consume(ctx)

	return async, nil
}

// Send asynchronously queues the payload for transmission. Note that this routine will reject sends
// beyond the allocated queue capacity.
func (a *Async) Send(data []byte) (int, error) {
	if a.closed {
		return 0, errors.New("transport", "async transport is closing")
	}

	select {
	case a.queue <- data:
	default:
		return 0, errors.New("transport", "async transport delivery queue is full")
	}

	return len(data), nil
}

// consume is a cancellable routine for dequeueing payloads from the queue and sending them over the
// backing transport.
func (a *Async) consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case payload := <-a.queue:
			a.Transport.Send(payload)
		}
	}
}

// Close closes all internal structures, drains any pending payloads, and closes the backing
// transport.
func (a *Async) Close() error {
	// Cancel the ongoing asynchronous consumption and prevent additional enqueues
	a.closed = true
	a.cancel()
	close(a.queue)

	// Drain any remaining payloads from the queue
	for payload := range a.queue {
		a.Transport.Send(payload)
	}

	// Close the underlying transport
	return a.Transport.Close()
}
