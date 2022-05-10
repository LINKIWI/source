# aperture

[![Go Reference](https://pkg.go.dev/badge/lib.kevinlin.info/aperture.svg)](https://pkg.go.dev/lib.kevinlin.info/aperture)

**Aperture** is a Go statsd client library. It offers the following features:

* Native, first-class support for metric tags
* Support for various underlying network transports, including Unix datagram sockets
* Rich networking abstractions including SOCKS-proxied transports and automatic reconnects
* Built-in logging of metrics emissions
* Configurable metric buffering/batching
* Opt-in asynchronous delivery for improved throughput performance
* Pluggable statsd protocol serialization implementations

Aperture doesn't claim to offer better throughput or resource efficiency compared to other statsd
client libraries. Its main goal is to provide common metrics usage abstractions out of the box, to
minimize the amount of code that needs to be implemented on top of the library in applications.

## Usage

To get the library,

```bash
$ go get -d lib.kevinlin.info/aperture
```

Instantiate a statsd client with a `Config` struct to customize configuration parameters.

```go
package main

import (
	"lib.kevinlin.info/aperture"
)

func main() {
	cfg := &aperture.Config{
		Address: "unix:///var/run/metrics/statsd.sock",
		Prefix:  "aperture",
	}

	statsd, err := aperture.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	statsd.Incr("client.instantiate", map[string]interface{}{"aperture": 1})
}
```
