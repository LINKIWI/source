package aperture

import (
	"net"
	"os"
	"path"
	"strings"
	"time"

	"golang.org/x/sys/unix"

	"lib.kevinlin.info/aperture/internal/data"
	"lib.kevinlin.info/aperture/internal/errors"
	"lib.kevinlin.info/aperture/logger"
	"lib.kevinlin.info/aperture/protocol"
)

// Config is the client-specified configuration for the statsd client.
type Config struct {
	// URI describing the statsd endpoint. Supported schemes are "udp", "tcp", "unix", and
	// "file" for UDP, TCP, UDS (Unix datagram socket), and local disk file transports,
	// respectively. A UDP or TCP transport should specify the network address, a Unix socket
	// transport should specify the path to a socket file, and a file transport should specify
	// the path to a preexisting or nonexistent file on disk.
	//
	// Multiple addresses may be specified, delimited by a comma, to indicate use of a
	// composite transport that will concurrently write metrics to all addresses. By default,
	// the composite transport operates at the strongest write consistency level; it waits until
	// all writes have completed and will consider the metric emission to have failed if any of
	// the underlying transports fail.
	//
	// For example:
	// A UDP transport should be shaped as "udp://localhost:8125".
	// A TCP transport should be shaped as "tcp://localhost:8125".
	// A UDS transport should be shaped as "unix:///var/run/metrics/statsd.sock".
	// A file transport should be shaped as "file:///var/log/metrics/statsd.log".
	// A composite transport that writes to multiple addresses may be shaped as
	// "udp://localhost:8125,tcp://localhost:8125,file:///var/log/metrics/statsd.log".
	Address string

	// Prefix to prepend to all emitted metrics.
	Prefix string

	// Suffix to append to all emitted metrics.
	Suffix string

	// Set of defaults tags to include on all emitted metrics.
	DefaultTags map[string]interface{}

	// Sample rate to share across all metric emissions for which sampling is supported.
	SampleRate float64

	// Used for buffering metric emissions. When non-negative, a buffered transport is used that
	// batches sending of metrics once the threshold queue size has been reached. When zero or
	// negative, buffering is disabled; each metric is sent over the transport immediately.
	BufferSize int

	// Used for asynchronous delivery of metrics. When non-negative, an asynchronous transport
	// is used that buffers outbound metrics into a bounded queue that is asynchronously drained
	// to the underlying transport. When zero or negative, asynchronous delivery is disabled;
	// each metric is sent over the transport immediately.
	AsyncQueueSize int

	// URI describing the address of a SOCKS5 proxy server through which the transport should be
	// established. Supported schemes are "tcp" and "unix" for TCP and Unix domain socket proxy
	// server listeners, respectively.
	//
	// Only statsd emissions through TCP transports support use of a proxy server; this setting
	// is ignored for all other transport types. When omitted, no proxy is used.
	Proxy string

	// Used for automatic and asynchronous connection reestablishment to the network transport
	// on metric emission errors. When nonzero, this duration defines the interval at which the
	// current transport's connection may be reestablished if the number of transmission errors
	// accumulated in the current window is nonzero. When zero, health probing of the network
	// transport is disabled at the library layer, and it is the responsibility of the client
	// application to reestablish a new instance as needed in response to write errors.
	TransportProbeInterval time.Duration

	// Requested time interval between keepalive probes for the underlying connection, relevant
	// only when using a reliable transport (i.e. TCP) with support from the underlying OS. When
	// zero or negative, keepalive probes are disabled.
	TransportKeepaliveInterval time.Duration

	// Absolute timeout for establishing/dialing a new connection, relevant only when using a
	// reliable transport (i.e. TCP). When zero, no application-layer timeout is applied.
	TransportConnectTimeout time.Duration

	// Absolute timeout for individual payload transmissions over the network, relevant only
	// when using a reliable transport (i.e. TCP). When zero, no application-layer transmission
	// timeout is applied.
	TransportSendTimeout time.Duration

	// Initialize the backing transport lazily on the first metric emission instead of eagerly
	// on client initialization. This has the effect of potentially masking transport errors
	// until late into the application lifecycle, though this may be desirable when using a
	// reliable transport (e.g. TCP or UDS) to eliminate the requirement that the backend be
	// available before the client is initialized. It is recommended that this option be used in
	// conjunction with active transport probing.
	LazyTransportInitialization bool

	// Specify a metrics logger used for logging all metrics emitted by this client. When
	// omitted, a noop logger is used, which effectively disables any logging.
	Logger logger.Logger

	// Override the line protocol serializer implementation. This can be used, for example, to
	// use a custom tags serialization mechanism, or to conform to a different metrics protocol
	// that shares statsd metric types. When omitted, the default statsd serializer is used.
	Serializer protocol.LineSerializer
}

// Validate the supplied configuration parameters.
func (c *Config) validate() error {
	if c.Address != "" {
		for _, addr := range strings.Split(c.Address, ",") {
			uri, err := data.ParseAddressURI(addr)
			if err != nil {
				return errors.New(
					"config",
					"malformed statsd server address URI",
					errors.Tag{Key: "addr", Value: addr},
				)
			}

			switch uri.Scheme {
			case "udp", "tcp":
				if _, _, err := net.SplitHostPort(uri.Authority); err != nil {
					return errors.Stack(
						"config",
						"failed to parse network address",
						err,
						errors.Tag{Key: "addr", Value: addr},
					)
				}
			case "unix":
				file, err := os.Stat(uri.Authority)
				if err != nil {
					return errors.Stack(
						"config",
						"failed to stat socket path",
						err,
						errors.Tag{Key: "addr", Value: addr},
					)
				} else if file.Mode()&os.ModeSocket != os.ModeSocket {
					return errors.New(
						"config",
						"socket path is not a Unix socket",
						errors.Tag{Key: "addr", Value: addr},
					)
				}
			case "file":
				file, err := os.Stat(uri.Authority)
				if err == nil && !file.Mode().IsRegular() {
					return errors.New(
						"config",
						"file path is not a regular file",
						errors.Tag{Key: "addr", Value: addr},
					)
				}

				if unix.Access(uri.Authority, unix.W_OK) != nil &&
					unix.Access(path.Dir(uri.Authority), unix.W_OK) != nil {
					return errors.New(
						"config",
						"both file path and its parent directory are not writable",
						errors.Tag{Key: "addr", Value: addr},
					)
				}
			default:
				return errors.New(
					"config",
					"unsupported statsd server address scheme",
					errors.Tag{Key: "addr", Value: addr},
				)
			}
		}
	}

	if c.Proxy != "" {
		uri, err := data.ParseAddressURI(c.Proxy)
		if err != nil {
			return errors.New(
				"config",
				"malformed proxy server address URI",
				errors.Tag{Key: "addr", Value: c.Proxy},
			)
		}

		switch uri.Scheme {
		case "tcp":
			if _, _, err := net.SplitHostPort(uri.Authority); err != nil {
				return errors.Stack(
					"config",
					"failed to parse proxy server address",
					err,
					errors.Tag{Key: "proxy", Value: c.Proxy},
				)
			}
		case "unix":
			file, err := os.Stat(uri.Authority)
			if err != nil {
				return errors.Stack(
					"config",
					"failed to stat proxy server socket path",
					err,
					errors.Tag{Key: "addr", Value: c.Proxy},
				)
			} else if file.Mode()&os.ModeSocket != os.ModeSocket {
				return errors.New(
					"config",
					"proxy server socket path is not a Unix socket",
					errors.Tag{Key: "addr", Value: c.Proxy},
				)
			}
		default:
			return errors.New(
				"config",
				"unsupported proxy server address scheme",
				errors.Tag{Key: "addr", Value: c.Proxy},
			)
		}
	}

	if c.SampleRate < 0.0 || c.SampleRate > 1.0 {
		return errors.New(
			"config",
			"sample rate must be in range [0.0, 1.0]",
			errors.Tag{Key: "sample_rate", Value: c.SampleRate},
		)
	}

	if c.BufferSize < 0 {
		return errors.New(
			"config",
			"buffer size must be non-negative",
			errors.Tag{Key: "buffer_size", Value: c.BufferSize},
		)
	}

	if c.AsyncQueueSize < 0 {
		return errors.New(
			"config",
			"async queue size must be non-negative",
			errors.Tag{Key: "async_queue_size", Value: c.AsyncQueueSize},
		)
	}

	if c.TransportProbeInterval < 0 {
		return errors.New(
			"config",
			"transport probe interval must be non-negative",
			errors.Tag{Key: "probe_interval", Value: c.TransportProbeInterval},
		)
	}

	return nil
}
