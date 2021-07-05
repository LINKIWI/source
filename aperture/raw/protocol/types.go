package protocol

import (
	"time"
)

// Identifier is the statsd protocol metric type identifier string.
type identifier string

// LineSerializer describes a protocol serializer for various metric types.
type LineSerializer interface {
	// SerializeGauge serializes a Gauge to a protocol-compliant line.
	SerializeGauge(gauge *Gauge) (string, error)

	// SerializeCounter serializes a Counter to a protocol-compliant line.
	SerializeCounter(counter *Counter) (string, error)

	// SerializeTimer serializes a Timer to a protocol-compliant line.
	SerializeTimer(timer *Timer) (string, error)

	// SerializeHistogram serializes a Histogram to a protocol-compliant line.
	SerializeHistogram(histogram *Histogram) (string, error)
}

var (
	// DefaultStatsdLineSerializer is a statsd protocol LineSerializer implementation that uses
	// default tag delimiters and joiners. The default delimiters are compatible with InfluxDB.
	DefaultStatsdLineSerializer = NewStatsdLineSerializer(defaultTagDelimiter, defaultTagJoiner)
	// DefaultInfluxLineSerializer is an InfluxDB line protocol LineSerializer implementation.
	DefaultInfluxLineSerializer = NewInfluxLineSerializer()
)

// The below struct types are formalizations of metric types defined in the statsd protocol.
// See: https://github.com/b/statsd_spec/blob/master/README.md

// Gauge is a statsd gauge metric.
type Gauge struct {
	Timestamp time.Time
	Name      string
	Value     float64
	Tags      map[string]string
}

// Counter is a statsd counter metric.
type Counter struct {
	Timestamp  time.Time
	Name       string
	Value      int64
	SampleRate float64
	Tags       map[string]string
}

// Timer is a statsd timing metric.
type Timer struct {
	Timestamp  time.Time
	Name       string
	Value      float64
	SampleRate float64
	Tags       map[string]string
}

// Histogram is a statsd histogram metric.
type Histogram struct {
	Timestamp time.Time
	Name      string
	Value     float64
	Tags      map[string]string
}
