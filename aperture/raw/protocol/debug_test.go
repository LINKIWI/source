package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	debugSerializer = NewDebugLineSerializer()
)

func TestDebugLineSerializerGauge(t *testing.T) {
	t.Parallel()

	gauge := &Gauge{
		Name:  "name",
		Value: 1.234,
	}
	line, err := debugSerializer.SerializeGauge(gauge)

	assert.NoError(t, err)
	assert.Equal(t, "name 1.234", line)
}

func TestDebugLineSerializerCounter(t *testing.T) {
	t.Parallel()

	counter := &Counter{
		Name:       "name",
		Value:      4,
		SampleRate: 0.7,
	}
	line, err := debugSerializer.SerializeCounter(counter)

	assert.NoError(t, err)
	assert.Equal(t, "name 4", line)
}

func TestDebugLineSerializerTimer(t *testing.T) {
	t.Parallel()

	timer := &Timer{
		Name:  "name",
		Value: 1.2340,
	}
	line, err := debugSerializer.SerializeTimer(timer)

	assert.NoError(t, err)
	assert.Equal(t, "name 1.234", line)
}

func TestDebugLineSerializerHistogram(t *testing.T) {
	t.Parallel()

	histogram := &Histogram{
		Name:  "name",
		Value: 1.234,
	}
	line, err := debugSerializer.SerializeHistogram(histogram)

	assert.NoError(t, err)
	assert.Equal(t, "name 1.234", line)
}

func TestDebugLineSerializerTags(t *testing.T) {
	t.Parallel()

	gauge := &Gauge{
		Name:  "name",
		Value: 1.234,
		Tags: map[string]string{
			"foo": "bar",
			"bar": "baz",
		},
	}
	line, err := debugSerializer.SerializeGauge(gauge)

	assert.NoError(t, err)
	assert.Equal(t, "name bar=baz foo=bar 1.234", line)
}
