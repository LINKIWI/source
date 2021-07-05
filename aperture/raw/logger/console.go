package logger

import (
	"log"

	"lib.kevinlin.info/aperture/protocol"
)

// Console is a Logger that logs outgoing metrics to the console via standard output.
type Console struct {
	prefix     string
	serializer protocol.LineSerializer
}

// NewConsole creates a console metrics logger with default runtime parameters. All log statements
// are produced by the debug line serializer.
func NewConsole() Logger {
	return &Console{
		prefix:     "[aperture]",
		serializer: protocol.NewDebugLineSerializer(),
	}
}

// LogGauge logs a human-consumable representation of a gauge to the console.
func (l *Console) LogGauge(gauge *protocol.Gauge) {
	line, _ := l.serializer.SerializeGauge(gauge)

	log.Println(l.prefix, "gauge:", line)
}

// LogCounter logs a human-consumable representation of a counter to the console.
func (l *Console) LogCounter(counter *protocol.Counter) {
	line, _ := l.serializer.SerializeCounter(counter)

	log.Println(l.prefix, "counter:", line)
}

// LogTimer logs a human-consumable representation of a timer to the console.
func (l *Console) LogTimer(timer *protocol.Timer) {
	line, _ := l.serializer.SerializeTimer(timer)

	log.Println(l.prefix, "timer:", line)
}

// LogHistogram logs a human-consumable representation of a histogram to the console.
func (l *Console) LogHistogram(histogram *protocol.Histogram) {
	line, _ := l.serializer.SerializeHistogram(histogram)

	log.Println(l.prefix, "histogram:", line)
}
