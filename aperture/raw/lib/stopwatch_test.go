package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// fakeTimeProvider implements timeProvider but returns static, pre-configured times.
type fakeTimeProvider struct {
	nows       []time.Time
	invocation int
}

// Now statefully returns the next time in the nows input array, advancing one index each time it is
// invoked.
func (p *fakeTimeProvider) Now() time.Time {
	if p.invocation >= len(p.nows) {
		return time.Time{}
	}

	defer func() { p.invocation++ }()

	return p.nows[p.invocation]
}

func TestFakeTimeProvider(t *testing.T) {
	t.Parallel()

	provider := &fakeTimeProvider{
		nows: []time.Time{
			time.Unix(1, 0),
			time.Unix(2, 0),
		},
	}

	assert.Equal(t, int64(1), provider.Now().Unix())
	assert.Equal(t, int64(2), provider.Now().Unix())
	assert.True(t, provider.Now().IsZero())
}

func TestNewStopwatch(t *testing.T) {
	t.Parallel()

	stopwatch := NewStopwatch()

	assert.Equal(t, running, stopwatch.state)
	assert.GreaterOrEqual(t, stopwatch.Elapsed().Nanoseconds(), int64(0))
}

func TestStopwatchElapsed(t *testing.T) {
	t.Parallel()

	stopwatch := &Stopwatch{
		provider: &fakeTimeProvider{
			nows: []time.Time{
				time.Unix(0, 0), // start
				time.Unix(1, 0), // first elapsed
				time.Unix(2, 0), // second elapsed
			},
		},
	}

	stopwatch.Start()
	stopwatch.Start() // idempotence
	assert.Equal(t, 1*time.Second, stopwatch.Elapsed())
	assert.Equal(t, 2*time.Second, stopwatch.Elapsed())
}

func TestStopwatchPause(t *testing.T) {
	t.Parallel()

	stopwatch := &Stopwatch{
		provider: &fakeTimeProvider{
			nows: []time.Time{
				time.Unix(0, 0), // first start
				time.Unix(1, 0), // first elapsed
				time.Unix(2, 0), // first pause
				time.Unix(3, 0), // second start
				time.Unix(4, 0), // second pause
			},
		},
	}

	stopwatch.Start()
	assert.Equal(t, 1*time.Second, stopwatch.Elapsed())
	stopwatch.Pause()
	stopwatch.Start()
	stopwatch.Pause()
	stopwatch.Pause() // idempotence
	assert.Equal(t, 3*time.Second, stopwatch.Elapsed())
	assert.Equal(t, 3*time.Second, stopwatch.Elapsed())
}

func TestStopwatchReset(t *testing.T) {
	t.Parallel()

	stopwatch := &Stopwatch{
		provider: &fakeTimeProvider{
			nows: []time.Time{
				time.Unix(0, 0), // start
				time.Unix(1, 0), // first elapsed
				time.Unix(2, 0), // first reset
				time.Unix(3, 0), // second elapsed
				time.Unix(4, 0), // third elapsed
				time.Unix(5, 0), // second reset
				time.Unix(6, 0), // fourth reset
				time.Unix(7, 0), // pause
				time.Unix(8, 0), // fifth reset
			},
		},
	}

	stopwatch.Start()
	assert.Equal(t, 1*time.Second, stopwatch.Elapsed())
	stopwatch.Reset()
	assert.Equal(t, 1*time.Second, stopwatch.Elapsed())
	assert.Equal(t, 2*time.Second, stopwatch.Elapsed())
	stopwatch.Reset()
	assert.Equal(t, 1*time.Second, stopwatch.Elapsed())
	stopwatch.Pause()
	stopwatch.Reset()
	assert.Equal(t, time.Duration(0), stopwatch.Elapsed())
}
