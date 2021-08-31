# zero

**Zero** is a process and socket manager for zero-downtime binary deploys.

It binds listening sockets on TCP or Unix networks and passes the associated file descriptors to
a forked service process. Zero manages the lifecycle of the underlying service process with standard
Unix process signals, allowing for true graceful zero-downtime restarts.

## Building

```bash
$ make
```

## Usage

```bash
$ ./bin/zero-$OS-$ARCH --config config.yaml
```

## Configuration

Zero requires a configuration file that describes the listeners to bind, the runtime parameters for
the child service process, and options for hot-reloading the binary during deployments.
An example configuration file is available in `config.example.yaml`.

One Zero process manages the lifecycle for a single service process. For a multi-tenant environment,
it is expected that there exists multiple Zero instances, one per service instance.
