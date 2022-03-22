package config

import (
	"io/fs"
	"net/http"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	filterTypes = map[string]bool{
		"authentication":  true,
		"body":            true,
		"compression":     true,
		"concurrency":     true,
		"header":          true,
		"identity":        true,
		"instrumentation": true,
		"ip":              true,
		"log":             true,
		"metadata":        true,
		"redirect":        true,
		"rewrite":         true,
		"size":            true,
		"static":          true,
	}
	metricSerializers = map[string]bool{
		"":       true,
		"statsd": true,
		"influx": true,
	}
	upstreamTransports = map[string]bool{
		"":      true,
		"http":  true,
		"http2": true,
	}
	tlsVerifyModes = map[string]bool{
		"":        true,
		"none":    true,
		"relaxed": true,
		"strict":  true,
	}
	tlsProtocolVersions = map[string]bool{
		"":         true,
		"lenient":  true,
		"standard": true,
		"modern":   true,
	}
	connectionProtocols = map[string]bool{
		"":        true,
		"haproxy": true,
	}
	healthCheckMethods = map[string]bool{
		"":                 true,
		http.MethodGet:     true,
		http.MethodHead:    true,
		http.MethodPost:    true,
		http.MethodPut:     true,
		http.MethodPatch:   true,
		http.MethodDelete:  true,
		http.MethodConnect: true,
		http.MethodOptions: true,
		http.MethodTrace:   true,
	}
)

// Filter describes universal configuration for a request filter.
type Filter struct {
	Name   string    `yaml:"name"`
	Type   string    `yaml:"type"`
	Params yaml.Node `yaml:"params"`
}

// UpstreamConnection describes configuration for an upstream connection from the server.
type UpstreamConnection struct {
	ConnectTimeout        time.Duration `yaml:"connect_timeout"`
	ReadTimeout           time.Duration `yaml:"read_timeout"`
	ReadHeaderTimeout     time.Duration `yaml:"read_header_timeout"`
	WriteTimeout          time.Duration `yaml:"write_timeout"`
	ExpectContinueTimeout time.Duration `yaml:"expect_continue_timeout"`
	IdleTimeout           time.Duration `yaml:"idle_timeout"`
	KeepaliveInterval     time.Duration `yaml:"keepalive_interval"`
	IdleLimit             int           `yaml:"idle_limit"`
	ReadBufferSize        int           `yaml:"read_buffer_size"`
	WriteBufferSize       int           `yaml:"write_buffer_size"`
}

// ServerConnection describes configuration for a downstream connection from a client.
type ServerConnection struct {
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
	ActiveLimit       int           `yaml:"active_limit"`
}

// HealthCheck describes parameters for an upstream health check.
type HealthCheck struct {
	Host    string `yaml:"host"`
	Method  string `yaml:"method"`
	Path    string `yaml:"path"`
	Headers []struct {
		Key   string `yaml:"key"`
		Value string `yaml:"value"`
	} `yaml:"headers"`
	Interval           time.Duration `yaml:"interval"`
	Jitter             time.Duration `yaml:"jitter"`
	EnableCircuitBreak bool          `yaml:"enable_circuit_break"`
	Status             *struct {
		Lower int `yaml:"lower"`
		Upper int `yaml:"upper"`
	} `yaml:"status"`
}

// Upstream describes one upstream HTTP server.
type Upstream struct {
	Name        string             `yaml:"name"`
	Address     *Address           `yaml:"address"`
	Proxy       *URL               `yaml:"proxy"`
	Transport   string             `yaml:"transport"`
	Connection  UpstreamConnection `yaml:"connection"`
	TLSContext  *TLSContext        `yaml:"tls_context"`
	HealthCheck *HealthCheck       `yaml:"health_check"`
}

// Route describes a routing policy for a request to an available upstream.
type Route struct {
	Upstream string `yaml:"upstream"`
	Match    struct {
		Host   *Regex `yaml:"host"`
		Method *Regex `yaml:"method"`
		Path   *Regex `yaml:"path"`
		Header *struct {
			Key   string `yaml:"key"`
			Value *Regex `yaml:"value"`
		} `yaml:"header"`
	} `yaml:"match"`
}

// VirtualHost describes one virtual host served by the server.
type VirtualHost struct {
	Name      string      `yaml:"name"`
	Host      *Regex      `yaml:"host"`
	Upstreams []*Upstream `yaml:"upstreams"`
	Routes    []*Route    `yaml:"routes"`
	Filters   []*Filter   `yaml:"filters"`
}

// TLSCertificate describes a the paths to a single TLS server key and certificate.
type TLSCertificate struct {
	Key         string `yaml:"key"`
	Certificate string `yaml:"certificate"`
	Authority   string `yaml:"authority"`
	Peer        *Regex `yaml:"peer"`
}

// TLSContext describes server and client TLS configuration.
type TLSContext struct {
	Certificates     []*TLSCertificate `yaml:"certificates"`
	VerifyMode       string            `yaml:"verify_mode"`
	VerifyPeer       *Regex            `yaml:"verify_peer"`
	HandshakeTimeout time.Duration     `yaml:"handshake_timeout"`
	ReloadInterval   time.Duration     `yaml:"reload_interval"`
	ProtocolVersion  string            `yaml:"protocol_version"`
}

// Authorization describes listener connection authorization of clients.
type Authorization struct {
	Sources        []CIDR      `yaml:"sources"`
	SocketFileMode fs.FileMode `yaml:"socket_file_mode"`
}

// Listener describes a single bound listener for the server.
type Listener struct {
	Name          string           `yaml:"name"`
	Address       *Address         `yaml:"address"`
	Protocol      string           `yaml:"protocol"`
	Connection    ServerConnection `yaml:"connection"`
	TLSContext    *TLSContext      `yaml:"tls_context"`
	Authorization Authorization    `yaml:"authorization"`
}

// Proxy describes proxy behavior options.
type Proxy struct {
	EnableErrorPropagation bool `yaml:"enable_error_propagation"`
	EnableErrorSemantics   bool `yaml:"enable_error_semantics"`
}

// Server describes configuration for the server.
type Server struct {
	Listeners    []*Listener    `yaml:"listeners"`
	VirtualHosts []*VirtualHost `yaml:"virtual_hosts"`
	Proxy        Proxy          `yaml:"proxy"`
}

// Metrics describes metrics output configuration.
type Metrics struct {
	Address    string `yaml:"address"`
	Prefix     string `yaml:"prefix"`
	Proxy      string `yaml:"proxy"`
	Serializer string `yaml:"serializer"`
}

// Log describes logging configuration.
type Log struct {
	Level   string   `yaml:"level"`
	Outputs []string `yaml:"outputs"`
}

// Application describes Courier-specific application configuration.
type Application struct {
	Instance  string   `yaml:"instance"`
	Metrics   *Metrics `yaml:"metrics"`
	Log       *Log     `yaml:"log"`
	SentryDSN *URL     `yaml:"sentry_dsn"`
}

// Config is the top-level application config.
type Config struct {
	Application *Application `yaml:"application"`
	Server      *Server      `yaml:"server"`
}
