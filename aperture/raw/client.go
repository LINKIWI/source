package aperture

import (
	"fmt"
	"strings"
	"time"

	"lib.kevinlin.info/aperture/internal/data"
	"lib.kevinlin.info/aperture/internal/errors"
	"lib.kevinlin.info/aperture/logger"
	"lib.kevinlin.info/aperture/protocol"
	"lib.kevinlin.info/aperture/transport"
)

// Statsd is an interface for a client providing metrics reporting to a statsd backend.
type Statsd interface {
	Count(name string, delta int64, tags map[string]interface{}) error
	Incr(name string, tags map[string]interface{}) error
	Decr(name string, tags map[string]interface{}) error
	Gauge(name string, value float64, tags map[string]interface{}) error
	Size(name string, bytes int64, tags map[string]interface{}) error
	Timing(name string, duration time.Duration, tags map[string]interface{}) error
	TimingMs(name string, duration float64, tags map[string]interface{}) error
	Histogram(name string, value float64, tags map[string]interface{}) error
	Close() error
}

// Client is the Aperture implementation of a statsd client.
type Client struct {
	cfg   *Config
	tport transport.Transport
}

// NewClient creates a client with the specified configuration.
func NewClient(cfg *Config) (*Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, errors.Stack("client", "client config validation failed", err)
	}

	/* Transport instantiation */

	var backends []transport.Factory

	tFactory := func() (transport.Transport, error) {
		return transport.NewNoop(), nil
	}

	if cfg.Address != "" {
		for _, addr := range strings.Split(cfg.Address, ",") {
			func(addr string) {
				// Instantiation of the transport client specified by the address

				uri, _ := data.ParseAddressURI(addr)

				switch uri.Scheme {
				case "udp":
					tFactory = func() (transport.Transport, error) {
						return transport.NewUDP(uri.Authority)
					}
				case "tcp":
					tFactory = func() (transport.Transport, error) {
						opts := transport.TCPSocketConfig{
							ConnectTimeout:    cfg.TransportConnectTimeout,
							SendTimeout:       cfg.TransportSendTimeout,
							KeepAliveInterval: cfg.TransportKeepaliveInterval,
						}

						if cfg.Proxy == "" {
							return transport.NewDirectTCP(
								uri.Authority,
								opts,
							)
						}

						proxy, _ := data.ParseAddressURI(cfg.Proxy)
						return transport.NewProxiedTCP(
							uri.Authority,
							proxy.Scheme,
							proxy.Authority,
							opts,
						)
					}
				case "unix":
					tFactory = func() (transport.Transport, error) {
						return transport.NewUDS(uri.Authority)
					}
				case "file":
					tFactory = func() (transport.Transport, error) {
						return transport.NewFile(uri.Authority)
					}
				}

				// Transport usage abstractions on top of backends

				backendTFactory := tFactory
				tFactory = func() (transport.Transport, error) {
					return transport.NewFramed(backendTFactory, []byte("\n"))
				}

				if cfg.LazyTransportInitialization {
					backendTFactory := tFactory
					tFactory = func() (transport.Transport, error) {
						return transport.NewLazy(backendTFactory)
					}
				}

				if cfg.TransportProbeInterval > 0 {
					backendTFactory := tFactory
					tFactory = func() (transport.Transport, error) {
						return transport.NewReconnecting(
							backendTFactory,
							cfg.TransportProbeInterval,
							1,
						)
					}
				}

				backends = append(backends, tFactory)
			}(addr)
		}
	}

	if len(backends) == 1 {
		tFactory = backends[0]
	} else if len(backends) > 1 {
		tFactory = func() (transport.Transport, error) {
			return transport.NewComposite(backends, transport.All)
		}
	}

	if cfg.AsyncQueueSize > 0 {
		backendTFactory := tFactory
		tFactory = func() (transport.Transport, error) {
			return transport.NewAsync(
				backendTFactory,
				cfg.AsyncQueueSize,
			)
		}
	}

	if cfg.BufferSize > 0 {
		backendTFactory := tFactory
		tFactory = func() (transport.Transport, error) {
			return transport.NewBuffered(backendTFactory, cfg.BufferSize)
		}
	}

	tport, err := tFactory()
	if err != nil {
		return nil, errors.Stack("client", "failed to create transport", err)
	}

	/* Configuration defaults */

	if cfg.Logger == nil {
		cfg.Logger = logger.NewNoop()
	}

	if cfg.Serializer == nil {
		cfg.Serializer = protocol.DefaultStatsdLineSerializer
	}

	if cfg.SampleRate == 0.0 {
		cfg.SampleRate = 1.0
	}

	return &Client{
		cfg:   cfg,
		tport: tport,
	}, nil
}

// Count emits a counter metric.
func (s *Client) Count(name string, delta int64, tags map[string]interface{}) error {
	counter := &protocol.Counter{
		Timestamp:  time.Now(),
		Name:       s.formatName(name),
		Value:      delta,
		SampleRate: s.cfg.SampleRate,
		Tags:       s.mergeTags(tags),
	}

	s.cfg.Logger.LogCounter(counter)

	line, err := s.cfg.Serializer.SerializeCounter(counter)
	if err != nil {
		return err
	}

	if _, err := s.tport.Send([]byte(line)); err != nil {
		return err
	}

	return nil
}

// Incr is a convenience method for emitting a counter with value 1.
func (s *Client) Incr(name string, tags map[string]interface{}) error {
	return s.Count(name, 1, tags)
}

// Decr is a convenience method for emitting a counter with value -1.
func (s *Client) Decr(name string, tags map[string]interface{}) error {
	return s.Count(name, -1, tags)
}

// Gauge emits a gauge metric.
func (s *Client) Gauge(name string, value float64, tags map[string]interface{}) error {
	gauge := &protocol.Gauge{
		Timestamp: time.Now(),
		Name:      s.formatName(name),
		Value:     value,
		Tags:      s.mergeTags(tags),
	}

	s.cfg.Logger.LogGauge(gauge)

	line, err := s.cfg.Serializer.SerializeGauge(gauge)
	if err != nil {
		return err
	}

	if _, err := s.tport.Send([]byte(line)); err != nil {
		return err
	}

	return nil
}

// Size is used to describe the size of a payload in byte units. In practice, this is merely an
// alias for emitting a timer metric.
func (s *Client) Size(name string, bytes int64, tags map[string]interface{}) error {
	// Size metrics share the same semantics as timing metrics; they are interpreted and
	// aggregated in the same way.
	return s.TimingMs(name, float64(bytes), tags)
}

// Timing is a convenience method for emitting a timer with a duration specified as type
// time.Duration.
func (s *Client) Timing(name string, duration time.Duration, tags map[string]interface{}) error {
	return s.TimingMs(name, float64(duration)/float64(time.Millisecond), tags)
}

// TimingMs emits a timer metric for a duration in milliseconds.
func (s *Client) TimingMs(name string, duration float64, tags map[string]interface{}) error {
	timer := &protocol.Timer{
		Timestamp:  time.Now(),
		Name:       s.formatName(name),
		Value:      duration,
		SampleRate: s.cfg.SampleRate,
		Tags:       s.mergeTags(tags),
	}

	s.cfg.Logger.LogTimer(timer)

	line, err := s.cfg.Serializer.SerializeTimer(timer)
	if err != nil {
		return err
	}

	if _, err := s.tport.Send([]byte(line)); err != nil {
		return err
	}

	return nil
}

// Histogram emits a histogram metric.
func (s *Client) Histogram(name string, value float64, tags map[string]interface{}) error {
	histogram := &protocol.Histogram{
		Timestamp: time.Now(),
		Name:      s.formatName(name),
		Value:     value,
		Tags:      s.mergeTags(tags),
	}

	s.cfg.Logger.LogHistogram(histogram)

	line, err := s.cfg.Serializer.SerializeHistogram(histogram)
	if err != nil {
		return err
	}

	if _, err := s.tport.Send([]byte(line)); err != nil {
		return err
	}

	return nil
}

// Close closes the underlying network transport.
func (s *Client) Close() error {
	return s.tport.Close()
}

// formatName adds a prefix and suffix, if specified, to the metric name.
func (s *Client) formatName(metric string) string {
	var formatted strings.Builder

	if s.cfg.Prefix != "" {
		formatted.WriteString(s.cfg.Prefix)
		formatted.WriteString(".")
	}

	formatted.WriteString(metric)

	if s.cfg.Suffix != "" {
		formatted.WriteString(".")
		formatted.WriteString(s.cfg.Suffix)
	}

	return formatted.String()
}

// mergeTags merges configuration-supplied default tag with the per-metric tags. Metric tags take
// precedence over client-global default tags (i.e., it can override default tags).
func (s *Client) mergeTags(tags map[string]interface{}) map[string]string {
	mergedTags := make(map[string]string, len(s.cfg.DefaultTags)+len(tags))

	for key, value := range s.cfg.DefaultTags {
		mergedTags[key] = fmt.Sprintf("%v", value)
	}

	if tags != nil {
		for key, value := range tags {
			mergedTags[key] = fmt.Sprintf("%v", value)
		}
	}

	return mergedTags
}
