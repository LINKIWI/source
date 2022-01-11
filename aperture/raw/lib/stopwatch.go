package lib

import (
	"time"
)

// stopwatchState describes the current state of the stopwatch.
// The stopwatch can be considered a finite state machine with only two states (running and paused),
// with bidirectional transitions between them.
type stopwatchState int

const (
	paused stopwatchState = iota
	running
)

// timeProvider is an interface for supplying details about time.
type timeProvider interface {
	// Now provides the current time.
	Now() time.Time
}

// realTimeProvider uses standard library functions to report real clock time.
type realTimeProvider struct{}

// Now retrieves the current real time.
func (p *realTimeProvider) Now() time.Time {
	return time.Now()
}

// Stopwatch reports elapsed time.
type Stopwatch struct {
	start    time.Time
	delta    time.Duration
	state    stopwatchState
	provider timeProvider
}

// NewStopwatch creates a stopwatch instance. Note that the stopwatch is started on instantiation.
func NewStopwatch() *Stopwatch {
	stopwatch := &Stopwatch{
		delta:    0,
		state:    paused,
		provider: &realTimeProvider{},
	}
	stopwatch.Start()

	return stopwatch
}

// Start starts the stopwatch. Start is an idempotent operation.
func (s *Stopwatch) Start() {
	if s.state == running {
		return
	}

	s.start = s.provider.Now()
	s.state = running
}

// Pause pauses the stopwatch; all future Elapsed invocations are expected to return the same
// duration. Pause is an idempotent operation.
func (s *Stopwatch) Pause() {
	if s.state == paused {
		return
	}

	s.delta = s.provider.Now().Add(s.delta).Sub(s.start)
	s.state = paused
}

// Reset resets the stopwatch; all future Elapsed invocations will report durations as if they had
// been started at the time Reset was invoked. Note that Reset does not create a state transition;
// the stopwatch must be manually started again if it is in a paused state.
func (s *Stopwatch) Reset() {
	s.start = s.provider.Now()
	s.delta = 0
}

// Elapsed returns a time.Duration indicating the amount of time that has elapsed since the
// stopwatch was started.
func (s *Stopwatch) Elapsed() time.Duration {
	// No need to calculate the additional elapsed time if the stopwatch is already paused.
	if s.state == paused {
		return s.delta
	}

	return s.provider.Now().Add(s.delta).Sub(s.start)
}
