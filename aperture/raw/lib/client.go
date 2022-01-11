package lib

import (
	"time"

	"lib.kevinlin.info/aperture"
)

// NoopStatsd implements the Statsd interface but noops on all APIs.
type NoopStatsd struct{}

// NewNoopStatsd is a convenience function for creating a NoopStatsd.
func NewNoopStatsd() *NoopStatsd {
	return &NoopStatsd{}
}

// Count is a noop.
func (n *NoopStatsd) Count(name string, delta int64, tags map[string]interface{}) error {
	return nil
}

// Incr is a noop.
func (n *NoopStatsd) Incr(name string, tags map[string]interface{}) error {
	return nil
}

// Decr is a noop.
func (n *NoopStatsd) Decr(name string, tags map[string]interface{}) error {
	return nil
}

// Gauge is a noop.
func (n *NoopStatsd) Gauge(name string, value float64, tags map[string]interface{}) error {
	return nil
}

// Size is a noop.
func (n *NoopStatsd) Size(name string, bytes int64, tags map[string]interface{}) error {
	return nil
}

// Timing is a noop.
func (n *NoopStatsd) Timing(name string, duration time.Duration, tags map[string]interface{}) error {
	return nil
}

// TimingMs is a noop.
func (n *NoopStatsd) TimingMs(name string, duration float64, tags map[string]interface{}) error {
	return nil
}

// Histogram is a noop.
func (n *NoopStatsd) Histogram(name string, value float64, tags map[string]interface{}) error {
	return nil
}

// Close is a noop.
func (n *NoopStatsd) Close() error {
	return nil
}

// AsyncStatsd wraps a Statsd implementation to write all metrics asynchronously.
type AsyncStatsd struct {
	aperture.Statsd
}

// NewAsyncStatsd creates an AsyncStatsd wrapping an existing Statsd implementation.
// Note that, due to the asynchronous nature of the client, write errors are not surfaced; all APIs
// return a nil error even if an error occurs internally.
func NewAsyncStatsd(statsd aperture.Statsd) *AsyncStatsd {
	return &AsyncStatsd{statsd}
}

// Count invokes the underlying Count in a goroutine.
func (a *AsyncStatsd) Count(name string, delta int64, tags map[string]interface{}) error {
	go a.Statsd.Count(name, delta, tags)
	return nil
}

// Incr invokes the underlying Incr in a goroutine.
func (a *AsyncStatsd) Incr(name string, tags map[string]interface{}) error {
	go a.Statsd.Incr(name, tags)
	return nil
}

// Decr invokes the underlying Decr in a goroutine.
func (a *AsyncStatsd) Decr(name string, tags map[string]interface{}) error {
	go a.Statsd.Decr(name, tags)
	return nil
}

// Gauge invokes the underlying Decr in a goroutine.
func (a *AsyncStatsd) Gauge(name string, value float64, tags map[string]interface{}) error {
	go a.Statsd.Gauge(name, value, tags)
	return nil
}

// Size invokes the underlying Size in a goroutine.
func (a *AsyncStatsd) Size(name string, bytes int64, tags map[string]interface{}) error {
	go a.Statsd.Size(name, bytes, tags)
	return nil
}

// Timing invokes the underlying Timing in a goroutine.
func (a *AsyncStatsd) Timing(name string, duration time.Duration, tags map[string]interface{}) error {
	go a.Statsd.Timing(name, duration, tags)
	return nil
}

// TimingMs invokes the underlying TimingMs in a goroutine.
func (a *AsyncStatsd) TimingMs(name string, duration float64, tags map[string]interface{}) error {
	go a.Statsd.TimingMs(name, duration, tags)
	return nil
}

// Histogram invokes the underlying Histogram in a goroutine.
func (a *AsyncStatsd) Histogram(name string, value float64, tags map[string]interface{}) error {
	go a.Statsd.Histogram(name, value, tags)
	return nil
}
