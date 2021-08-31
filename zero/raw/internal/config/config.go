package config

import (
	"fmt"
	"os"
	"regexp"
	"syscall"

	"gopkg.in/yaml.v3"
)

// New creates a configuration object from a path on disk.
func New(path string) (*Config, error) {
	var cfg *Config

	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf(
			"config: error opening config file on disk: path=%s err=%v",
			path,
			err,
		)
	}

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf(
			"config: error deserializing config file contents: err=%v",
			err,
		)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config: error validating config: err=%v", err)
	}

	return cfg, nil
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if len(c.Service.Listeners) == 0 {
		return fmt.Errorf("config: no listeners specified")
	}

	for _, listener := range c.Service.Listeners {
		if listener.Name == "" {
			return fmt.Errorf("config: listener name missing")
		} else if listener.Address == nil {
			return fmt.Errorf("config: listener address missing")
		} else if listener.SocketFD < 3 {
			return fmt.Errorf(
				"config: listener socket file descriptor must be >= 3: fd=%d",
				listener.SocketFD,
			)
		}
	}

	if c.Service.Runtime.Name == "" {
		return fmt.Errorf("config: runtime name identifier missing")
	} else if c.Service.Runtime.Path == "" {
		return fmt.Errorf("config: runtime binary path missing")
	}

	for _, env := range c.Service.Runtime.Environment {
		if !regexp.MustCompile("[a-zA-Z_]+[a-zA-Z0-9_]*").MatchString(env.Key) {
			return fmt.Errorf("config: environment variable key is missing or invalid")
		}
	}

	if c.Service.Options.ReloadSignal == nil {
		return fmt.Errorf("config: zero process reload signal missing")
	} else if c.Service.Options.ReloadSignal.Signal == syscall.SIGINT {
		return fmt.Errorf("config: SIGINT is not permitted as the service reload signal")
	} else if c.Service.Options.ReloadSignal.Signal == syscall.SIGTERM {
		return fmt.Errorf("config: SIGTERM is not permitted as the service reload signal")
	} else if c.Service.Options.ShutdownSignal == nil {
		return fmt.Errorf("config: service process reload signal missing")
	}

	return nil
}
