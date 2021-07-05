package transport

import (
	"sync"
	"sync/atomic"
	"time"
)

// Reconnecting is an thread-safe abstraction over a Transport that transparently attempts to
// reestablish connections to the underlying transport at regular intervals if the transmission
// error count exceeds a defined threshold.
type Reconnecting struct {
	backendFactory Factory
	duration       time.Duration
	threshold      int64
	signal         chan bool
	errs           int64
	mutex          sync.Mutex
	Transport
}

// NewReconnecting creates an automatically reconnecting transport over an existing transport
// backend. Note that it is expected that the callee supply a transport factory that may be invoked
// lazily, since the client will reestablish connections by simply closing the existing transport
// and creating a new one.
func NewReconnecting(backendFactory Factory, duration time.Duration, threshold int64) (Transport, error) {
	backend, err := backendFactory()
	if err != nil {
		return backend, err
	}

	tport := &Reconnecting{
		backendFactory: backendFactory,
		duration:       duration,
		threshold:      threshold,
		signal:         make(chan bool),
		errs:           0,
		Transport:      backend,
	}

	// Start an asynchronous timer to probe the current state of the transport and create a new
	// instance as necessary.
	if duration > 0 {
		go tport.probe()
	}

	return tport, nil
}

// Send defers to the underlying transport to transmit the data. Internally, the number of send
// errors from the backing transport is atomically tracked for evaluation in an asynchronous check
// path. Note that failed writes are silently ignored (considered lost) and not buffered for retry.
func (t *Reconnecting) Send(data []byte) (int, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	n, err := t.Transport.Send(data)
	if err != nil {
		atomic.AddInt64(&t.errs, 1)
	}

	return n, err
}

// Close defers to the underlying transport to close the transport. It also requests a graceful stop
// to the internal ticker for evaluating the current transport error count.
func (t *Reconnecting) Close() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.signal <- true

	return t.Transport.Close()
}

// probe periodically checks the current error count at an interval defined by duration. If the
// error count exceeds the error threshold, it re-instantiates the underlying backend transport.
// Periodic probing stops when the stop signal channel is populated.
func (t *Reconnecting) probe() {
	ticker := time.NewTicker(t.duration)

	for {
		select {
		case <-t.signal:
			ticker.Stop()
			return
		case <-ticker.C:
			if atomic.LoadInt64(&t.errs) >= t.threshold {
				t.mutex.Lock()

				// Only reset the transport if we can successfully create another
				// instance. Otherwise, silently ignore the error and allow the
				// transmission error count to pile up. The next iteration of the
				// loop will attempt instantiation again.
				if tport, err := t.backendFactory(); err == nil {
					t.Transport.Close()
					t.Transport = tport
					atomic.StoreInt64(&t.errs, 0)
				}

				t.mutex.Unlock()
			}
		}
	}
}
