package config

import (
	"time"
)

var (
	connectionLoadBalancers = map[string]bool{
		"":            true,
		"none":        true,
		"failover":    true,
		"round-robin": true,
		"sni":         true,
	}
	tlsVerifyModes = map[string]bool{
		"":        true,
		"none":    true,
		"relaxed": true,
		"strict":  true,
	}
)

// Target describes the target TCP server.
type Target struct {
	Name              string        `yaml:"name"`
	Address           *Address      `yaml:"address"`
	Proxy             *Address      `yaml:"proxy"`
	ConnectAttempts   int           `yaml:"connect_attempts"`
	ConnectBackoff    time.Duration `yaml:"connect_backoff"`
	ConnectTimeout    time.Duration `yaml:"connect_timeout"`
	ResolveTimeout    time.Duration `yaml:"resolve_timeout"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	KeepaliveInterval time.Duration `yaml:"keepalive_interval"`
	TLSContext        *TLSContext   `yaml:"tls_context"`
}

// TLSCertificate describes a TLS keypair with an optional certificate authority, all in PEM format.
type TLSCertificate struct {
	Key         string `yaml:"key"`
	Certificate string `yaml:"certificate"`
	Authority   string `yaml:"authority"`
}

// TLSContext describes TLS termination properties for the local server listener.
type TLSContext struct {
	Certificates         []TLSCertificate `yaml:"certificates"`
	ApplicationProtocols []string         `yaml:"application_protocols"`
	VerifyMode           string           `yaml:"verify_mode"`
}

// Listener describes the local server listener.
type Listener struct {
	Address           *Address      `yaml:"address"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	AuthorizedSources []CIDR        `yaml:"authorized_sources"`
	TLSContext        *TLSContext   `yaml:"tls_context"`
}

// Options describes available options for one proxy instance.
type Options struct {
	ConnectionLimit        int           `yaml:"connection_limit"`
	ConnectionLifetime     time.Duration `yaml:"connection_lifetime"`
	ConnectionLog          string        `yaml:"connection_log"`
	ConnectionLoadBalancer string        `yaml:"connection_load_balancer"`
	EnableProxyHeader      bool          `yaml:"enable_proxy_header"`
	EnableDualStack        bool          `yaml:"enable_dual_stack"`
}

// Proxy describes a single TCP proxy.
type Proxy struct {
	Name     string    `yaml:"name"`
	Listener *Listener `yaml:"listener"`
	Targets  []*Target `yaml:"targets"`
	Options  Options   `yaml:"options"`
}

// Metrics describes metrics client configuration.
type Metrics struct {
	Address string `yaml:"address"`
	Proxy   string `yaml:"proxy"`
	Prefix  string `yaml:"prefix"`
}

// Log describes logging configuration.
type Log struct {
	Level string `yaml:"level"`
}

// Application describes application configuration.
type Application struct {
	SentryDSN string   `yaml:"sentry_dsn"`
	Metrics   *Metrics `yaml:"metrics"`
	Log       *Log     `yaml:"log"`
}

// Server describes server configuration.
type Server struct {
	Proxies []Proxy `yaml:"proxies"`
}

// Config is the top-level Bitflow configuration struct.
type Config struct {
	Application Application `yaml:"application"`
	Server      Server      `yaml:"server"`
}
