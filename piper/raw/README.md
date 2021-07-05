# piper

Piper is a daemon for realtime streaming of log files to Kafka.

### Building

```bash
$ make
```

### Usage

```bash
$ ./bin/piper-$OS-$ARCH --config config.toml
```

### Configuration

Every Piper instance must be launched with a valid configuration file that describes the log files
to watch and their associated Kafka targets.

An example configuration file with documentation for all supported directives is available in
`config.example.toml`.
