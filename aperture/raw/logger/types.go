package logger

import (
	"lib.kevinlin.info/aperture/protocol"
)

// Logger provides interfaces for logging outgoing metrics.
type Logger interface {
	// LogGauge logs the emission of a gauge.
	LogGauge(gauge *protocol.Gauge)

	// LogCounter logs the emission of a counter.
	LogCounter(counter *protocol.Counter)

	// LogTimer logs the emission of a timer.
	LogTimer(timer *protocol.Timer)

	// LogHistogram logs the emission of a histogram.
	LogHistogram(histogram *protocol.Histogram)
}
