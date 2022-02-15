package protocol

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

const (
	// statsd default tag delimiters, compatible with Telegraf/InfluxDB
	defaultTagDelimiter = ","
	defaultTagJoiner    = "="

	// statsd line protocol name, value, and type delimiters
	valueDelimiter = ":"
	typeDelimiter  = "|"

	// statsd protocol metric type constants
	gaugeType     identifier = "g"
	counterType   identifier = "c"
	timerType     identifier = "ms"
	histogramType identifier = "h"
)

// StatsdLineSerializer implements LineSerializer for the statsd protocol.
type StatsdLineSerializer struct {
	tagDelimiter string
	tagJoiner    string
}

// NewStatsdLineSerializer creates a statsd line serializer with custom serialization options.
func NewStatsdLineSerializer(tagDelimiter string, tagJoiner string) LineSerializer {
	return &StatsdLineSerializer{
		tagDelimiter: tagDelimiter,
		tagJoiner:    tagJoiner,
	}
}

// SerializeGauge serializes a statsd gauge.
func (s *StatsdLineSerializer) SerializeGauge(gauge *Gauge) (string, error) {
	return s.serializeMetric(gauge.Name, gauge.Value, gauge.Tags, gaugeType), nil
}

// SerializeCounter serializes a statsd counter.
func (s *StatsdLineSerializer) SerializeCounter(counter *Counter) (string, error) {
	metric := s.serializeMetric(counter.Name, counter.Value, counter.Tags, counterType)

	if counter.SampleRate == 1.0 {
		return metric, nil
	}

	return fmt.Sprintf("%s%s@%v", metric, typeDelimiter, counter.SampleRate), nil
}

// SerializeTimer serializes a statsd timer.
func (s *StatsdLineSerializer) SerializeTimer(timer *Timer) (string, error) {
	metric := s.serializeMetric(timer.Name, timer.Value, timer.Tags, timerType)

	if timer.SampleRate == 1.0 {
		return metric, nil
	}

	return fmt.Sprintf("%s%s@%v", metric, typeDelimiter, timer.SampleRate), nil
}

// SerializeHistogram serializes a statsd histogram.
func (s *StatsdLineSerializer) SerializeHistogram(histogram *Histogram) (string, error) {
	return s.serializeMetric(histogram.Name, histogram.Value, histogram.Tags, histogramType), nil
}

// serializeMetric is a common utility to format tagged metrics in accordance with statsd
// protocol conventions.
func (s *StatsdLineSerializer) serializeMetric(
	name string,
	value interface{},
	tags map[string]string,
	id identifier,
) string {
	/* Serialize tags */

	idx := 0
	tagComponents := make([]string, len(tags))

	for tagKey, tagValue := range tags {
		tagComponents[idx] += url.QueryEscape(tagKey)
		tagComponents[idx] += s.tagJoiner
		tagComponents[idx] += url.QueryEscape(tagValue)

		idx++
	}

	// Sort tag components to achieve deterministic serialization outputs
	if len(tagComponents) > 1 {
		sort.Strings(tagComponents)
	}

	/* Assemble metric name, value, and type */

	var metric strings.Builder

	// Some characters, like colons, are incompatible with the statsd protocol.
	// This standardizes on URL escaping to encode such characters that may appear in the metric
	// name or tag keys/values.
	metric.WriteString(url.QueryEscape(name))

	if len(tagComponents) > 0 {
		metric.WriteString(s.tagDelimiter)
		metric.WriteString(strings.Join(tagComponents, s.tagDelimiter))
	}

	metric.WriteString(valueDelimiter)

	switch v := value.(type) {
	case float64:
		metric.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
	case int64:
		metric.WriteString(strconv.Itoa(int(v)))
	default:
		metric.WriteString(fmt.Sprintf("%v", v))
	}

	metric.WriteString(typeDelimiter)
	metric.WriteString(string(id))

	return metric.String()
}
