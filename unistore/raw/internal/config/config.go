package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"unistore/schemas/common"
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

	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf(
			"config: error deserializing config file contents: path=%s err=%v",
			path,
			err,
		)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation error: err=%v", err)
	}

	return cfg, nil
}

// String dumps the parsed configuration as formatted YAML.
func (c *Config) String() string {
	marshaled, err := yaml.Marshal(c)
	if err != nil {
		return "<nil>"
	}

	return string(marshaled)
}

// validate validates the parsed configuration and returns an error on any validation failures.
func (c *Config) validate() error {
	if c.Server == nil && c.Client == nil {
		return fmt.Errorf("config: neither server nor client block is defined")
	} else if c.Server != nil && c.Client != nil {
		return fmt.Errorf("config: exactly one of server or client blocks can be defined")
	}

	if c.Server != nil {
		if c.Server.Listener.Address == nil {
			return fmt.Errorf("config: listener address not defined")
		}

		if composite := c.Server.Storage.Composite; composite != nil {
			if len(composite.Backends) == 0 {
				return fmt.Errorf("config: composite backend must have at least one backend")
			} else if _, ok := common.Dispatch_value[strings.ToUpper(composite.ReadDispatch)]; !ok {
				return fmt.Errorf(
					"config: unknown read dispatch policy: policy=%s",
					composite.ReadDispatch,
				)
			} else if _, ok := common.Dispatch_value[strings.ToUpper(composite.WriteDispatch)]; !ok {
				return fmt.Errorf(
					"config: unknown write dispatch policy: policy=%s",
					composite.WriteDispatch,
				)
			}

			for _, backend := range composite.Backends {
				bid, ok := common.Backend_value[strings.ToUpper(backend)]
				if !ok {
					return fmt.Errorf(
						"config: unknown identifier in composite backend children: backend=%s candidates=%v",
						backend,
						common.Backend_value,
					)
				}

				switch common.Backend(bid) {
				case common.Backend_COMPOSITE:
					return fmt.Errorf(
						"config: composite backend cannot depend on itself: backend=%s",
						common.Backend(bid),
					)
				case common.Backend_DISK:
					if c.Server.Storage.Disk == nil {
						return fmt.Errorf(
							"config: composite backend depends on disk backend, but it is not enabled: backend=%v",
							common.Backend(bid),
						)
					}
				case common.Backend_UNISTORE:
					if c.Server.Storage.Unistore == nil {
						return fmt.Errorf(
							"config: composite backend depends on Unistore backend, but it is not enabled: backend=%v",
							common.Backend(bid),
						)
					}
				case common.Backend_B2:
					if c.Server.Storage.B2 == nil {
						return fmt.Errorf(
							"config: composite backend depends on B2 backend, but it is not enabled: backend=%v",
							common.Backend(bid),
						)
					}
				}
			}
		}

		if disk := c.Server.Storage.Disk; disk != nil {
			if disk.Root == "" {
				return fmt.Errorf("config: disk backend root path not defined")
			}

			if checksum := disk.Checksum; checksum != nil {
				if _, ok := common.Checksum_value[strings.ToUpper(checksum.Algorithm)]; !ok {
					return fmt.Errorf(
						"config: unknown disk checksum algorithm: algorithm=%s candidates=%v",
						checksum.Algorithm,
						common.Checksum_value,
					)
				}
			}

			if compression := disk.Compression; compression != nil {
				if _, ok := common.Compression_value[strings.ToUpper(compression.Algorithm)]; !ok {
					return fmt.Errorf(
						"config: unknown disk compression algorithm: algorithm=%s candidates=%v",
						compression.Algorithm,
						common.Compression_value,
					)
				}
			}

			if encryption := disk.Encryption; encryption != nil {
				if encryption.Name == "" {
					return fmt.Errorf("config: disk encryption keypair name not defined")
				} else if _, ok := common.Encryption_value[strings.ToUpper(encryption.Mechanism)]; !ok {
					return fmt.Errorf(
						"config: unknown disk encryption mechanism: mechanism=%s candidates=%v",
						encryption.Mechanism,
						common.Encryption_value,
					)
				}
			}
		}

		if unistore := c.Server.Storage.Unistore; unistore != nil {
			if unistore.Address == "" {
				return fmt.Errorf("config: Unistore remote address not defined")
			}

			if checksum := unistore.Checksum; checksum != nil {
				if _, ok := common.Checksum_value[strings.ToUpper(checksum.Algorithm)]; !ok {
					return fmt.Errorf(
						"config: unknown Unistore checksum algorithm: algorithm=%s candidates=%v",
						checksum.Algorithm,
						common.Checksum_value,
					)
				}
			}

			if compression := unistore.Compression; compression != nil {
				if _, ok := common.Compression_value[strings.ToUpper(compression.Algorithm)]; !ok {
					return fmt.Errorf(
						"config: unknown Unistore compression algorithm: algorithm=%s candidates=%v",
						compression.Algorithm,
						common.Compression_value,
					)
				}
			}

			if encryption := unistore.Encryption; encryption != nil {
				if encryption.Name == "" {
					return fmt.Errorf("config: Unistore encryption keypair name not defined")
				} else if _, ok := common.Encryption_value[strings.ToUpper(encryption.Mechanism)]; !ok {
					return fmt.Errorf(
						"config: unknown Unistore encryption mechanism: mechanism=%s candidates=%v",
						encryption.Mechanism,
						common.Encryption_value,
					)
				}
			}
		}

		if b2 := c.Server.Storage.B2; b2 != nil {
			if b2.ApplicationKeyID == "" {
				return fmt.Errorf("config: B2 application key ID not defined")
			} else if b2.ApplicationKey == "" {
				return fmt.Errorf("config: B2 application key not defined")
			}

			if checksum := b2.Checksum; checksum != nil {
				if _, ok := common.Checksum_value[strings.ToUpper(checksum.Algorithm)]; !ok {
					return fmt.Errorf(
						"config: unknown B2 checksum algorithm: algorithm=%s candidates=%v",
						checksum.Algorithm,
						common.Checksum_value,
					)
				}
			}

			if compression := b2.Compression; compression != nil {
				if _, ok := common.Compression_value[strings.ToUpper(compression.Algorithm)]; !ok {
					return fmt.Errorf(
						"config: unknown B2 compression algorithm: algorithm=%s candidates=%v",
						compression.Algorithm,
						common.Compression_value,
					)
				}
			}

			if encryption := b2.Encryption; encryption != nil {
				if encryption.Name == "" {
					return fmt.Errorf("config: B2 encryption keypair name not defined")
				} else if _, ok := common.Encryption_value[strings.ToUpper(encryption.Mechanism)]; !ok {
					return fmt.Errorf(
						"config: unknown B2 encryption mechanism: mechanism=%s candidates=%v",
						encryption.Mechanism,
						common.Encryption_value,
					)
				}
			}
		}
	}

	if c.Client != nil {
		for _, store := range c.Client.Stores {
			if store.Name == "" {
				return fmt.Errorf("config: store name not defined")
			} else if store.Address == "" {
				return fmt.Errorf("config: store address not defined: name=%s", store.Name)
			} else if store.Backend == "" {
				return fmt.Errorf("config: store backend is not defined: name=%s", store.Name)
			}

			if _, ok := common.Backend_value[strings.ToUpper(store.Backend)]; !ok {
				return fmt.Errorf(
					"config: unknown store backend: backend=%s candidates=%v",
					store.Backend,
					common.Backend_value,
				)
			}
		}
	}

	return nil
}
