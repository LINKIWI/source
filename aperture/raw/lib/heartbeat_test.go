package lib

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterHeartbeat(t *testing.T) {
	t.Parallel()

	var executed int

	heartbeat := NewStatelessHeartbeat("dummy", func() error {
		executed++

		if executed == 3 {
			return fmt.Errorf("heartbeat")
		}

		return nil
	})

	err := <-RegisterHeartbeat(heartbeat, 1*time.Nanosecond, 0)

	assert.Equal(t, 3, executed)
	assert.Error(t, err)
}

func TestUptimeHeartbeat(t *testing.T) {
	t.Parallel()

	heartbeat := NewUptimeHeartbeat(NewNoopStatsd())

	err := heartbeat.Run()
	assert.NoError(t, err)
}

func TestRuntimeStatsHeartbeat(t *testing.T) {
	t.Parallel()

	heartbeat := NewRuntimeStatsHeartbeat(NewNoopStatsd())

	err := heartbeat.Run()
	assert.NoError(t, err)
}

func TestResourceUsageHeartbeat(t *testing.T) {
	t.Parallel()

	heartbeat := NewResourceUsageHeartbeat(NewNoopStatsd())

	err := heartbeat.Run()
	assert.NoError(t, err)
}
