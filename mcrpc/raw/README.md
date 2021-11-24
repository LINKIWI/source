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
	"context"
	"net"

	"lib.kevinlin.info/mcrpc"
	"lib.kevinlin.info/mcrpc/lib"
	"lib.kevinlin.info/mcrpc/protocol"
)

type dummyHandler struct {
	lib.NoopHandler
}

func (d *dummyHandler) Version(ctx context.Context, request *protocol.VersionRequest) (*protocol.VersionResponse, error) {
	return &protocol.VersionResponse{Version: "dummy"}, nil
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

Generate code (requires the [`ragel`](https://www.colm.net/open-source/ragel/) binary):

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
