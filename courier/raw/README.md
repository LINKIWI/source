# courier

**Courier** is an HTTP service proxy appropriate for use both as an internet-facing ingress reverse
proxy and as a service-to-service mesh proxy. It offers:

* Composable configuration based on pluggable filters
* Rich, extensible metrics and logging
* Support for HTTP/2 transports, including gRPC
* TLS termination and origination with support for mutual authentication
* Support for forward proxies behind reverse proxy upstreams

## Building

```bash
$ make
```

## Usage

Start a Courier server:

```bash
$ ./bin/courier-$OS-$ARCH --config config.yaml
```

Validate a configuration file:

```bash
$ ./bin/courier-$OS-$ARCH --config config.yaml --validate
```

## Configuration

Courier requires a configuration file that describes server listeners and virtual hosts. Each
virtual host is characterized by one upstream and an optional chain of filters that inspect or
modify proxied requests and responses.

An example configuration file with documentation for all supported directives is available in
`config.example.yaml`.
