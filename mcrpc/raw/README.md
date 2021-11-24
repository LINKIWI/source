# mcrpc

**mcrpc** is a library for building [memcache protocol](https://github.com/memcached/memcached/blob/master/doc/protocol.txt)
servers with arbitrary backing implementations.

## Usage

Add the library module dependency:

```bash
$ go get lib.kevinlin.info/mcrpc
```

Define a server handler and start a server:

```go
package main

import (
	"net"

	"lib.kevinlin.info/mcrpc"
	"lib.kevinlin.info/mcrpc/lib"
)

type dummyHandler struct {
	lib.NoopHandler
}

func main() {
	ln, err := net.Listen("tcp", "localhost:11211")
	if err != nil {
		panic(err)
	}

	srv, err := mcrpc.NewServer(&dummyHandler{})
	if err != nil {
		panic(err)
	}

	srv.Serve(ln)
}

```

## Development

Generate code:

```bash
$ make generate
```

Build the library:

```bash
$ make build
```

Run tests:

```bash
$ make lint
$ make test
```
