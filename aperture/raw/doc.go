// Package aperture provides a statsd client library with native tags integration, support for
// multiple underlying network transports, configurable buffering/batching behavior, and more.
// By default, Aperture implements the statsd protocol, but can be extended by the client to support
// any protocol (including non-standard statsd specifications) that shares the same metric types
// (gauges, counters, timers, and histograms).
package aperture
