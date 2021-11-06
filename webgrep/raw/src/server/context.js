import os from 'os';
import path from 'path';
import winston from 'winston';
import CacheClient, { MemoryCache, NoopCache, RedisCache } from 'server/clients/cache';
import ConfigClient from 'server/clients/config';
import GRPCClient, { RoundRobinGRPCLoadBalancer, FailoverGRPCLoadBalancer, StaticGRPCLoadBalancer } from 'server/clients/grpc';
import MetricsClient, { NoopMetricsClient, StatsdMetricsClient } from 'server/clients/metrics';
import SourceClient, { NoopSourceBackend, LocalSourceBackend, GitlabSourceBackend } from 'server/clients/source';
import AdminLogic from 'server/logic/admin';
import MetaLogic from 'server/logic/meta';
import SearchLogic from 'server/logic/search';
import SourceLogic from 'server/logic/source';
import ViewLogic from 'server/logic/view';

/**
 * Format a fully qualified path to a schemas file.
 *
 * @param {String} schema Name of the schemas filename.
 * @return {String} Fully qualified path to the schemas path.
 */
const schemasPath = (schema) => path.join(__dirname, `../shared/schemas/${schema}`);

/**
 * Server-side context shared by all incoming request handlers.
 */
export default class Context {
  constructor(options) {
    this.config = new ConfigClient(options.config);

    this.metrics = new MetricsClient(
      (() => {
        const backends = Object.keys(this.config.get('server.metrics') || {});
        if (!backends.length) {
          return new NoopMetricsClient();
        }

        const [backend] = backends;

        switch (backend) {
          case 'statsd':
            return new StatsdMetricsClient(
              this.config.get('server.metrics.statsd.addr') || '',
              this.config.get('server.metrics.statsd.prefix') || 'webgrep',
              this.config.get('server.metrics.statsd.sample_rate') || 1.0,
            );
          default:
            return new NoopMetricsClient();
        }
      })(),
      { host: os.hostname(), version: process.env.VERSION },
    );

    this.cache = (() => {
      const noop = new CacheClient(
        new NoopCache(),
        '',
        0,
        new MetricsClient(new NoopMetricsClient(), {}),
      );

      const backends = Object.keys(this.config.get('server.cache') || {});
      if (!backends.length) {
        return noop;
      }

      const [backend] = backends;

      switch (backend) {
        case 'redis':
          return new CacheClient(
            new RedisCache(
              this.config.get('server.cache.redis.socket'),
              this.config.get('server.cache.redis.timeout') || 0,
            ),
            this.config.get('server.cache.redis.prefix') || 'webgrep',
            this.config.get('server.cache.redis.ttl') || 10 * 60 * 1000,
            this.metrics,
          );
        case 'memory':
          return new CacheClient(
            new MemoryCache(),
            '',
            this.config.get('server.cache.memory.ttl') || 10 * 60 * 1000,
            this.metrics,
          );
        default:
          return noop;
      }
    })();

    this.service = {
      codesearch: (() => {
        const backends = this.config.get('server.livegrep.servers').map((server) => new GRPCClient(
          server.address,
          server.authority || '',
          'CodeSearch',
          schemasPath('livegrep.proto'),
          this.metrics,
          { requestTimeout: this.config.get('server.livegrep.request_timeout') || 3000 },
        ));

        switch (this.config.get('server.livegrep.load_balancing_policy')) {
          case 'round-robin':
            return new RoundRobinGRPCLoadBalancer(backends);
          case 'failover':
            return new FailoverGRPCLoadBalancer(backends);
          default:
            return new StaticGRPCLoadBalancer(backends);
        }
      })(),
    };

    this.source = (() => {
      const noop = new SourceClient(
        new NoopSourceBackend(),
        new MetricsClient(new NoopMetricsClient(), {}),
      );

      const backends = Object.keys(this.config.get('server.source') || {});
      if (!backends.length) {
        return noop;
      }

      const [backend] = backends;

      switch (backend) {
        case 'local':
          return new SourceClient(
            new LocalSourceBackend(this.config.get('server.source.local.repositories') || []),
            this.metrics,
          );
        case 'gitlab':
          return new SourceClient(
            new GitlabSourceBackend(
              this.config.get('server.source.gitlab.base_url') || 'https://gitlab.com',
              this.config.get('server.source.gitlab.socket_path'),
              this.config.get('server.source.gitlab.access_token'),
              {
                key: this.config.get('server.source.gitlab.tls_key'),
                cert: this.config.get('server.source.gitlab.tls_cert'),
                ca: this.config.get('server.source.gitlab.tls_ca_cert'),
              },
            ),
            this.metrics,
          );
        default:
          return noop;
      }
    })();

    this.logic = {
      admin: new AdminLogic(this),
      meta: new MetaLogic(this),
      search: new SearchLogic(this),
      source: new SourceLogic(this),
      view: new ViewLogic(this),
    };

    this.log = winston.createLogger({
      level: options.verbosity,
      format: winston.format.combine(
        winston.format.timestamp(),
        winston.format.splat(),
        winston.format.printf((fmt) =>
          `${fmt.timestamp} ${fmt.level.toUpperCase()}\t${fmt.message}`),
      ),
      transports: [
        new winston.transports.Console({
          timestamp: () => Date.now(),
        }),
      ],
    });
  }
}
