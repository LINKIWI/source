package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/getsentry/sentry-go"
	"lib.kevinlin.info/aperture"
	"lib.kevinlin.info/aperture/lib"

	"bitflow/internal/config"
	"bitflow/internal/log"
	"bitflow/internal/meta"
	"bitflow/internal/server"
)

var (
	flagConfig = flag.String("config", os.Getenv("BITFLOW_CONFIG"), "path to the Bitflow config file")
)

func init() {
	flag.Parse()
}

func main() {
	var err error
	var metrics aperture.Statsd = lib.NewNoopStatsd()

	cfg, err := config.New(*flagConfig)
	if err != nil {
		panic(err)
	}

	logCtx, err := log.Init(cfg.Application.Log)
	if err != nil {
		panic(err)
	}

	if cfg.Application.Metrics != nil {
		logCtx.Debug.Printf(
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
		logCtx.Debug.Printf(
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

		defer sentry.Recover()
	}

	logCtx.Info.Printf(
		"main: starting Bitflow proxy server: config=%s version=%s proxies=%d",
		*flagConfig,
		meta.Version,
		len(cfg.Server.Proxies),
	)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	srv := &server.Server{
		Config:  cfg.Server,
		Metrics: metrics,
		Log:     logCtx,
	}

	go func() {
		<-shutdown
		go srv.Close()

		<-shutdown
		logCtx.Warn.Printf(
			"main: received interrupt signal during graceful shutdown phase; " +
				"forcing immediate ungraceful shutdown",
		)
		os.Exit(1)
	}()

	if err := srv.Serve(); err != nil {
		panic(err)
	}

	logCtx.Info.Printf("main: completed graceful shutdown: proxies=%d", len(cfg.Server.Proxies))
}
