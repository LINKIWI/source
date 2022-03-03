package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

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

// Validate validates the configuration object.
func (c *Config) Validate() error {
	for _, proxy := range c.Server.Proxies {
		if proxy.Name == "" {
			return fmt.Errorf("config: missing proxy name")
		} else if proxy.Listener == nil {
			return fmt.Errorf("config: missing proxy listener: proxy=%s", proxy.Name)
		} else if proxy.Listener.Address == nil {
			return fmt.Errorf(
				"config: missing proxy listener address: proxy=%s",
				proxy.Name,
			)
		} else if len(proxy.Targets) == 0 {
			return fmt.Errorf(
				"config: missing at least one proxy target specification: proxy=%s",
				proxy.Name,
			)
		} else if _, ok := connectionLoadBalancers[proxy.Options.ConnectionLoadBalancer]; !ok {
			return fmt.Errorf(
				"config: unknown connection load balancing policy: "+
					"proxy=%s policy=%s candidates=%v",
				proxy.Name,
				proxy.Options.ConnectionLoadBalancer,
				connectionLoadBalancers,
			)
		}

		if proxy.Listener.TLSContext != nil {
			if _, ok := tlsVerifyModes[proxy.Listener.TLSContext.VerifyMode]; !ok {
				return fmt.Errorf(
					"config: unknown TLS client verification mode: proxy=%s mode=%s candidates=%v",
					proxy.Name,
					proxy.Listener.TLSContext.VerifyMode,
					tlsVerifyModes,
				)
			} else if len(proxy.Listener.TLSContext.Certificates) == 0 {
				return fmt.Errorf(
					"config: TLS context specified without any server certificates: proxy=%s",
					proxy.Name,
				)
			}

			for _, certificate := range proxy.Listener.TLSContext.Certificates {
				if certificate.Key == "" {
					return fmt.Errorf(
						"config: keypair specification missing path to private key: proxy=%s",
						proxy.Name,
					)
				} else if certificate.Certificate == "" {
					return fmt.Errorf(
						"config: keypair specification missing path to certificate: proxy=%s",
						proxy.Name,
					)
				}
			}
		}

		for _, target := range proxy.Targets {
			if target.Name == "" {
				return fmt.Errorf(
					"config: missing proxy target name: proxy=%s",
					proxy.Name,
				)
			} else if target.Address == nil {
				return fmt.Errorf(
					"config: missing proxy target address: proxy=%s target=%s",
					proxy.Name,
					target.Name,
				)
			}

			if target.TLSContext != nil {
				if _, ok := tlsVerifyModes[target.TLSContext.VerifyMode]; !ok {
					return fmt.Errorf(
						"config: unknown TLS server verification mode: "+
							"proxy=%s target=%s mode=%s candidates=%v",
						proxy.Name,
						target.Name,
						target.TLSContext.VerifyMode,
						tlsVerifyModes,
					)
				}
			}
		}
	}

	return nil
}

// ClientConfig derives a *tls.Config from the declared client TLS origination configuration.
func (t *TLSContext) ClientConfig() (*tls.Config, error) {
	var certificates []tls.Certificate

	caPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf(
			"config: error retrieving system default CA certificate store: err=%v",
			err,
		)
	}

	for _, cert := range t.Certificates {
		keypair, err := tls.LoadX509KeyPair(cert.Certificate, cert.Key)
		if err != nil {
			return nil, fmt.Errorf("config: error loading X.509 keypair: err=%v", err)
		}

		certificates = append(certificates, keypair)

		if cert.Authority != "" {
			ca, err := os.ReadFile(cert.Authority)
			if err != nil {
				return nil, fmt.Errorf(
					"config: error reading certificate authority file: err=%v",
					err,
				)
			}

			caPool.AppendCertsFromPEM(ca)
		}
	}

	return &tls.Config{
		Certificates:       certificates,
		RootCAs:            caPool,
		MinVersion:         tls.VersionTLS13,
		NextProtos:         t.ApplicationProtocols,
		InsecureSkipVerify: t.VerifyMode == "none",
		ClientSessionCache: tls.NewLRUClientSessionCache(128),
	}, nil
}

// ServerConfig derives a *tls.Config from the declared server TLS termination configuration.
func (t *TLSContext) ServerConfig() (*tls.Config, error) {
	var certificates []tls.Certificate
	var verifyMode tls.ClientAuthType

	caPool := x509.NewCertPool()

	for _, cert := range t.Certificates {
		keypair, err := tls.LoadX509KeyPair(cert.Certificate, cert.Key)
		if err != nil {
			return nil, fmt.Errorf("config: error loading X.509 keypair: err=%v", err)
		}

		certificates = append(certificates, keypair)

		if cert.Authority != "" {
			ca, err := os.ReadFile(cert.Authority)
			if err != nil {
				return nil, fmt.Errorf(
					"config: error reading certificate authority file: err=%v",
					err,
				)
			}

			caPool.AppendCertsFromPEM(ca)
		}
	}

	switch t.VerifyMode {
	case "none":
		verifyMode = tls.NoClientCert
	case "relaxed":
		verifyMode = tls.VerifyClientCertIfGiven
	case "strict":
		verifyMode = tls.RequireAndVerifyClientCert
	default:
		verifyMode = tls.NoClientCert
	}

	return &tls.Config{
		Certificates: certificates,
		ClientCAs:    caPool,
		ClientAuth:   verifyMode,
		NextProtos:   t.ApplicationProtocols,
		MinVersion:   tls.VersionTLS13,
	}, nil
}
