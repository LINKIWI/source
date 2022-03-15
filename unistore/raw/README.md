# unistore

**Unistore** is an object storage server with pluggable middleware layers and storage backends.
It can be accessed directly over gRPC or via the first-party command line client, `uni`. The gRPC
service schema is defined in `schemas/service/unistore.proto` and can also be obtained over gRPC
reflection.

Unistore can be used as both source of truth storage and as a proxy to other object storage
backends. Supported backends include:

* disk (source of truth storage on a local filesystem)
* Unistore (use another Unistore server as a backend)
* Backblaze B2 (remote cloud storage)

Configurable middleware layers include:

* Bucket remapping/aliasing
* Bucket authorization policies with per-RPC granularity
* Automatic checksum injection (reads) and validation (writes)
* In-flight compression with configurable compression algorithms and levels
* In-flight encryption with configurable encryption mechanisms
* In-flight stream chunk reshaping for optimal performance
* Structured request and response logging

## Building

Building from source requires a Go toolchain.

```bash
$ make -j
```

## Usage

Start a Unistore server:

```bash
$ ./bin/unistore/unistore-$OS-$ARCH --config server.yaml
```

Use the Uni CLI tool:

```bash
$ ./bin/unistore-cli/unistore-cli-$OS-$ARCH --config client.yaml
```

## Configuration

A configuration file is required for both server and client usage. An example configuration file is
available in `config.example.yaml`.
