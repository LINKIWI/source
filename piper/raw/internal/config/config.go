package config

import (
	"fmt"
	"net"
	"os"

	"github.com/BurntSushi/toml"

	"piper/internal/meta"
)

// New creates and validates a Config instance from a file path on disk.
func New(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg *Config
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the parsed configuration struct and returns a non-nil error on validation
// errors.
func (c *Config) Validate() error {
	for _, relay := range c.Relays {
		if relay.Name == "" {
			return fmt.Errorf("config: relay name is not populated")
		} else if relay.LogFile.Pattern == "" {
			return fmt.Errorf("config: relay log file path is not populated")
		} else if len(relay.Delimiter) > 1 {
			return fmt.Errorf("config: relay delimiter must be omitted or exactly one character")
		} else if relay.BufferLength < 0 {
			return fmt.Errorf("config: buffer length must be a nonnegative integer")
		} else if _, _, err := net.SplitHostPort(relay.ProxyAddress); relay.ProxyAddress != "" && err != nil {
			return fmt.Errorf("config: failed to parse proxy server address: err=%v", err)
		} else if _, _, err := net.SplitHostPort(relay.KafkaAddress); err != nil {
			return fmt.Errorf("config: failed to parse Kafka address: err=%v", err)
		}

		switch relay.Serializer {
		case "", SerializerLine, SerializerJSON:
		default:
			return fmt.Errorf("config: unknown relay serializer identifier")
		}

		switch relay.SeekPosition {
		case "", SeekPositionStart, SeekPositionEnd:
		default:
			return fmt.Errorf("config: unknown seek position identifier")
		}

		for _, tagID := range relay.TagIdentifiers {
			switch tagID {
			case TagHostname, TagVersion:
			default:
				return fmt.Errorf("config: unknown relay tag identifier: tag=%s", tagID)
			}
		}
	}

	return nil
}

// DefaultMetricsTags returns a mapping of static tags to use as the default tags set included with
// all outgoing emitted metrics.
func (p *Piper) DefaultMetricsTags() (map[string]interface{}, error) {
	tags := map[string]interface{}{
		"version": meta.Version,
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	tags["host"] = hostname

	return tags, nil
}

// Tags returns a mapping of derived static tags associated with the relay based on the list of
// enabled tag identifiers. A non-nil error is returned if an error is encountered while computing
// any tag value.
func (r *Relay) Tags() (map[string]string, error) {
	tags := make(map[string]string)

	for _, tagID := range r.TagIdentifiers {
		switch tagID {
		case "hostname":
			hostname, err := os.Hostname()
			if err != nil {
				return nil, err
			}
			tags["hostname"] = hostname

		case "version":
			tags["version"] = meta.Version
		}
	}

	return tags, nil
}
