# bitflow

**Bitflow** is a flexible TCP proxy. It offers:

* Support for multiple independent concurrent proxy instances
* Support for TCP, Unix domain socket, and inherited file descriptor socket transports for both
  listeners and targets
* Support for TLS termination with SNI-based proxy routing capabilities
* Detailed observability into connection behavior
* Highly configurable connection circuit breaking and timeouts
* Optional use of the HAProxy PROXY protocol header
* IPv6 support and seamless in-proxy IPv4/IPv6 translation

## Building

```bash
$ make
```

## Usage

```bash
$ ./bin/bitflow-$OS-$ARCH --config config.yaml
```

## Configuration

Bitflow requires a configuration file that describes one or more proxies as a listener and target.
An example configuration file is available in `config.example.yaml`.

The following configuration directives are recognized:

|Key|Required|Description|
|-|-|-|
|`application.sentry_dsn`|No|Sentry DSN URL for error reporting (omit to disable)|
|`application.metrics.address`|No|statsd listener address for metrics reporting (omit to disable)|
|`application.metrics.prefix`|No|statsd metrics prefix for all emitted metrics|
|`application.log.level`|No|Application log level; choose from `[debug info warn error]` (omit to disable)|
|`server.proxies[].name`|Yes|Proxy instance name|
|`server.proxies[].listener.address`|Yes|Proxy listener address|
|`server.proxies[].listener.read_timeout`|No|Listener socket read timeout|
|`server.proxies[].listener.write_timeout`|No|Listener socket write timeout|
|`server.proxies[].listener.authorized_sources[]`|No|List of CIDR-notation IP address blocks permitted to establish connections|
|`server.proxies[].listener.tls_context.verify_mode`|No|Mutual TLS authentication mode for client handshakes; choose from `[none relaxed strict]`|
|`server.proxies[].listener.tls_context.certificates[].key`|No|Path to a PEM-format TLS server private key|
|`server.proxies[].listener.tls_context.certificates[].certificate`|No|Path to a PEM-format TLS server certificate|
|`server.proxies[].listener.tls_context.certificates[].authority`|No|Path to a PEM-format TLS CA certificate for client certificate validation|
|`server.proxies[].listener.tls_context.application_protocols[]`|No|List of server-supported application protocols to advertise for ALPN at handshake time|
|`server.proxies[].targets[].name`|Yes|Proxy target name|
|`server.proxies[].targets[].address`|Yes|Proxy target address|
|`server.proxies[].targets[].proxy`|No|SOCKS5 proxy for the target connection|
|`server.proxies[].targets[].connect_backoff`|No|Seed delay for an exponential backoff connection retry policy|
|`server.proxies[].targets[].connect_attempts`|No|Number of times to attempt a connection when connection fails|
|`server.proxies[].targets[].connect_timeout`|No|Timeout for creating a connection to the target|
|`server.proxies[].targets[].resolve_timeout`|No|Timeout for resolving hostnames; used only for dual stack dynamic selection|
|`server.proxies[].targets[].write_timeout`|No|Target socket write timeout|
|`server.proxies[].targets[].read_timeout`|No|Target socket read timeout|
|`server.proxies[].targets[].keepalive_interval`|No|TCP keepalive interval; 0 to use default, -1 to disable, any positive value as a custom interval|
|`server.proxies[].targets[].tls_context.verify_mode`|No|Mutual TLS authentication mode for server handshakes; choose from `[none relaxed strict]`|
|`server.proxies[].targets[].tls_context.certificates[].key`|No|Path to a PEM-format TLS client private key|
|`server.proxies[].targets[].tls_context.certificates[].certificate`|No|Path to a PEM-format TLS client certificate|
|`server.proxies[].targets[].tls_context.certificates[].authority`|No|Path to a PEM-format TLS CA certificate for server certificate validation|
|`server.proxies[].targets[].tls_context.application_protocols[]`|No|List of client-supported application protocols to negotiate for ALPN at handshake time|
|`server.proxies[].options.connection_limit`|No|Maximum allowed active open connections; accepted connections above this limit are rejected|
|`server.proxies[].options.connection_lifetime`|No|Maximum allowed duration of a client connection|
|`server.proxies[].options.connection_log`|No|Optional path to a log file on disk for JSON-formatted connection access logs|
|`server.proxies[].options.connection_load_balancer`|No|Load balancing policy for distributing connections among multiple targets; choose from `[none failover round-robin sni]`|
|`server.proxies[].options.enable_proxy_header`|No|Enable use of the HAProxy PROXY protocol header in upstream connections|
|`server.proxies[].options.enable_dual_stack`|No|Enable use of downstream IP networks for upstream connections|
