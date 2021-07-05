# proton

[![Go Reference](https://pkg.go.dev/badge/lib.kevinlin.info/proton.svg)](https://pkg.go.dev/lib.kevinlin.info/proton)

**Proton** is a Go client library for HTTP servers compliant with the Supercharged API framework.
It supports:

* Requests over both HTTP and bidirectional websockets
* Pluggable `http.Client` backends
* Transparent instrumentation with an [Aperture](https://lib.kevinlin.info/aperture) metrics client

## Usage

To get the library,

```bash
$ go get -d lib.kevinlin.info/proton
```

Instantiate a Proton client instance with a `proton.Config` struct and the `proton.NewConfig`
wrapper function:

```go
package main

import (
	"net/url"

	"lib.kevinlin.info/proton"
)

type exampleResponse struct {
	Value string `json:"value"`
	Num   int    `json:"num"`
}

func main() {
	cfg, err := proton.NewConfig(&proton.Config{
		BaseURL: &url.URL{
			Scheme: "https",
			Host:   "supercharged.example.com",
		},
		ClientID: "client",
	})
	if err != nil {
		panic(err)
	}

	client := proton.NewClient(cfg)

	var resp exampleResponse

	if err := client.DoHTTP("GET", "/endpoint", nil, &resp); err != nil {
		panic(err)
	}
}

```
