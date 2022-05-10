package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"unistore/internal/config"
	"unistore/internal/log"
	"unistore/internal/meta"
	"unistore/internal/metrics"
	"unistore/internal/server"
)

// validateCmd validates args and flags before running any command.
func validateCmd(cmd *cobra.Command, args []string) error {
	flagConfig, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("main: error reading command line flags: err=%v", err)
	}

	if flagConfig == "" {
		return fmt.Errorf("main: required configuration file not specified")
	}

	return nil
}

// runServer runs the main Unistore server.
func runServer(cmd *cobra.Command, args []string) error {
	flagConfig, _ := cmd.Flags().GetString("config")

	cfg, err := config.New(flagConfig)
	if err != nil {
		return fmt.Errorf("main: error creating configuration: err=%v", err)
	}

	if cfg.Application.Log != nil {
		if err := log.Init(cfg.Application.Log); err != nil {
			panic(fmt.Errorf("main: error initializing logging: err=%v", err))
		}
	}

	zap.L().Info(
		"initializing unistore server",
		zap.String("config", flagConfig),
		zap.String("version", meta.Version),
	)

	if cfg.Application.Metrics != nil {
		zap.L().Info(
			"enabling metrics reporting",
			zap.String("address", cfg.Application.Metrics.Address),
			zap.String("proxy", cfg.Application.Metrics.Proxy),
			zap.String("prefix", cfg.Application.Metrics.Prefix),
		)

		if err := metrics.Init(cfg.Application.Metrics); err != nil {
			zap.L().Error("failed to create metrics client", zap.Error(err))
			return err
		}
	}

	if cfg.Application.Errors != nil {
		zap.L().Info(
			"enabling error reporting",
			zap.String("sentry_dsn", cfg.Application.Errors.SentryDSN),
		)

		opts := sentry.ClientOptions{
			Dsn:              cfg.Application.Errors.SentryDSN,
			Release:          meta.Version,
			AttachStacktrace: true,
		}

		if err := sentry.Init(opts); err != nil {
			zap.L().Error("failed to create Sentry client", zap.Error(err))
			return err
		}

		defer sentry.Recover()
	}

	srv, err := server.New(cfg.Server)
	if err != nil {
		zap.L().Error("error during server initialization", zap.Error(err))
		return err
	}

	closed := make(chan bool)

	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

		sig := <-shutdown
		zap.L().Info("initiating graceful server shutdown", zap.Stringer("signal", sig))

		go func() {
			sig = <-shutdown
			zap.L().Warn(
				"received shutdown signal during graceful shutdown phase; "+
					"forcing immediate ungraceful shutdown",
				zap.Stringer("signal", sig),
			)
			os.Exit(1)
		}()

		if err := srv.Close(); err != nil {
			zap.L().Error("error during server close", zap.Error(err))
		}

		if err := metrics.Client.Close(); err != nil {
			zap.L().Error("error during metrics client close", zap.Error(err))
		}

		closed <- true
	}()

	if err := srv.Serve(context.Background()); err != nil {
		zap.L().Error("error during serve", zap.Error(err))
		return err
	}

	// Additionally wait for the server to complete (or at least attempt) its full close
	// routine, even if the server's serving path has already exited.
	<-closed

	zap.L().Info("completed graceful shutdown")

	return nil
}

// runConfigShow reads, parses, and prints configuration.
func runConfigShow(cmd *cobra.Command, args []string) error {
	flagConfig, _ := cmd.Flags().GetString("config")

	cfg, err := config.New(flagConfig)
	if err != nil {
		return fmt.Errorf("main: error creating configuration: err=%v", err)
	}

	fmt.Print(cfg)

	return nil
}

// runConfigValidate validates the specified configuration and returns an error on validation
// failures.
func runConfigValidate(cmd *cobra.Command, args []string) error {
	flagConfig, _ := cmd.Flags().GetString("config")

	_, err := config.New(flagConfig)
	if err != nil {
		return fmt.Errorf("main: error creating configuration: err=%v", err)
	}

	return nil
}
