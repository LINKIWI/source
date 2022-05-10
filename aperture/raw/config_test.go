package aperture

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidationFailure(t *testing.T) {
	t.Parallel()

	cases := []*Config{
		// URI parse failure
		{
			Address: "invalid",
		},
		// Unsupported URI scheme
		{
			Address: "unknown://localhost:8125",
		},
		// Address parse failure
		{
			Address: "udp://localhost",
		},
		// Socket stat failure
		{
			Address: "unix:///",
		},
		// File stat of a non-regular file
		{
			Address: "file:///",
		},
		// Invalid buffer size
		{
			Address:    "udp://localhost:8125",
			BufferSize: -1,
		},
		// Invalid async queue size
		{
			Address:        "udp://localhost:8125",
			AsyncQueueSize: -1,
		},
		// Proxy address parse failure
		{
			Address: "udp://localhost:8125",
			Proxy:   "localhost:3125",
		},
		// Unsupported proxy server URI scheme
		{
			Address: "udp://localhost:8125",
			Proxy:   "udp://localhost:3125",
		},
		// Invalid transport probe interval
		{
			Address:                "udp://localhost:8125",
			TransportProbeInterval: -1 * time.Second,
		},
		// Invalid sample rate
		{
			Address:    "udp://localhost:8125",
			SampleRate: 1.1,
		},
		{
			Address:    "udp://localhost:8125",
			SampleRate: -1.2,
		},
		// Invalid single address in a composite address
		{
			Address: ",",
		},
		{
			Address: "udp://localhost:8125,tcp://",
		},
		{
			Address: "udp://localhost:8125,tcp://localhost:8125,",
		},
		{
			Address: "udp://localhost:8125,",
		},
	}

	for _, cfg := range cases {
		assert.Error(t, cfg.validate())
	}
}

func TestConfigValidationSuccess(t *testing.T) {
	t.Parallel()

	cases := []*Config{
		{},
		{
			Address: "udp://localhost:8125",
		},
		{
			Address: "tcp://localhost:8125",
		},
		{
			Address: "tcp://localhost:8125,udp://localhost:8125",
		},
		{
			Address: "udp://localhost:8125",
			Proxy:   "tcp://localhost:3125",
		},
		{
			Address:    "udp://localhost:8125",
			BufferSize: 1,
		},
		{
			Address:    "udp://localhost:8125",
			SampleRate: 0.8,
			BufferSize: 2,
		},
		{
			Address: "udp://localhost:8125",
			Prefix:  "bar",
			DefaultTags: map[string]interface{}{
				"host": "foo",
				"ncpu": 6,
			},
		},
	}

	for _, cfg := range cases {
		assert.NoError(t, cfg.validate())
	}
}
