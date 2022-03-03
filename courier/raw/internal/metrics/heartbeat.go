package metrics

import (
	"time"

	"go.uber.org/zap"
	"lib.kevinlin.info/aperture"
	"lib.kevinlin.info/aperture/lib"
)

const (
	// DefaultHeartbeatInterval is the default time interval for running heartbeat metric tasks.
	DefaultHeartbeatInterval = 10 * time.Second
)

const (
	metricHeartbeatRegister = "heartbeat.meta.register"
	metricHeartbeatStart    = "heartbeat.meta.start"
	metricHeartbeatError    = "heartbeat.meta.error"
	metricHeartbeatDuration = "heartbeat.meta.duration"
)

// Heartbeat describes a routine task that regularly sends metrics.
type Heartbeat interface {
	// Run runs the underlying heartbeat business logic. If an error is returned, the heartbeat
	// is stopped and no further executions will occur unless the heartbeat is re-registered.
	Run(client aperture.Statsd) error
}

// NewHeartbeat is a convenience function for instantiating a simple inline callback as a stateless
// heartbeat implementation.
func NewHeartbeat(run func(client aperture.Statsd) error) Heartbeat {
	return &statelessHeartbeat{run}
}

// RegisterHeartbeat schedules periodic execution of a new heartbeat task.
func RegisterHeartbeat(name string, interval time.Duration, heartbeat Heartbeat) {
	ticker := time.NewTicker(interval)
	done := make(chan bool)
	tags := map[string]interface{}{"name": name}

	Client.Incr(metricHeartbeatRegister, tags)
	zap.L().Debug(
		"registered metrics heartbeat task",
		zap.String("name", name),
		zap.Duration("interval", interval),
	)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				Client.Incr(metricHeartbeatStart, tags)
				stopwatch := lib.NewStopwatch()

				if err := heartbeat.Run(Client); err != nil {
					Client.Incr(metricHeartbeatError, tags)
					done <- true
				}

				Client.Timing(metricHeartbeatDuration, stopwatch.Elapsed(), tags)
			}
		}
	}()
}
