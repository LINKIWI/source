package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"

	"courier/internal/config"
	"courier/internal/log"
	"courier/internal/meta"
	"courier/internal/metrics"
	"courier/internal/server"
)

var (
	cfg *config.Config
)

var (
	flagConfig   = flag.String("config", os.Getenv("COURIER_CONFIG"), "path to the courier config file")
	flagValidate = flag.Bool("validate", false, "validate the configuration file and exit")
)

func init() {
	var err error

	flag.Parse()

	cfg, err = config.New(*flagConfig)
	if err != nil {
		panic(err)
	}

	if cfg.Application.Log != nil {
		if err := log.Init(cfg); err != nil {
			panic(err)
		}
	}
}

func main() {
	zap.L().Info(
		"initializing courier server",
		zap.String("version", meta.Version),
		zap.String("config", *flagConfig),
		zap.String("hash", cfg.Hash()),
	)

	for _, listener := range cfg.Server.Listeners {
		zap.L().Info(
			"registered server listener",
			zap.String("name", listener.Name),
			zap.String("address", listener.Address.String()),
			zap.Any("connection", listener.Connection),
			zap.Any("tls_context", listener.TLSContext),
		)
	}
	for _, vhost := range cfg.Server.VirtualHosts {
		zap.L().Info(
			"registered virtual host",
			zap.String("name", vhost.Name),
			zap.String("host", vhost.Host.String()),
			zap.Int("upstreams", len(vhost.Upstreams)),
			zap.Int("routes", len(vhost.Routes)),
		)

		for _, upstream := range vhost.Upstreams {
			zap.L().Debug(
				"registered virtual host upstream",
				zap.String("vhost", vhost.Name),
				zap.String("name", upstream.Name),
				zap.Stringer("address", upstream.Address),
				zap.Stringer("proxy", upstream.Proxy),
				zap.String("transport", upstream.Transport),
				zap.Any("connection", upstream.Connection),
				zap.Any("tls_context", upstream.TLSContext),
			)
		}

		for _, route := range vhost.Routes {
			zap.L().Debug(
				"registered virtual host route",
				zap.String("vhost", vhost.Name),
				zap.String("upstream", route.Upstream),
				zap.Any("match", route.Match),
			)
		}
	}

	if *flagValidate {
		zap.L().Info("config file successfully validated", zap.String("config", *flagConfig))
		return
	}

	if cfg.Application.Metrics != nil {
		zap.L().Info(
			"initializing metrics subsystem",
			zap.String("address", cfg.Application.Metrics.Address),
			zap.String("prefix", cfg.Application.Metrics.Prefix),
			zap.String("serializer", cfg.Application.Metrics.Serializer),
		)

		if err := metrics.Init(cfg); err != nil {
			zap.L().Fatal("failed to initialize metrics subsystem", zap.Error(err))
		}
	}

	if cfg.Application.SentryDSN != nil {
		zap.L().Info(
			"enabling error reporting",
			zap.String("sentry_dsn", cfg.Application.SentryDSN.String()),
		)

		opts := sentry.ClientOptions{
			Dsn:              cfg.Application.SentryDSN.String(),
			Release:          meta.Version,
			AttachStacktrace: true,
		}

		if err := sentry.Init(opts); err != nil {
			zap.L().Fatal("failed to create Sentry client", zap.Error(err))
		}

		defer sentry.Recover()
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT)
	signal.Notify(shutdown, syscall.SIGTERM)

	zap.L().Info("initializing HTTP server")
	srv, err := server.New(cfg)
	if err != nil {
		zap.L().Fatal("failed to start server", zap.Error(err))
	}

	go func() {
		<-shutdown
		zap.L().Info("initiating graceful server shutdown")
		go srv.Close()

		<-shutdown
		zap.L().Warn(
			"received interrupt signal during graceful shutdown phase; " +
				"forcing immediate ungraceful shutdown",
		)
		os.Exit(1)
	}()

	zap.L().Info("starting indefinite serving")
	if err := srv.Serve(); err != nil {
		zap.L().Fatal("failed to serve", zap.Error(err))
	}

	zap.L().Info("completed graceful server shutdown")
}
