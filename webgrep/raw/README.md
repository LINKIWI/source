# webgrep

webgrep is a web client for livegrep.

## Building

To build a standalone `webgrep` binary,

```
$ npm run build
$ npm run release
$ ./dist/bin/webgrep --help
```

## Configuration

See `config.example.yaml` for an example configuration file. The config file location can be specified at runtime with the `--config` flag.

The following documents each field and its expected value:

|Key|Required|Description|
|-|-|-|
|`server.listen.transport`|Yes|Network transport on which to bind the server HTTP listener; one of `tcp` or `unix`|
|`server.listen.address`|Yes|IP address and port pair for TCP transports, or a socket path for Unix domain socket transports|
|`server.livegrep.load_balancing_policy`|No|Load balancing policy to use among multiple livegrep replicas; one of `round-robin` (rotate through all servers on each request), `failover` (prefer earlier servers, failing over to next servers on request errors) (default exclusive use of first server with no load balancing)|
|`server.livegrep.request_timeout`|No|Per-request RPC invocation timeout in milliseconds (default 3000)|
|`server.livegrep.servers[].authority`|No|Livegrep codesearch gRPC server HTTP/2 authority header|
|`server.livegrep.servers[].address`|Yes|Livegrep codesearch gRPC server address|
|`server.metrics.statsd.addr`|No|Address of the statsd server or listener for metrics reporting (default disabled)|
|`server.metrics.statsd.prefix`|No|String prefix for all emitted statsd metrics (default `webgrep`)|
|`server.metrics.statsd.sample_rate`|No|statsd reporting sample rate (default `1.0`)|
|`server.logging.winston.level`|No|Winston console output log level (default `info`)|
|`server.logging.winston.output`|No|Winston console output file descriptor (default `stdout`)|
|`server.logging.supercharged.enabled`|No|Enable Supercharged-internal framework logging (default `false`)|
|`server.cache.redis.ttl`|No|TTL for Redis cache keys in milliseconds (default 10 minutes)|
|`server.cache.redis.timeout`|No|Timeout for all Redis requests (default disabled)|
|`server.cache.redis.prefix`|No|Prefix to apply to Redis cache keys (default `webgrep`)|
|`server.cache.redis.socket`|No|Redis Unix socket path|
|`server.cache.memory.ttl`|No|TTL for in-memory cache keys in milliseconds (default 10 minutes)|
|`server.source.local.repositories[].name`|No|Local on-disk repository name for source code viewer|
|`server.source.local.repositories[].path`|No|Local on-disk repository directory path|
|`server.source.gitlab.base_url`|No|GitLab instance base URL for source code viewer (default `https://gitlab.com`)|
|`server.source.gitlab.socket_path`|No|Path to a Unix domain socket for accessing GitLab API|
|`server.source.gitlab.access_token`|No|GitLab API access token|
|`server.source.gitlab.tls_key`|No|TLS client key path for mutual authentication with GitLab|
|`server.source.gitlab.tls_cert`|No|TLS client certificate path for mutual authentication with GitLab|
|`server.source.gitlab.tls_ca_cert`|No|TLS CA certificate path for server authentication of GitLab|
|`client.site.logo`|No|URL to an image to use as the logo shown in the upper-left corner|
|`client.site.title`|No|Formatting string for creating browser page titles (default `webgrep`)|
|`client.site.about`|No|About text shown in the Admin section|
|`client.site.banner.title`|No|Header informational banner title|
|`client.site.banner.description`|No|Header informational banner text description|
|`client.site.banner.link`|No|Header informational banner external link title|
|`client.site.banner.href`|No|Header informational banner external link target URL|
|`client.search.name`|No|Search engine plugin name|
|`client.search.description`|No|Search engine plugin description|
|`client.search.base_url`|No|Instance base URL used by the search engine plugin for search executions|
|`client.options.metadata_poll_interval`|No|Poll interval for new server metadata in milliseconds|
|`client.resources[].title`|No|Title for an arbitrary external link shown in the Admin section|
|`client.resources[].href`|No|URL for an associated external link shown in the Admin section|
