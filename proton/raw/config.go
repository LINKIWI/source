package proton

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"lib.kevinlin.info/aperture"
	"lib.kevinlin.info/aperture/lib"
)

var (
	// defaultConfig is a Config instance with reasonable defaults.
	defaultConfig = &Config{
		ClientID:      "proton",
		ClientVersion: "none",
		Backend:       http.DefaultClient,
		Metrics:       lib.NewNoopStatsd(),
		Logger:        log.New(io.Discard, "", 0),
	}
)

// Config describes parameters for creating a Proton Supercharged client.
//
// Most clients should use the NewConfig function in order to assign proper defaults and perform
// automatic parameter validation:
//
//	cfg, err := proton.NewConfig(&proton.Config{BaseURL: ...})
//	if err != nil { ... }
type Config struct {
	// BaseURL points to the URL of the Supercharged HTTP server.
	BaseURL *url.URL

	// ClientID is an optional client identifier attached to all outgoing requests. When
	// omitted, a default Proton client identifier is used.
	ClientID string

	// ClientVersion is an optional string describing the client application version attached to
	// all outgoing requests.
	ClientVersion string

	// Backend is an *http.Client backend to use for making requests. When omitted, the default
	// globally shared net/http client is used.
	Backend *http.Client

	// Metrics is an optional aperture.Statsd implementation used for emitting metrics for all
	// outgoing requests. When omitted, instrumentation features are disabled.
	Metrics aperture.Statsd

	// Logger is an optional *log.Logger standard library logger for logging internal library
	// events, including request/response lifecycles.
	Logger *log.Logger
}

// NewConfig creates a Config with proper defaults assigned for omitted configuration parameters.
// An error is returned for any configuration validation errors.
func NewConfig(cfg *Config) (*Config, error) {
	if cfg.ClientID == "" {
		cfg.ClientID = defaultConfig.ClientID
	}
	if cfg.ClientVersion == "" {
		cfg.ClientVersion = defaultConfig.ClientVersion
	}
	if cfg.Backend == nil {
		cfg.Backend = defaultConfig.Backend
	}
	if cfg.Metrics == nil {
		cfg.Metrics = defaultConfig.Metrics
	}
	if cfg.Logger == nil {
		cfg.Logger = defaultConfig.Logger
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate returns an error for validation errors in the passed config.
func (c *Config) validate() error {
	if c.BaseURL == nil {
		return errors.New("proton: base URL is missing")
	} else if c.BaseURL.Scheme == "" || c.BaseURL.Host == "" {
		return errors.New("proton: base URL is malformed and missing fields")
	} else if !regexp.MustCompile("^[a-zA-Z0-9-_]*$").MatchString(c.ClientID) {
		return errors.New("proton: client ID uses illegal characters")
	}

	return nil
}
