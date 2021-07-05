package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"lib.kevinlin.info/aperture"
)

// slowStatsd wraps a NoopStatsd but injects an artificial delay in its emissions APIs.
type slowStatsd struct {
	delay time.Duration
	NoopStatsd
}

func (s *slowStatsd) Count(name string, delta int64, tags map[string]interface{}) error {
	time.Sleep(s.delay)
	return nil
}

func TestNoopStastd(t *testing.T) {
	t.Parallel()

	var noop aperture.Statsd = NewNoopStatsd()

	assert.NoError(t, noop.Count("", 0, nil))
	assert.NoError(t, noop.Incr("", nil))
	assert.NoError(t, noop.Decr("", nil))
	assert.NoError(t, noop.Gauge("", 0, nil))
	assert.NoError(t, noop.Size("", 0, nil))
	assert.NoError(t, noop.Timing("", 0, nil))
	assert.NoError(t, noop.TimingMs("", 0, nil))
	assert.NoError(t, noop.Histogram("", 0, nil))
	assert.NoError(t, noop.Close())
}

func TestAsyncStatsd(t *testing.T) {
	t.Parallel()

	var async aperture.Statsd = NewAsyncStatsd(&slowStatsd{delay: 10 * time.Minute})

	done := make(chan bool)
	go func() {
		assert.NoError(t, async.Count("", 0, nil))
		done <- true
	}()

	// A slow Count invocation should still return immediately due to the asynchronous write;
	// if the call doesn't complete on the order of seconds, it's highly likely it's blocking
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		assert.Fail(t, "expected nonblocking invocation: failed to return within timeout")
	}

	assert.NoError(t, async.Incr("", nil))
	assert.NoError(t, async.Decr("", nil))
	assert.NoError(t, async.Gauge("", 0, nil))
	assert.NoError(t, async.Size("", 0, nil))
	assert.NoError(t, async.Timing("", 0, nil))
	assert.NoError(t, async.TimingMs("", 0, nil))
	assert.NoError(t, async.Histogram("", 0, nil))
	assert.NoError(t, async.Close())
}
