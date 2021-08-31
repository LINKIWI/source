package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"lib.kevinlin.info/aperture"
	"lib.kevinlin.info/aperture/lib"

	"zero/internal/config"
	"zero/internal/meta"
	"zero/pkg/manager"
)

var (
	flagConfig = flag.String("config", os.Getenv("ZERO_CONFIG"), "path to the Zero config file")
)

func init() {
	flag.Parse()

	log.SetFlags(log.Lshortfile | log.Lmsgprefix | log.Ldate | log.Ltime)
	log.SetPrefix("[zero] ")
}

func main() {
	var metrics aperture.Statsd = lib.NewNoopStatsd()
	var exit sync.WaitGroup

	log.Printf(
		"main: starting zero process manager: config=%s version=%s",
		*flagConfig,
		meta.Version,
	)

	cfg, err := config.New(*flagConfig)
	if err != nil {
		panic(err)
	}

	if cfg.Application.Metrics != nil {
		log.Printf(
			"main: enabling metrics reporting: address=%s prefix=%s",
			cfg.Application.Metrics.Address,
			cfg.Application.Metrics.Prefix,
		)

		metrics, err = aperture.NewClient(&aperture.Config{
			Address:                     cfg.Application.Metrics.Address,
			Prefix:                      cfg.Application.Metrics.Prefix,
			Proxy:                       cfg.Application.Metrics.Proxy,
			TransportProbeInterval:      10 * time.Second,
			LazyTransportInitialization: true,
			DefaultTags: map[string]interface{}{
				"version": meta.Version,
			},
		})
		if err != nil {
			panic(err)
		}

		metrics = lib.NewAsyncStatsd(metrics)
	}

	if cfg.Application.SentryDSN != "" {
		log.Printf(
			"main: enabling error reporting: sentry_dsn=%s",
			cfg.Application.SentryDSN,
		)

		opts := sentry.ClientOptions{
			Dsn:              cfg.Application.SentryDSN,
			Release:          meta.Version,
			Transport:        sentry.NewHTTPSyncTransport(),
			AttachStacktrace: true,
		}

		if err := sentry.Init(opts); err != nil {
			panic(err)
		}

		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("name", cfg.Service.Runtime.Name)
			scope.SetTag("directory", cfg.Service.Runtime.Directory)
			scope.SetTag("path", cfg.Service.Runtime.Path)
		})

		defer sentry.Recover()
	}

	for _, listener := range cfg.Service.Listeners {
		log.Printf(
			"main: creating process listener: name=%s address=%v fd=%d",
			listener.Name,
			listener.Address,
			listener.SocketFD,
		)
	}

	tags := map[string]interface{}{
		"name": cfg.Service.Runtime.Name,
		"path": cfg.Service.Runtime.Path,
	}

	m, err := manager.NewManager(cfg.Service)
	if err != nil {
		panic(err)
	}

	reload := make(chan os.Signal, 1)
	signal.Notify(reload, cfg.Service.Options.ReloadSignal.Signal)

	go func() {
		var reloads int

		for {
			<-reload

			log.Printf(
				"main: processing zero process reload signal: signal=%s",
				cfg.Service.Options.ReloadSignal,
			)

			reloads++
			stopwatch := lib.NewStopwatch()
			metrics.Incr("process.reload", tags)
			metrics.Gauge("process.reload.total", float64(reloads), tags)

			stats, err := m.Reload()
			if err != nil {
				log.Printf("main: error during process reload: err=%v", err)
				sentry.CaptureException(err)
			} else {
				log.Println("main: completed process reload successfully")
			}

			metrics.Timing("process.reload.duration", stopwatch.Elapsed(), tags)

			if stats != nil {
				rusage := stats.SysUsage().(*syscall.Rusage)

				metrics.Gauge("process.stats.system_time", stats.SystemTime().Seconds(), tags)
				metrics.Gauge("process.stats.user_time", stats.UserTime().Seconds(), tags)
				metrics.Gauge("process.stats.max_rss", float64(1024*rusage.Maxrss), tags)
				metrics.Gauge("process.stats.minor_faults", float64(rusage.Minflt), tags)
				metrics.Gauge("process.stats.major_faults", float64(rusage.Majflt), tags)
				metrics.Gauge("process.stats.disk_reads", float64(rusage.Inblock), tags)
				metrics.Gauge("process.stats.disk_writes", float64(rusage.Oublock), tags)
				metrics.Gauge("process.stats.signals", float64(rusage.Nsignals), tags)
				metrics.Gauge("process.stats.voluntary_context_switches", float64(rusage.Nvcsw), tags)
				metrics.Gauge("process.stats.involuntary_context_switches", float64(rusage.Nivcsw), tags)
			}
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT)
	signal.Notify(shutdown, syscall.SIGTERM)
	exit.Add(1)

	go func() {
		defer exit.Done()

		<-shutdown
		log.Println("main: starting zero graceful shutdown")

		stopwatch := lib.NewStopwatch()
		metrics.Incr("process.shutdown", tags)

		if err := m.Close(); err != nil {
			log.Printf("main: error during graceful shutdown of child: err=%v", err)
			sentry.CaptureException(err)
		} else {
			log.Println("main: completed graceful shutdown successfully")
		}

		metrics.Timing("process.shutdown.duration", stopwatch.Elapsed(), tags)
	}()

	log.Printf(
		"main: starting service process: directory=%s path=%s args=%v",
		cfg.Service.Runtime.Directory,
		cfg.Service.Runtime.Path,
		cfg.Service.Runtime.Args,
	)

	if err := m.Start(); err != nil {
		panic(err)
	}

	metrics.Incr("process.start", tags)

	go func() {
		uptime := lib.NewStopwatch()
		tags := map[string]interface{}{
			"name": cfg.Service.Runtime.Name,
			"path": cfg.Service.Runtime.Path,
		}

		for {
			if process, err := m.Process(); err == nil {
				tags["pid"] = process.Pid
				metrics.Gauge("manager.uptime", uptime.Elapsed().Seconds(), tags)
			}

			time.Sleep(10 * time.Second)
		}
	}()

	log.Printf(
		"main: started service process; waiting for reload or shutdown: "+
			"reload=%v uptime=%v shutdown=%v timeout=%v",
		cfg.Service.Options.ReloadSignal,
		cfg.Service.Options.ReloadUptime,
		cfg.Service.Options.ShutdownSignal,
		cfg.Service.Options.ShutdownTimeout,
	)

	exit.Wait()
}
