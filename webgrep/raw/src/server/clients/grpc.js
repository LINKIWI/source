import async from 'async';
import * as grpc from '@grpc/grpc-js';
import { stopwatch } from 'shared/util/instrumentation';

/**
 * Abstract superclass describing a client implementing a per-request load balancing policy among
 * multiple, identical replicas of a gRPC client.
 */
class GRPCLoadBalancer {
  /**
   * Create a GRPCLoadBalancer with an array of backends (gRPC clients).
   *
   * @param {Array} backends Array of GRPCClient instances.
   */
  constructor(backends) {
    this.backends = backends;
  }

  /**
   * Invoke an RPC method. It is expected that subclasses implement this method.
   *
   * @param {String} method Service method name.
   * @param {Object} req Request object.
   * @param {Function} cb RPC stub callback function.
   */
  rpc(method, req, cb) {}  // eslint-disable-line

  /**
   * Execute a single unary RPC on all backend clients, in parallel. Arguments are proxied
   * transparently directly to the backends, skipping all load balancing.
   *
   * @param {String} method Service method name.
   * @param {Object} req Request object.
   * @param {Function} cb RPC stub callback function.
   */
  broadcastRPC(method, req, cb) {
    async.parallel(this.backends.map((backend) => (done) => backend.rpc(method, req, done)), cb);
  }
}

/**
 * API-compliant gRPC load balancer implementation that statically uses the first backend
 * exclusively for all requests.
 */
export class StaticGRPCLoadBalancer extends GRPCLoadBalancer {
  rpc(method, req, cb) {
    return this.backends[0].rpc(method, req, cb);
  }
}

/**
 * Abstraction over multiple gRPC clients that rotates among all clients in round-robin order for
 * each request.
 */
export class RoundRobinGRPCLoadBalancer extends GRPCLoadBalancer {
  rotation = 0;  // Current round robin index

  rpc(method, req, cb) {
    try {
      return this.backends[this.rotation].rpc(method, req, cb);
    } finally {
      this.rotation = (this.rotation + 1) % this.backends.length;
    }
  }
}

/**
 * Abstraction over multiple gRPC clients that prefers earlier clients, failing over to later
 * clients when the preceding client encounters any error during the stub invocation.
 */
export class FailoverGRPCLoadBalancer extends GRPCLoadBalancer {
  rpc(method, req, cb) {
    const failoverRPC = (idx) => this.backends[idx].rpc(method, req, (err, ...args) => {
      if (err) {
        // All backends exhausted; propagate the error back to the client
        if (idx === this.backends.length) {
          return cb(err, ...args);
        }

        return failoverRPC(idx + 1);
      }

      return cb(err, ...args);
    });

    return failoverRPC(0);
  }
}

/**
 * Client for interacting with gRPC services.
 */
export default class GRPCClient {
  /**
   * Create a gRPC service client, with access to an RPC stub as well as the protobuf data
   * structures defined alongside the service.
   *
   * @param {String} name Unique identifier for this gRPC service client instance.
   * @param {String} address Address of the gRPC server.
   * @param {String} authority HTTP/2 authority pseudo-header for requests to the gRPC server.
   * @param {String} service Name of the gRPC service as defined in the proto file.
   * @param {MetricsClient} metrics Metrics client instance for automatic instrumentation of RPC
   *                                stub invocations.
   * @param {Object} options Optional object of customized options to apply.
   */
  constructor(name, address, authority, service, metrics, options = {}) {
    const Service = service;
    const opts = {
      'grpc.default_authority': authority,
      ...options.maxRecvMessageSize && {
        'grpc.max_receive_message_length': options.maxRecvMessageSize,
      },
      ...options.maxSendMessageSize && {
        'grpc.max_send_message_length': options.maxSendMessageSize,
      },
    };

    this._stub = address &&
      service &&
      new Service(address, grpc.credentials.createInsecure(), opts);
    this._name = name;
    this._metrics = metrics;
    this._options = options;
  }

  /**
   * Invoke an RPC method.
   *
   * @param {String} method Service method name.
   * @param {Object} req Request object.
   * @param {Function} cb RPC stub callback function.
   */
  rpc(method, req, cb) {
    const duration = stopwatch();
    const tags = {
      service: this._name.toLowerCase(),
      method,
    };
    const opts = {
      deadline: new Date(Date.now() + this._options.requestTimeout || 3000),
    };

    this._metrics.increment('event.rpc.invoke', tags);

    return this._stub[method](req, opts, (err, ...args) => {
      if (err) {
        this._metrics.increment('event.rpc.error', tags);
      }

      this._metrics.timing('latency.rpc.invoke', duration.elapsed(), tags);

      return cb(err, ...args);
    });
  }
}
