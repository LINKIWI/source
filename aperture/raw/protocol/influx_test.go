package protocol

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInfluxLineSerializerGauge(t *testing.T) {
	t.Parallel()

	cases := []struct {
		gauge      *Gauge
		serialized string
	}{
		{
			&Gauge{Name: "name", Value: 1.234, Timestamp: time.Unix(100, 0)},
			"name,metric_type=gauge value=1.234 100000000000",
		},
		{
			&Gauge{Name: "name", Value: 1.12121212121313, Timestamp: time.Unix(100, 0)},
			"name,metric_type=gauge value=1.12121212121313 100000000000",
		},
		{
			&Gauge{Name: "name", Value: 1.123000, Timestamp: time.Unix(100, 0)},
			"name,metric_type=gauge value=1.123 100000000000",
		},
		{
			&Gauge{Name: "name", Value: 12345678912345, Timestamp: time.Unix(100, 0)},
			"name,metric_type=gauge value=12345678912345 100000000000",
		},
	}

	for _, testCase := range cases {
		line, err := DefaultInfluxLineSerializer.SerializeGauge(testCase.gauge)

		assert.NoError(t, err)
		assert.Equal(t, testCase.serialized, line)
	}
}

func TestInfluxLineSerializerCounter(t *testing.T) {
	t.Parallel()

	cases := []struct {
		counter    *Counter
		serialized string
	}{
		{
			&Counter{Name: "name", Value: 1, Timestamp: time.Unix(100, 0)},
			"name,metric_type=counter value=1 100000000000",
		},
		{
			&Counter{Name: "name", Value: -1, Timestamp: time.Unix(100, 0)},
			"name,metric_type=counter value=-1 100000000000",
		},
		{
			&Counter{Name: "name", Value: 5, Timestamp: time.Unix(100, 0)},
			"name,metric_type=counter value=5 100000000000",
		},
	}

	for _, testCase := range cases {
		line, err := DefaultInfluxLineSerializer.SerializeCounter(testCase.counter)

		assert.NoError(t, err)
		assert.Equal(t, testCase.serialized, line)
	}
}

func TestInfluxLineSerializerTimer(t *testing.T) {
	t.Parallel()

	timer := &Timer{
		Name:      "name",
		Value:     1.2340,
		Timestamp: time.Unix(100, 0),
	}
	line, err := DefaultInfluxLineSerializer.SerializeTimer(timer)

	assert.NoError(t, err)
	assert.Equal(t, "name,metric_type=timer value=1.234 100000000000", line)
}

func TestInfluxLineSerializerHistogram(t *testing.T) {
	t.Parallel()

	histogram := &Histogram{
		Name:      "name",
		Value:     1.234,
		Timestamp: time.Unix(100, 0),
	}
	line, err := DefaultInfluxLineSerializer.SerializeHistogram(histogram)

	assert.NoError(t, err)
	assert.Equal(t, "name,metric_type=histogram value=1.234 100000000000", line)
}

func TestInfluxLineSerializerTags(t *testing.T) {
	t.Parallel()

	gauge := &Gauge{
		Name:      "name",
		Value:     1.234,
		Timestamp: time.Unix(100, 0),
		Tags: map[string]string{
			"foo": "bar",
			"bar": "baz",
		},
	}
	line, err := DefaultInfluxLineSerializer.SerializeGauge(gauge)

	assert.NoError(t, err)
	assert.Equal(t, "name,bar=baz,foo=bar,metric_type=gauge value=1.234 100000000000", line)
}

func TestInfluxLineSerializerSpecialCharacterEscape(t *testing.T) {
	t.Parallel()

	gauge := &Gauge{
		Name:      "\\na,m e",
		Value:     1.234,
		Timestamp: time.Unix(100, 0),
		Tags: map[string]string{
			"\\,foo": "b=ar",
			"b=ar":   "ba z",
		},
	}
	line, err := DefaultInfluxLineSerializer.SerializeGauge(gauge)

	assert.NoError(t, err)
	assert.Equal(
		t,
		"\\\\na\\,m\\ e,\\\\\\,foo=b\\=ar,b\\=ar=ba\\ z,metric_type=gauge value=1.234 100000000000",
		line,
	)
}
