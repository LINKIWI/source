package config

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/proxy"
	"gopkg.in/yaml.v3"

	"courier/internal/util"
)

// New creates a configuration object from a path on disk.
func New(path string) (*Config, error) {
	var cfg *Config

	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, &util.Error{
			Namespace: "config",
			Message:   "error opening config file on disk",
			Tags: map[string]interface{}{
				"path": path,
			},
			StackedError: err,
		}
	}

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	if err := decoder.Decode(&cfg); err != nil {
		return nil, &util.Error{
			Namespace: "config",
			Message:   "error deserializing config file contents",
			Tags: map[string]interface{}{
				"path": path,
			},
			StackedError: err,
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, &util.Error{
			Namespace: "config",
			Message:   "error validating config",
			Tags: map[string]interface{}{
				"path": path,
			},
			StackedError: err,
		}
	}

	return cfg, nil
}

// Hash generates a deterministic and unique identifier for the currently active config.
func (c *Config) Hash() string {
	encoded, _ := yaml.Marshal(c)
	return fmt.Sprintf("%x", sha256.Sum256(encoded))
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	vhosts := make(map[string]bool)

	if c.Application == nil {
		return &util.Error{
			Namespace: "config",
			Message:   "missing application configuration",
		}
	} else if c.Application.Instance == "" {
		return &util.Error{
			Namespace: "config",
			Message:   "missing application instance identifier",
		}
	} else if c.Application.Metrics != nil {
		if _, ok := metricSerializers[c.Application.Metrics.Serializer]; !ok {
			return &util.Error{
				Namespace: "config",
				Message:   "unknown metric serializer",
				Tags: map[string]interface{}{
					"serializer": c.Application.Metrics.Serializer,
				},
			}
		}
	}

	if c.Server == nil {
		return &util.Error{
			Namespace: "config",
			Message:   "missing server configuration",
		}
	} else if len(c.Server.Listeners) == 0 {
		return &util.Error{
			Namespace: "config",
			Message:   "missing server listeners",
		}
	}

	for _, listener := range c.Server.Listeners {
		if listener.Name == "" {
			return &util.Error{
				Namespace: "config",
				Message:   "listener is missing name",
			}
		} else if listener.Address == nil {
			return &util.Error{
				Namespace: "config",
				Message:   "listener missing address",
				Tags: map[string]interface{}{
					"name": listener.Name,
				},
			}
		} else if _, ok := connectionProtocols[listener.Protocol]; !ok {
			return &util.Error{
				Namespace: "config",
				Message:   "unknown listener protocol",
				Tags: map[string]interface{}{
					"name":     listener.Name,
					"protocol": listener.Protocol,
				},
			}
		}

		if listener.TLSContext != nil {
			if listener.TLSContext.Certificates == nil || len(listener.TLSContext.Certificates) == 0 {
				return &util.Error{
					Namespace: "config",
					Message:   "missing server TLS certificates",
					Tags: map[string]interface{}{
						"name": listener.Name,
					},
				}
			} else if _, ok := tlsVerifyModes[listener.TLSContext.VerifyMode]; !ok {
				return &util.Error{
					Namespace: "config",
					Message:   "unknown TLS client verification mode",
					Tags: map[string]interface{}{
						"name":       listener.Name,
						"mode":       listener.TLSContext.VerifyMode,
						"candidates": tlsVerifyModes,
					},
				}
			} else if listener.TLSContext.VerifyPeer != nil && listener.TLSContext.VerifyMode != "strict" {
				return &util.Error{
					Namespace: "config",
					Message:   "TLS verify mode must be strict to enable peer identity verification",
					Tags: map[string]interface{}{
						"name": listener.Name,
						"mode": listener.TLSContext.VerifyMode,
						"peer": listener.TLSContext.VerifyPeer,
					},
				}
			} else if _, ok := tlsProtocolVersions[listener.TLSContext.ProtocolVersion]; !ok {
				return &util.Error{
					Namespace: "config",
					Message:   "unknown TLS protocol version constraint",
					Tags: map[string]interface{}{
						"name":       listener.Name,
						"constraint": listener.TLSContext.ProtocolVersion,
						"candidates": tlsProtocolVersions,
					},
				}
			}

			for _, cert := range listener.TLSContext.Certificates {
				if cert.Key == "" {
					return &util.Error{
						Namespace: "config",
						Message:   "TLS certificate is missing key path",
						Tags: map[string]interface{}{
							"name": listener.Name,
						},
					}
				} else if cert.Certificate == "" {
					return &util.Error{
						Namespace: "config",
						Message:   "TLS certificate is missing certificate path",
						Tags: map[string]interface{}{
							"name": listener.Name,
						},
					}
				}
			}
		}
	}

	for _, vhost := range c.Server.VirtualHosts {
		upstreams := make(map[string]bool)

		if _, ok := vhosts[vhost.Name]; ok {
			return &util.Error{
				Namespace: "config",
				Message:   "duplicated virtual host name",
				Tags: map[string]interface{}{
					"name": vhost.Name,
				},
			}
		} else if len(vhost.Upstreams) == 0 {
			return &util.Error{
				Namespace: "config",
				Message:   "virtual host missing upstream configuration",
				Tags: map[string]interface{}{
					"name": vhost.Name,
				},
			}
		} else if len(vhost.Routes) == 0 {
			return &util.Error{
				Namespace: "config",
				Message:   "virtual host missing route configuration",
				Tags: map[string]interface{}{
					"name": vhost.Name,
				},
			}
		}

		vhosts[vhost.Name] = true

		for _, upstream := range vhost.Upstreams {
			if upstream.Name == "" {
				return &util.Error{
					Namespace: "config",
					Message:   "upstream missing name",
					Tags: map[string]interface{}{
						"vhost": vhost.Name,
					},
				}
			} else if _, ok := upstreams[upstream.Name]; ok {
				return &util.Error{
					Namespace: "config",
					Message:   "duplicated upstream name",
					Tags: map[string]interface{}{
						"vhost":    vhost.Name,
						"upstream": upstream.Name,
					},
				}
			} else if upstream.Address == nil {
				return &util.Error{
					Namespace: "config",
					Message:   "upstream missing address",
					Tags: map[string]interface{}{
						"vhost":    vhost.Name,
						"upstream": upstream.Name,
					},
				}
			} else if _, ok := upstreamTransports[upstream.Transport]; !ok {
				return &util.Error{
					Namespace: "config",
					Message:   "unknown virtual host transport",
					Tags: map[string]interface{}{
						"vhost":     vhost.Name,
						"upstream":  upstream.Name,
						"transport": upstream.Transport,
					},
				}
			}

			if upstream.HealthCheck != nil {
				if _, ok := healthCheckMethods[upstream.HealthCheck.Method]; !ok {
					return &util.Error{
						Namespace: "config",
						Message:   "unknown HTTP health check method",
						Tags: map[string]interface{}{
							"vhost":    vhost.Name,
							"upstream": upstream.Name,
							"method":   upstream.HealthCheck.Method,
						},
					}
				} else if upstream.HealthCheck.Interval <= 0 {
					return &util.Error{
						Namespace: "config",
						Message:   "health check interval must be a positive value",
						Tags: map[string]interface{}{
							"vhost":    vhost.Name,
							"upstream": upstream.Name,
							"interval": upstream.HealthCheck.Interval,
						},
					}
				}
			}

			if upstream.TLSContext != nil {
				if _, ok := tlsVerifyModes[upstream.TLSContext.VerifyMode]; !ok {
					return &util.Error{
						Namespace: "config",
						Message:   "unknown TLS server verification mode",
						Tags: map[string]interface{}{
							"vhost":      vhost.Name,
							"upstream":   upstream.Name,
							"mode":       upstream.TLSContext.VerifyMode,
							"candidates": tlsVerifyModes,
						},
					}
				} else if upstream.TLSContext.VerifyPeer != nil && upstream.TLSContext.VerifyMode != "strict" {
					return &util.Error{
						Namespace: "config",
						Message:   "TLS verify mode must be strict to enable peer identity verification",
						Tags: map[string]interface{}{
							"vhost":    vhost.Name,
							"upstream": upstream.Name,
							"mode":     upstream.TLSContext.VerifyMode,
							"peer":     upstream.TLSContext.VerifyPeer,
						},
					}
				} else if _, ok := tlsProtocolVersions[upstream.TLSContext.ProtocolVersion]; !ok {
					return &util.Error{
						Namespace: "config",
						Message:   "unknown TLS protocol version constraint",
						Tags: map[string]interface{}{
							"vhost":      vhost.Name,
							"upstream":   upstream.Name,
							"constraint": upstream.TLSContext.ProtocolVersion,
							"candidates": tlsProtocolVersions,
						},
					}
				}
			}

			upstreams[upstream.Name] = true
		}

		for _, route := range vhost.Routes {
			if route.Upstream == "" {
				return &util.Error{
					Namespace: "config",
					Message:   "route upstream name is missing",
					Tags: map[string]interface{}{
						"vhost":      vhost.Name,
						"candidates": upstreams,
					},
				}
			} else if _, ok := upstreams[route.Upstream]; !ok {
				return &util.Error{
					Namespace: "config",
					Message:   "unknown route upstream",
					Tags: map[string]interface{}{
						"vhost":      vhost.Name,
						"upstream":   route.Upstream,
						"candidates": upstreams,
					},
				}
			} else if route.Match.Header != nil {
				if route.Match.Header.Key == "" {
					return &util.Error{
						Namespace: "config",
						Message:   "route match header key is missing",
						Tags: map[string]interface{}{
							"vhost":    vhost.Name,
							"upstream": route.Upstream,
						},
					}
				} else if route.Match.Header.Value == nil {
					return &util.Error{
						Namespace: "config",
						Message:   "route match header value is missing",
						Tags: map[string]interface{}{
							"vhost":    vhost.Name,
							"upstream": route.Upstream,
							"key":      route.Match.Header.Key,
						},
					}
				}
			}
		}

		for _, filter := range vhost.Filters {
			if filter.Name == "" {
				return &util.Error{
					Namespace: "config",
					Message:   "missing filter name",
					Tags:      map[string]interface{}{"vhost": vhost.Name},
				}
			} else if filter.Type == "" {
				return &util.Error{
					Namespace: "config",
					Message:   "missing filter type",
					Tags:      map[string]interface{}{"vhost": vhost.Name},
				}
			} else if _, ok := filterTypes[filter.Type]; !ok {
				return &util.Error{
					Namespace: "config",
					Message:   "unknown filter type",
					Tags: map[string]interface{}{
						"vhost": vhost.Name,
						"type":  filter.Type,
					},
				}
			}
		}
	}

	return nil
}

// ClientTransport creates a new http.RoundTripper transport derived from the upstream connection
// configuration. Note that this method initializes new resources and the returned transport should
// thus be reused across multiple connections for this upstream in order to avoid resource leaks.
func (u *Upstream) ClientTransport() (http.RoundTripper, error) {
	var err error
	var tlsCfg *tls.Config

	if u.TLSContext != nil {
		zap.L().Debug(
			"creating client TLS origination context",
			zap.String("upstream", u.Address.String()),
		)

		tlsCfg, err = u.TLSContext.ClientConfig()
		if err != nil {
			return nil, &util.Error{
				Namespace:    "config",
				Message:      "error creating client TLS origination context",
				StackedError: err,
			}
		}
	}

	switch u.Transport {
	case "http2":
		zap.L().Debug(
			"creating client HTTP/2 transport",
			zap.String("upstream", u.Address.String()),
			zap.Stringer("proxy", u.Proxy),
			zap.Any("connection", u.Connection),
		)

		return &http2.Transport{
			AllowHTTP:       true,
			TLSClientConfig: tlsCfg,
			DialTLS: func(network string, addr string, cfg *tls.Config) (net.Conn, error) {
				var dialer proxy.Dialer
				var cancel context.CancelFunc

				dialer = &net.Dialer{
					Timeout:   u.Connection.ConnectTimeout,
					KeepAlive: u.Connection.KeepaliveInterval,
				}

				if u.TLSContext != nil {
					dialer = &tls.Dialer{
						NetDialer: dialer.(*net.Dialer),
						Config:    cfg,
					}
				}

				if u.Proxy != nil {
					dialer, err = proxy.FromURL(u.Proxy.URL, dialer)
					if err != nil {
						return nil, &util.Error{
							Namespace: "config",
							Message:   "error creating HTTP/2 proxy dialer",
							Tags: map[string]interface{}{
								"name":  u.Name,
								"addr":  u.Address,
								"proxy": u.Proxy,
							},
							StackedError: err,
						}
					}
				}

				network, addr = u.Address.Address()
				conn, err := dialer.Dial(network, addr)
				if err != nil {
					return nil, &util.Error{
						Namespace: "config",
						Message:   "error dialing upstream address over HTTP/2",
						Tags: map[string]interface{}{
							"name": u.Name,
							"addr": u.Address,
						},
						StackedError: err,
					}
				}

				if u.TLSContext != nil {
					ctx := context.Background()
					if u.TLSContext.HandshakeTimeout > 0 {
						ctx, cancel = context.WithTimeout(
							context.Background(),
							u.TLSContext.HandshakeTimeout,
						)
						defer cancel()
					}

					if err := conn.(*tls.Conn).HandshakeContext(ctx); err != nil {
						return nil, &util.Error{
							Namespace: "config",
							Message:   "error during upstream TLS handshake",
							Tags: map[string]interface{}{
								"name": u.Name,
								"addr": u.Address,
							},
							StackedError: err,
						}
					}
				}

				return newTimeoutConn(
					u.Connection.ReadTimeout,
					u.Connection.WriteTimeout,
					conn,
				), nil
			},
		}, nil
	default:
		var tlsHandshakeTimeout time.Duration

		zap.L().Debug(
			"creating client HTTP/1.1 transport",
			zap.String("upstream", u.Address.String()),
			zap.Stringer("proxy", u.Proxy),
			zap.Any("connection", u.Connection),
		)

		if u.TLSContext != nil {
			tlsHandshakeTimeout = u.TLSContext.HandshakeTimeout
		}

		return &http.Transport{
			ResponseHeaderTimeout: u.Connection.ReadHeaderTimeout,
			ExpectContinueTimeout: u.Connection.ExpectContinueTimeout,
			IdleConnTimeout:       u.Connection.IdleTimeout,
			TLSHandshakeTimeout:   tlsHandshakeTimeout,
			MaxIdleConnsPerHost:   u.Connection.IdleLimit,
			DisableKeepAlives:     u.Connection.KeepaliveInterval < 0,
			ReadBufferSize:        u.Connection.ReadBufferSize,
			WriteBufferSize:       u.Connection.WriteBufferSize,
			TLSClientConfig:       tlsCfg,
			// Even if an HTTP/2 upstream transport is not explicitly requested, attempt
			// an upgrade regardless and gracefully fall back to HTTP/1.1 when the
			// upstream does not support HTTP/2 or otherwise refuses the upgrade.
			ForceAttemptHTTP2: true,
			DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
				var dialer proxy.Dialer

				dialer = &net.Dialer{
					Timeout:   u.Connection.ConnectTimeout,
					KeepAlive: u.Connection.KeepaliveInterval,
				}

				if u.Proxy != nil {
					dialer, err = proxy.FromURL(u.Proxy.URL, dialer)
					if err != nil {
						return nil, &util.Error{
							Namespace: "config",
							Message:   "error creating proxy dialer",
							Tags: map[string]interface{}{
								"name":  u.Name,
								"addr":  u.Address,
								"proxy": u.Proxy,
							},
							StackedError: err,
						}
					}
				}

				network, addr = u.Address.Address()
				conn, err := dialer.Dial(network, addr)
				if err != nil {
					return nil, &util.Error{
						Namespace: "config",
						Message:   "error dialing upstream address",
						Tags: map[string]interface{}{
							"name": u.Name,
							"addr": u.Address,
						},
						StackedError: err,
					}
				}

				return newTimeoutConn(
					u.Connection.ReadTimeout,
					u.Connection.WriteTimeout,
					conn,
				), nil
			},
		}, nil
	}
}

// ClientConfig transforms the TLS context into a tls.Config appropriate for purposes of TLS
// origination from a client.
func (t *TLSContext) ClientConfig() (*tls.Config, error) {
	var verifyConnection func(tls.ConnectionState) error
	var minVersion uint16

	cs, err := newCertificateStore(t, func() (*x509.CertPool, error) {
		return x509.SystemCertPool()
	})
	if err != nil {
		return nil, &util.Error{
			Namespace:    "config",
			Message:      "failed to initialize client certificate store",
			StackedError: err,
		}
	}

	if t.VerifyPeer != nil {
		pv := &peerVerifier{
			pattern: t.VerifyPeer,
			roots:   cs.cas,
		}

		verifyConnection = pv.verifyServer
	}

	switch t.ProtocolVersion {
	case "lenient":
		minVersion = tls.VersionTLS11
	case "standard":
		minVersion = tls.VersionTLS12
	case "modern":
		minVersion = tls.VersionTLS13
	default:
		minVersion = tls.VersionTLS12
	}

	return &tls.Config{
		GetClientCertificate: cs.getClientCertificate,
		RootCAs:              cs.cas,
		// Defer to the overridden connection verification in place of default validation
		// when explicit peer verification is enabled.
		InsecureSkipVerify: t.VerifyMode == "none" ||
			t.VerifyMode == "relaxed" ||
			t.VerifyPeer != nil,
		MinVersion:         minVersion,
		ClientSessionCache: tls.NewLRUClientSessionCache(128),
		VerifyConnection:   verifyConnection,
	}, nil
}

// ServerConfig transforms the TLS context into a tls.Config appropriate for purposes of TLS
// termination at a server.
func (t *TLSContext) ServerConfig() (*tls.Config, error) {
	var verifyMode tls.ClientAuthType
	var verifyConnection func(tls.ConnectionState) error
	var minVersion uint16

	// Explicitly use a nil (empty) seed CA certificate pool in order to use only the CAs
	// specified in TLS configuration for client certificate authentication.
	cs, err := newCertificateStore(t, nil)
	if err != nil {
		return nil, &util.Error{
			Namespace:    "config",
			Message:      "failed to initialize server certificate store",
			StackedError: err,
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

	if t.VerifyPeer != nil {
		pv := &peerVerifier{
			pattern: t.VerifyPeer,
			roots:   cs.cas,
		}

		verifyConnection = pv.verifyClient
	}

	switch t.ProtocolVersion {
	case "lenient":
		minVersion = tls.VersionTLS11
	case "standard":
		minVersion = tls.VersionTLS12
	case "modern":
		minVersion = tls.VersionTLS13
	default:
		minVersion = tls.VersionTLS12
	}

	return &tls.Config{
		GetCertificate:   cs.getServerCertificate,
		ClientCAs:        cs.cas,
		ClientAuth:       verifyMode,
		MinVersion:       minVersion,
		VerifyConnection: verifyConnection,
	}, nil
}
