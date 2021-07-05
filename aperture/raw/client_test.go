package aperture

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"lib.kevinlin.info/aperture/internal/test"
	"lib.kevinlin.info/aperture/protocol"
	"lib.kevinlin.info/aperture/transport"
)

func TestNewClientInvalidConfig(t *testing.T) {
	t.Parallel()

	client, err := NewClient(&Config{Address: "invalid"})

	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestNewClientConfigDefaults(t *testing.T) {
	t.Parallel()

	client, err := NewClient(&Config{})

	assert.NoError(t, err)
	assert.Equal(t, transport.NewNoop(), client.tport)
	assert.Equal(t, protocol.DefaultStatsdLineSerializer, client.cfg.Serializer)
	assert.Equal(t, 1.0, client.cfg.SampleRate)
}

func TestClientMetricNamePrefixSuffix(t *testing.T) {
	t.Parallel()

	client, err := NewClient(&Config{Address: "udp://localhost:8125"})
	assert.NoError(t, err)
	assert.Equal(t, "name", client.formatName("name"))

	client, err = NewClient(&Config{Address: "udp://localhost:8125", Prefix: "prefix"})
	assert.NoError(t, err)
	assert.Equal(t, "prefix.name", client.formatName("name"))

	client, err = NewClient(&Config{Address: "udp://localhost:8125", Suffix: "suffix"})
	assert.NoError(t, err)
	assert.Equal(t, "name.suffix", client.formatName("name"))

	client, err = NewClient(&Config{
		Address: "udp://localhost:8125",
		Prefix:  "prefix",
		Suffix:  "suffix",
	})
	assert.NoError(t, err)
	assert.Equal(t, "prefix.name.suffix", client.formatName("name"))
}

func TestClientMetricTags(t *testing.T) {
	t.Parallel()

	client, err := NewClient(&Config{Address: "udp://localhost:8125"})
	assert.NoError(t, err)
	assert.Equal(
		t,
		map[string]string{"foo": "bar", "baz": "1"},
		client.mergeTags(map[string]interface{}{"foo": "bar", "baz": 1}),
	)

	client, err = NewClient(&Config{
		Address: "udp://localhost:8125",
		DefaultTags: map[string]interface{}{
			"foo": "baz",
			"bar": "foo",
		},
	})
	assert.NoError(t, err)
	assert.Equal(
		t,
		map[string]string{"foo": "bar", "baz": "1", "bar": "foo"},
		client.mergeTags(map[string]interface{}{"foo": "bar", "baz": 1}),
	)
}

func TestClientCount(t *testing.T) {
	t.Parallel()

	tport := test.NewMockTransport()
	client, _ := NewClient(&Config{Address: "udp://localhost:8125"})
	client.tport = tport

	client.Count("name", 4, nil)

	assert.Equal(t, 1, tport.NumTransmissions())
	assert.Contains(t, tport.Transmissions(), []byte("name:4|c"))
}

func TestClientIncr(t *testing.T) {
	t.Parallel()

	tport := test.NewMockTransport()
	client, _ := NewClient(&Config{Address: "udp://localhost:8125"})
	client.tport = tport

	client.Incr("name", nil)

	assert.Equal(t, 1, tport.NumTransmissions())
	assert.Contains(t, tport.Transmissions(), []byte("name:1|c"))
}

func TestClientDecr(t *testing.T) {
	t.Parallel()

	tport := test.NewMockTransport()
	client, _ := NewClient(&Config{Address: "udp://localhost:8125"})
	client.tport = tport

	client.Decr("name", nil)

	assert.Equal(t, 1, tport.NumTransmissions())
	assert.Contains(t, tport.Transmissions(), []byte("name:-1|c"))
}

func TestClientGauge(t *testing.T) {
	t.Parallel()

	tport := test.NewMockTransport()
	client, _ := NewClient(&Config{Address: "udp://localhost:8125"})
	client.tport = tport

	client.Gauge("name", -1, nil)
	client.Gauge("name", 1.0, nil)
	client.Gauge("name", 1.012345, nil)

	assert.Equal(t, 3, tport.NumTransmissions())
	assert.Contains(t, tport.Transmissions(), []byte("name:-1|g"))
	assert.Contains(t, tport.Transmissions(), []byte("name:1|g"))
	assert.Contains(t, tport.Transmissions(), []byte("name:1.012345|g"))
}

func TestClientSize(t *testing.T) {
	t.Parallel()

	tport := test.NewMockTransport()
	client, _ := NewClient(&Config{Address: "udp://localhost:8125"})
	client.tport = tport

	client.Size("name", 1, nil)
	client.Size("name", 12345, nil)

	assert.Equal(t, 2, tport.NumTransmissions())
	assert.Contains(t, tport.Transmissions(), []byte("name:1|ms"))
	assert.Contains(t, tport.Transmissions(), []byte("name:12345|ms"))
}

func TestClientTiming(t *testing.T) {
	t.Parallel()

	tport := test.NewMockTransport()
	client, _ := NewClient(&Config{Address: "udp://localhost:8125"})
	client.tport = tport

	client.Timing("name", 2*time.Second, nil)
	client.Timing("name", 2*time.Millisecond, nil)
	client.Timing("name", 123456789*time.Microsecond, nil)

	assert.Equal(t, 3, tport.NumTransmissions())
	assert.Contains(t, tport.Transmissions(), []byte("name:2000|ms"))
	assert.Contains(t, tport.Transmissions(), []byte("name:2|ms"))
	assert.Contains(t, tport.Transmissions(), []byte("name:123456.789|ms"))
}

func TestClientTimingMs(t *testing.T) {
	t.Parallel()

	tport := test.NewMockTransport()
	client, _ := NewClient(&Config{Address: "udp://localhost:8125"})
	client.tport = tport

	client.TimingMs("name", 12345, nil)
	client.TimingMs("name", 12345.012345, nil)

	assert.Equal(t, 2, tport.NumTransmissions())
	assert.Contains(t, tport.Transmissions(), []byte("name:12345|ms"))
	assert.Contains(t, tport.Transmissions(), []byte("name:12345.012345|ms"))
}

func TestClientHistogram(t *testing.T) {
	t.Parallel()

	tport := test.NewMockTransport()
	client, _ := NewClient(&Config{Address: "udp://localhost:8125"})
	client.tport = tport

	client.Histogram("name", 12345, nil)
	client.Histogram("name", 12345.012345, nil)

	assert.Equal(t, 2, tport.NumTransmissions())
	assert.Contains(t, tport.Transmissions(), []byte("name:12345|h"))
	assert.Contains(t, tport.Transmissions(), []byte("name:12345.012345|h"))
}
