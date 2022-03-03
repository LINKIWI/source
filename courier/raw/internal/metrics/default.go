package metrics

import (
	"runtime"
	"runtime/metrics"
	"syscall"
	"time"

	"lib.kevinlin.info/aperture"

	"courier/internal/config"
)

// initDefaultHeartbeats registers a set of default heartbeats.
func initDefaultHeartbeats(cfg *config.Config) {
	RegisterHeartbeat("uptime", DefaultHeartbeatInterval, newUptimeHeartbeat())
	RegisterHeartbeat("config", DefaultHeartbeatInterval, newConfigHeartbeat(cfg))
	RegisterHeartbeat("runtime_stats", DefaultHeartbeatInterval, newRuntimeStatsHeartbeat())
	RegisterHeartbeat("resource_usage", DefaultHeartbeatInterval, newResourceUsageHeartbeat())
}

// statelessHeartbeat is a Heartbeat wrapper for a stateless callback function.
type statelessHeartbeat struct {
	run func(client aperture.Statsd) error
}

// Run invokes the supplied callback.
func (s *statelessHeartbeat) Run(client aperture.Statsd) error {
	return s.run(client)
}

// uptimeHeartbeat tracks the uptime of the application.
type uptimeHeartbeat struct {
	start time.Time
}

// newUptimeHeartbeat creates a new uptimeHeartbeat.
func newUptimeHeartbeat() Heartbeat {
	return &uptimeHeartbeat{start: time.Now()}
}

// Run emits a gauge for the amount of time since the application started.
func (u *uptimeHeartbeat) Run(client aperture.Statsd) error {
	client.Gauge("uptime.seconds", time.Since(u.start).Seconds(), nil)
	return nil
}

// configHeartbeat reports metrics about the currently loaded application configuration.
type configHeartbeat struct {
	cfg *config.Config
}

// newConfigHeartbeat creates a new configHeartbeat.
func newConfigHeartbeat(cfg *config.Config) Heartbeat {
	return &configHeartbeat{cfg}
}

// Run emits a gauge for the current configuration version, identified by its unique hash.
func (c *configHeartbeat) Run(client aperture.Statsd) error {
	client.Gauge("config.version", 1, map[string]interface{}{"hash": c.cfg.Hash()})
	return nil
}

// runtimeStatsHeartbeat reports Go runtime stats.
type runtimeStatsHeartbeat struct {
	samples       []metrics.Sample
	sampleMetrics map[string]string
}

// newRuntimeStatsHeartbeat creates a new runtimeStatsHeartbeat with a set of default runtime stats.
func newRuntimeStatsHeartbeat() Heartbeat {
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
		samples:       samples,
		sampleMetrics: sampleMetrics,
	}
}

// Run emits gauges describing goroutines and memory usage.
func (r *runtimeStatsHeartbeat) Run(client aperture.Statsd) error {
	tags := map[string]interface{}{"go_version": runtime.Version()}

	metrics.Read(r.samples)

	for _, sample := range r.samples {
		switch sample.Value.Kind() {
		case metrics.KindUint64:
			client.Gauge(r.sampleMetrics[sample.Name], float64(sample.Value.Uint64()), tags)
		case metrics.KindFloat64:
			client.Gauge(r.sampleMetrics[sample.Name], sample.Value.Float64(), tags)
		}
	}

	return nil
}

// resourceUsageHeartbeat reports process resource usage indicated by the getrusage syscall.
type resourceUsageHeartbeat struct {
	usage *syscall.Rusage
}

// newResourceUsageHeartbeat creates a new resourceUsageHeartbeat.
func newResourceUsageHeartbeat() Heartbeat {
	return &resourceUsageHeartbeat{usage: new(syscall.Rusage)}
}

// Run emits gauges describing current process resource usage.
func (r *resourceUsageHeartbeat) Run(client aperture.Statsd) error {
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, r.usage); err != nil {
		return err
	}

	client.Gauge("resource.usage.cpu.user", float64(r.usage.Utime.Nano()), nil)
	client.Gauge("resource.usage.cpu.system", float64(r.usage.Stime.Nano()), nil)
	client.Gauge("resource.usage.memory.max_rss", float64(1024*r.usage.Maxrss), nil)
	client.Gauge("resource.usage.memory.page_faults.minor", float64(r.usage.Minflt), nil)
	client.Gauge("resource.usage.memory.page_faults.major", float64(r.usage.Majflt), nil)
	client.Gauge("resource.usage.io.reads", float64(r.usage.Inblock), nil)
	client.Gauge("resource.usage.io.writes", float64(r.usage.Oublock), nil)
	client.Gauge("resource.usage.load.context_switch.voluntary", float64(r.usage.Nvcsw), nil)
	client.Gauge("resource.usage.load.context_switch.involuntary", float64(r.usage.Nivcsw), nil)

	return nil
}
