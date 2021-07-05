package logger

import (
	"lib.kevinlin.info/aperture/protocol"
)

// Noop is a Logger that noops for all metric logs.
type Noop struct{}

// NewNoop creates a new default Noop logger.
func NewNoop() Logger {
	return &Noop{}
}

// LogGauge noops.
func (l *Noop) LogGauge(gauge *protocol.Gauge) {}

// LogCounter noops.
func (l *Noop) LogCounter(counter *protocol.Counter) {}

// LogTimer noops.
func (l *Noop) LogTimer(timer *protocol.Timer) {}

// LogHistogram noops.
func (l *Noop) LogHistogram(histogram *protocol.Histogram) {}
