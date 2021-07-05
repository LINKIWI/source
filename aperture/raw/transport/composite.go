package transport

import (
	"lib.kevinlin.info/aperture/internal/errors"
)

// writeConsistencyLevel expresses the desired consistency level for individual transmissions.
type writeConsistencyLevel int

const (
	// None does not block until the dispatched writes are completed and thus considers all
	// writes to be immediately successful.
	None writeConsistencyLevel = iota
	// First blocks only until the first transport (according to the order in the transport
	// backends list) completes and uses its result as the result of the composite write.
	First
	// Any blocks until the first of (1) any transport completes successfully or (2) all
	// transports fail. In other words, the composite write is considered successful if one or
	// more individual transports succeed; the composite write is considered failed if all
	// transports fail.
	Any
	// Majority blocks until the first of (1) the majority of transports complete successfully
	// or (2) all transports complete without reaching the required majority threshold for
	// successful writes. In other words, the composite write is considered successful as long
	// as any majority of transports completes successfully. On the other hand, the composite
	// write is considered failed if any transport write fails without reaching a majority of
	// successful writes in total.
	Majority
	// All blocks until all transports complete and returns the worst result of all transports;
	// i.e. the composite write is considered failed if any of the individual writes failed.
	All
)

// Composite is an abstraction over multiple Transports that allows concurrent writes to multiple
// transports on each send.
type Composite struct {
	consistency writeConsistencyLevel
	backends    []Transport
}

// NewComposite create a new composite transport over multiple existing Transport backends.
func NewComposite(backendFactories []Factory, consistency writeConsistencyLevel) (Transport, error) {
	var transports []Transport

	for _, backendFactory := range backendFactories {
		tport, err := backendFactory()
		if err != nil {
			return nil, err
		}

		transports = append(transports, tport)
	}

	return &Composite{backends: transports, consistency: consistency}, nil
}

// Send concurrently dispatches the data payload to all transport backends and blocks until the
// desired write consistency level is satisfied. Note that, due to the concurrent nature of the
// composite write, the number of sent bytes should be considered meaningless regardless of the
// desired write consistency level.
func (t *Composite) Send(data []byte) (int, error) {
	var errs []error
	successes := 0
	results := make(chan error, len(t.backends))

	if len(t.backends) == 0 {
		return 0, nil
	}

	// Special case: for First write consistency, synchronously dispatch to the first backend
	// and asynchronously dispatch to the rest, ignoring their results.
	if t.consistency == First {
		defer func() {
			if len(t.backends) > 1 {
				for _, tport := range t.backends[1:] {
					go tport.Send(data)
				}
			}
		}()

		return t.backends[0].Send(data)
	}

	// Asynchronously dispatch to all backends and accumulate their results in a shared channel.
	for _, tport := range t.backends {
		go func(tport Transport) {
			_, err := tport.Send(data)
			results <- err
		}(tport)
	}

	if t.consistency == None {
		return 0, nil
	}

	for i := 0; i < len(t.backends); i++ {
		err := <-results
		if err == nil {
			successes++
		} else {
			errs = append(errs, err)
		}

		// Any write consistency allows early successful termination as long as one success
		// occurs.
		if t.consistency == Any && successes > 0 {
			return 0, err
		}

		// Majority write consistency requires a threshold minimum number of successful
		// writes for early successful termination.
		if t.consistency == Majority && float32(successes)/float32(len(t.backends)) >= 0.5 {
			return 0, nil
		}
	}

	if t.consistency == All && len(errs) == 0 {
		return 0, nil
	}

	return 0, errors.New(
		"transport",
		"send in composite transport failed to satisfy desired write consistency",
		errors.Tag{Key: "consistency", Value: t.consistency},
		errors.Tag{Key: "errs", Value: errs},
	)
}

// Close attempts to close all underlying transports, accumulating errors for individual failed
// closures.
func (t *Composite) Close() error {
	var errs []error

	for _, tport := range t.backends {
		if err := tport.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(
		"transport",
		"closing composite transport had multiple errors",
		errors.Tag{Key: "errs", Value: errs},
	)
}
