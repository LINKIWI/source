package protocol

import (
	"fmt"
	"sort"
	"strings"
)

// DebugLineSerializer is a LineSerializer that produces human-consumable lines summarizing each
// metric's parameters for purposes of debugging.
type DebugLineSerializer struct{}

// NewDebugLineSerializer creates a new debugging line serializer.
func NewDebugLineSerializer() LineSerializer {
	return &DebugLineSerializer{}
}

// SerializeGauge serializes a statsd gauge.
func (s *DebugLineSerializer) SerializeGauge(gauge *Gauge) (string, error) {
	return s.serializeMetric(gauge.Name, gauge.Value, gauge.Tags), nil
}

// SerializeCounter serializes a statsd counter.
func (s *DebugLineSerializer) SerializeCounter(counter *Counter) (string, error) {
	return s.serializeMetric(counter.Name, counter.Value, counter.Tags), nil
}

// SerializeTimer serializes a statsd timer.
func (s *DebugLineSerializer) SerializeTimer(timer *Timer) (string, error) {
	return s.serializeMetric(timer.Name, timer.Value, timer.Tags), nil
}

// SerializeHistogram serializes a statsd histogram.
func (s *DebugLineSerializer) SerializeHistogram(histogram *Histogram) (string, error) {
	return s.serializeMetric(histogram.Name, histogram.Value, histogram.Tags), nil
}

// serializeMetric is a common utility to format tagged metrics for human readability.
// The output format is "metric.name tag=value value"
func (s *DebugLineSerializer) serializeMetric(name string, value interface{}, tags map[string]string) string {
	var components []string
	var tagKeys []string

	// Sort tag keys to achieve deterministic serialization outputs
	for tagKey := range tags {
		tagKeys = append(tagKeys, tagKey)
	}
	sort.Strings(tagKeys)

	components = append(components, name)

	for _, tagKey := range tagKeys {
		components = append(components, fmt.Sprintf("%s=%s", tagKey, tags[tagKey]))
	}

	components = append(components, fmt.Sprintf("%v", value))

	return strings.Join(components, " ")
}
