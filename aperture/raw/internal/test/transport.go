package test

import (
	"fmt"
)

// MockTransportState describes the open/closed state of the mock transport.
type MockTransportState int

const (
	// ActiveTransport implies that the transport is open (like an established TCP connection)
	ActiveTransport MockTransportState = iota
	// ClosedTransport implies that the transport is closed (like a closed TCP connection)
	ClosedTransport
)

// MockTransport is a fake transport.Transport implementation for use in unit tests.
type MockTransport struct {
	transmissions [][]byte
	state         MockTransportState
}

// NewMockTransport creates a new mock transport.
func NewMockTransport() *MockTransport {
	return &MockTransport{
		state: ActiveTransport,
	}
}

// Send records the payload in internal state.
func (t *MockTransport) Send(data []byte) (int, error) {
	t.transmissions = append(t.transmissions, data)

	return len(data), nil
}

// Close changes the internal transport state.
func (t *MockTransport) Close() error {
	t.state = ClosedTransport

	return nil
}

// NumTransmissions returns the number of transmissions (i.e. Send invocations) received.
func (t *MockTransport) NumTransmissions() int {
	return len(t.transmissions)
}

// Transmissions returns all recorded transmissions.
func (t *MockTransport) Transmissions() [][]byte {
	return t.transmissions
}

// State reads the current transport state.
func (t *MockTransport) State() MockTransportState {
	return t.state
}

// ErrorTransport is a fake transport.Transport implementation that mimics MockTransport but
// deliberately fails every send with an error.
type ErrorTransport struct {
	*MockTransport
}

// NewErrorTransport creates a new mock error transport.
func NewErrorTransport() *ErrorTransport {
	return &ErrorTransport{
		MockTransport: NewMockTransport(),
	}
}

// Send records the attempted transmission but always returns a static error.
func (e *ErrorTransport) Send(data []byte) (int, error) {
	e.MockTransport.Send(data)

	return 0, fmt.Errorf("")
}

// StalledTransport is a fake transport.Transport implementation that mimics MockTransport but
// hangs indefinitely on every send to simulate a stalled connection.
type StalledTransport struct {
	*MockTransport
}

// NewStalledTransport creates a new mock stalled transport.
func NewStalledTransport() *StalledTransport {
	return &StalledTransport{
		MockTransport: NewMockTransport(),
	}
}

// Send records the attempted transmission but stalls indefinitely.
func (s *StalledTransport) Send(data []byte) (int, error) {
	s.MockTransport.Send(data)

	<-make(chan struct{})

	return 0, nil
}
