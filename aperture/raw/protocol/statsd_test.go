package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatsdLineSerializerGauge(t *testing.T) {
	t.Parallel()

	cases := []struct {
		gauge      *Gauge
		serialized string
	}{
		{
			&Gauge{Name: "name", Value: 1.234},
			"name:1.234|g",
		},
		{
			&Gauge{Name: "name", Value: 1.12121212121313},
			"name:1.12121212121313|g",
		},
		{
			&Gauge{Name: "name", Value: 1.123000},
			"name:1.123|g",
		},
		{
			&Gauge{Name: "name", Value: 12345678912345},
			"name:12345678912345|g",
		},
	}

	for _, testCase := range cases {
		line, err := DefaultStatsdLineSerializer.SerializeGauge(testCase.gauge)

		assert.NoError(t, err)
		assert.Equal(t, testCase.serialized, line)
	}
}

func TestStatsdLineSerializerCounter(t *testing.T) {
	t.Parallel()

	cases := []struct {
		counter    *Counter
		serialized string
	}{
		{
			&Counter{Name: "name", Value: 1, SampleRate: 1.0},
			"name:1|c",
		},
		{
			&Counter{Name: "name", Value: -1, SampleRate: 1.0},
			"name:-1|c",
		},
		{
			&Counter{Name: "name", Value: 5, SampleRate: 0.8},
			"name:5|c|@0.8",
		},
	}

	for _, testCase := range cases {
		line, err := DefaultStatsdLineSerializer.SerializeCounter(testCase.counter)

		assert.NoError(t, err)
		assert.Equal(t, testCase.serialized, line)
	}
}

func TestStatsdLineSerializerTimer(t *testing.T) {
	t.Parallel()

	timer := &Timer{
		Name:       "name",
		Value:      1.2340,
		SampleRate: 1.0,
	}
	line, err := DefaultStatsdLineSerializer.SerializeTimer(timer)

	assert.NoError(t, err)
	assert.Equal(t, "name:1.234|ms", line)
}

func TestStatsdLineSerializerHistogram(t *testing.T) {
	t.Parallel()

	histogram := &Histogram{
		Name:  "name",
		Value: 1.234,
	}
	line, err := DefaultStatsdLineSerializer.SerializeHistogram(histogram)

	assert.NoError(t, err)
	assert.Equal(t, "name:1.234|h", line)
}

func TestStatsdLineSerializerTags(t *testing.T) {
	t.Parallel()

	gauge := &Gauge{
		Name:  "name",
		Value: 1.234,
		Tags: map[string]string{
			"foo": "bar",
			"bar": "baz",
		},
	}
	line, err := DefaultStatsdLineSerializer.SerializeGauge(gauge)

	assert.NoError(t, err)
	assert.Equal(t, "name,bar=baz,foo=bar:1.234|g", line)
}

func TestStatsdLineSerializerCustomTagDelimiters(t *testing.T) {
	t.Parallel()

	serializer := NewStatsdLineSerializer(" ", ":")

	gauge := &Gauge{
		Name:  "name",
		Value: 1.234,
		Tags: map[string]string{
			"foo": "bar",
			"bar": "baz",
		},
	}
	line, err := serializer.SerializeGauge(gauge)

	assert.NoError(t, err)
	assert.Equal(t, "name bar:baz foo:bar:1.234|g", line)
}

func TestStatsdLineSerializerQueryEscape(t *testing.T) {
	t.Parallel()

	gauge := &Gauge{
		Name:  "name:with:colons",
		Value: 1.234,
		Tags: map[string]string{
			"foo": "bar baz",
		},
	}
	line, err := DefaultStatsdLineSerializer.SerializeGauge(gauge)

	assert.NoError(t, err)
	assert.Equal(t, "name%3Awith%3Acolons,foo=bar+baz:1.234|g", line)
}
