package config

import (
	"time"
)

// Listener describes a network listener to pass to the service process.
type Listener struct {
	Name     string   `yaml:"name"`
	Address  *Address `yaml:"address"`
	SocketFD int      `yaml:"socket_fd"`
}

// Runtime describes service runtime parameters.
type Runtime struct {
	Name        string   `yaml:"name"`
	Directory   string   `yaml:"directory"`
	Path        string   `yaml:"path"`
	Args        []string `yaml:"args"`
	Environment []struct {
		Key   string `yaml:"key"`
		Value string `yaml:"value"`
	} `yaml:"environment"`
}

// Options describes parameters for operations related to process reloads.
type Options struct {
	ReloadSignal    *Signal       `yaml:"reload_signal"`
	ReloadUptime    time.Duration `yaml:"reload_uptime"`
	ShutdownSignal  *Signal       `yaml:"shutdown_signal"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

// Metrics describes a statsd server for application metrics.
type Metrics struct {
	Address string `yaml:"address"`
	Proxy   string `yaml:"proxy"`
	Prefix  string `yaml:"prefix"`
}

// Application describes application-level configuration.
type Application struct {
	SentryDSN string   `yaml:"sentry_dsn"`
	Metrics   *Metrics `yaml:"metrics"`
}

// Service describes configuration for the underlying managed service.
type Service struct {
	Listeners []Listener `yaml:"listeners"`
	Runtime   Runtime    `yaml:"runtime"`
	Options   Options    `yaml:"options"`
}

// Config is the top-level application configuration.
type Config struct {
	Application Application `yaml:"application"`
	Service     Service     `yaml:"service"`
}
