package lib

import (
	"context"
	"math/rand"
	"runtime"
	"runtime/metrics"
	"syscall"
	"time"

	"lib.kevinlin.info/aperture"
)

const (
	// DefaultHeartbeatInterval is the default time interval for running heartbeat metric tasks.
	DefaultHeartbeatInterval = 10 * time.Second
	// DefaultHeartbeatJitter is the default jitter interval to pause before each invocation of
	// a heartbeat task.
	DefaultHeartbeatJitter = 0
)

// Heartbeat describes a routine task that regularly sends metrics.
type Heartbeat interface {
	// Name returns a string identifier for this heartbeat.
	Name() string

	// Run runs the underlying heartbeat business logic. If an error is returned, the heartbeat
	// is stopped and no further executions will occur unless the heartbeat is re-registered.
	Run() error
}

// RegisterHeartbeat schedules periodic execution of a new heartbeat task.
func RegisterHeartbeat(heartbeat Heartbeat, interval time.Duration, jitter time.Duration) chan error {
	var err error

	if interval == 0 {
		interval = DefaultHeartbeatInterval
	}

	ctx, cancel := context.WithCancel(context.Background())
	errC := make(chan error)
	ticker := time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				errC <- err
				return
			case <-ticker.C:
				if jitter > 0 {
					time.Sleep(time.Duration(float64(jitter) * rand.Float64()))
				}

				if err = heartbeat.Run(); err != nil {
					cancel()
				}
			}
		}
	}()

	return errC
}

// statelessHeartbeat is a Heartbeat wrapper for a stateless callback function.
type statelessHeartbeat struct {
	name string
	run  func() error
}

// NewStatelessHeartbeat creates a new statelessHeartbeat.
func NewStatelessHeartbeat(name string, run func() error) Heartbeat {
	return &statelessHeartbeat{
		name: name,
		run:  run,
	}
}

// Name returns the name of this heartbeat.
func (s *statelessHeartbeat) Name() string {
	return s.name
}

// Run runs the stateless callback function.
func (s *statelessHeartbeat) Run() error {
	return s.run()
}

// uptimeHeartbeat tracks the uptime of the application.
type uptimeHeartbeat struct {
	client aperture.Statsd
	start  time.Time
}

// NewUptimeHeartbeat creates a new uptimeHeartbeat.
func NewUptimeHeartbeat(client aperture.Statsd) Heartbeat {
	return &uptimeHeartbeat{
		start:  time.Now(),
		client: client,
	}
}

// Name returns a constant name for the uptime heartbeat.
func (u *uptimeHeartbeat) Name() string {
	return "uptime"
}

// Run emits a gauge for the amount of time since the application started.
func (u *uptimeHeartbeat) Run() error {
	tags := map[string]interface{}{"heartbeat": u.Name()}

	u.client.Gauge("uptime.seconds", time.Since(u.start).Seconds(), tags)

	return nil
}

// runtimeStatsHeartbeat reports Go runtime stats.
type runtimeStatsHeartbeat struct {
	client        aperture.Statsd
	samples       []metrics.Sample
	sampleMetrics map[string]string
}

// NewRuntimeStatsHeartbeat creates a new runtimeStatsHeartbeat with a set of default runtime stats.
func NewRuntimeStatsHeartbeat(client aperture.Statsd) Heartbeat {
	var samples []metrics.Sample

	sampleMetrics := map[string]string{
		"/gc/cycles/total:gc-cycles":          "runtime.gc.cycles",
		"/gc/heap/goal:bytes":                 "runtime.gc.heap.goal_bytes",
		"/gc/heap/objects:objects":            "runtime.gc.heap.objects",
		"/memory/classes/heap/free:bytes":     "runtime.memory.heap.free_bytes",
		"/memory/classes/heap/objects:bytes":  "runtime.memory.heap.object_bytes",
		"/memory/classes/heap/released:bytes": "runtime.memory.heap.released_bytes",
		"/memory/classes/heap/stacks:bytes":   "runtime.memory.heap.stack_bytes",
		"/memory/classes/heap/unused:bytes":   "runtime.memory.heap.unused_bytes",
		"/memory/classes/total:bytes":         "runtime.memory.total_bytes",
		"/sched/goroutines:goroutines":        "runtime.scheduler.goroutines",
	}

	for name := range sampleMetrics {
		samples = append(samples, metrics.Sample{Name: name})
	}

	return &runtimeStatsHeartbeat{
		client:        client,
		samples:       samples,
		sampleMetrics: sampleMetrics,
	}
}

// Name returns a constant name for the runtime stats heartbeat.
func (r *runtimeStatsHeartbeat) Name() string {
	return "runtime_stats"
}

// Run emits gauges describing goroutines and memory usage.
func (r *runtimeStatsHeartbeat) Run() error {
	tags := map[string]interface{}{
		"heartbeat":  r.Name(),
		"go_version": runtime.Version(),
		"go_os":      runtime.GOOS,
		"go_arch":    runtime.GOARCH,
	}

	metrics.Read(r.samples)

	for _, sample := range r.samples {
		switch sample.Value.Kind() {
		case metrics.KindUint64:
			r.client.Gauge(
				r.sampleMetrics[sample.Name],
				float64(sample.Value.Uint64()),
				tags,
			)
		case metrics.KindFloat64:
			r.client.Gauge(
				r.sampleMetrics[sample.Name],
				sample.Value.Float64(),
				tags,
			)
		}
	}

	return nil
}

// resourceUsageHeartbeat reports process resource usage indicated by the getrusage syscall.
type resourceUsageHeartbeat struct {
	client aperture.Statsd
	usage  *syscall.Rusage
}

// NewResourceUsageHeartbeat creates a new resourceUsageHeartbeat.
func NewResourceUsageHeartbeat(client aperture.Statsd) Heartbeat {
	return &resourceUsageHeartbeat{
		client: client,
		usage:  new(syscall.Rusage),
	}
}

// Name returns a constant name for the resource usage heartbeat.
func (r *resourceUsageHeartbeat) Name() string {
	return "resource_usage"
}

// Run emits gauges describing current process resource usage.
func (r *resourceUsageHeartbeat) Run() error {
	tags := map[string]interface{}{
		"heartbeat":  r.Name(),
		"go_version": runtime.Version(),
		"go_os":      runtime.GOOS,
		"go_arch":    runtime.GOARCH,
	}

	if err := syscall.Getrusage(syscall.RUSAGE_SELF, r.usage); err != nil {
		return err
	}

	r.client.Gauge("resource.usage.cpu.user", float64(r.usage.Utime.Nano()), tags)
	r.client.Gauge("resource.usage.cpu.system", float64(r.usage.Stime.Nano()), tags)
	r.client.Gauge("resource.usage.memory.max_rss", float64(1024*r.usage.Maxrss), tags)
	r.client.Gauge("resource.usage.memory.page_faults.minor", float64(r.usage.Minflt), tags)
	r.client.Gauge("resource.usage.memory.page_faults.major", float64(r.usage.Majflt), tags)
	r.client.Gauge("resource.usage.io.reads", float64(r.usage.Inblock), tags)
	r.client.Gauge("resource.usage.io.writes", float64(r.usage.Oublock), tags)
	r.client.Gauge("resource.usage.load.context_switch.voluntary", float64(r.usage.Nvcsw), tags)
	r.client.Gauge("resource.usage.load.context_switch.involuntary", float64(r.usage.Nivcsw), tags)

	return nil
}
