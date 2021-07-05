package protocol

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"lib.kevinlin.info/aperture/internal/data"
)

const (
	// Static field name for the metric value
	valueFieldName = "value"
	// Static tag name for the metric type tag
	metricTypeTagName = "metric_type"
)

var (
	// String sequence replacement maps for escaping different components within a metric line.
	// The protocol has a limited set of characters that must be explicitly escaped with a
	// leading backslash.
	// Reference: https://docs.influxdata.com/influxdb/v2.0/reference/syntax/line-protocol/#special-characters
	measurementEscapeReplacements = []data.Replacement{
		{From: "\\", To: "\\\\"},
		{From: ",", To: "\\,"},
		{From: " ", To: "\\ "},
	}
	tagEscapeReplacements = []data.Replacement{
		{From: "\\", To: "\\\\"},
		{From: "=", To: "\\="},
		{From: ",", To: "\\,"},
		{From: " ", To: "\\ "},
	}

	// Map of statsd metric types to human-consumable names.
	// Used exclusively for augmenting outgoing metric tags, since there is no per-type statsd
	// aggregation when writing measurements to InfluxDB directly.
	metricTypes = map[identifier]string{
		gaugeType:     "gauge",
		counterType:   "counter",
		timerType:     "timer",
		histogramType: "histogram",
	}
)

// InfluxLineSerializer implements LineSerializer for the InfluxDB line protocol.
type InfluxLineSerializer struct{}

// NewInfluxLineSerializer creates a InfluxDB line protocol serializer.
func NewInfluxLineSerializer() LineSerializer {
	return &InfluxLineSerializer{}
}

// SerializeGauge serializes a statsd gauge.
func (i *InfluxLineSerializer) SerializeGauge(gauge *Gauge) (string, error) {
	return i.serializeMetric(gauge.Name, gauge.Value, gauge.Tags, gauge.Timestamp, gaugeType), nil
}

// SerializeCounter serializes a statsd counter.
func (i *InfluxLineSerializer) SerializeCounter(counter *Counter) (string, error) {
	return i.serializeMetric(counter.Name, counter.Value, counter.Tags, counter.Timestamp, counterType), nil
}

// SerializeTimer serializes a statsd timer.
func (i *InfluxLineSerializer) SerializeTimer(timer *Timer) (string, error) {
	return i.serializeMetric(timer.Name, timer.Value, timer.Tags, timer.Timestamp, timerType), nil
}

// SerializeHistogram serializes a statsd histogram.
func (i *InfluxLineSerializer) SerializeHistogram(histogram *Histogram) (string, error) {
	return i.serializeMetric(histogram.Name, histogram.Value, histogram.Tags, histogram.Timestamp, histogramType), nil
}

// serializeMetric is a common utility to convert a statsd struct into an InfluxDB protocol line.
func (i *InfluxLineSerializer) serializeMetric(
	name string,
	value interface{},
	tags map[string]string,
	timestamp time.Time,
	id identifier,
) string {
	/* Serialize tags */

	tagComponents := []string{fmt.Sprintf("%s=%s", metricTypeTagName, metricTypes[id])}

	for tagKey, tagValue := range tags {
		tagComponents = append(
			tagComponents,
			fmt.Sprintf(
				"%s=%s",
				data.EscapeString(tagKey, tagEscapeReplacements),
				data.EscapeString(tagValue, tagEscapeReplacements),
			),
		)
	}

	// Sort tag components to achieve deterministic serialization outputs
	sort.Strings(tagComponents)

	/* Serialize field value */

	field := ""

	switch value := value.(type) {
	case float64:
		field = strconv.FormatFloat(value, 'f', -1, 64)
	case int64:
		field = strconv.Itoa(int(value))
	default:
		field = fmt.Sprintf("%v", value)
	}

	return fmt.Sprintf(
		"%s,%s %s=%s %d",
		data.EscapeString(name, measurementEscapeReplacements),
		strings.Join(tagComponents, ","),
		valueFieldName,
		field,
		timestamp.UnixNano(),
	)
}
