import async from 'async';
import protobufjs from 'protobufjs';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
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

  /**
   * Retrieve the Protobuf definition associated with the backends. This assumes that all backends
   * are homogeneous, and thus provides only the first backend's Protobuf definition.
   *
   * @return {Object} Service Protobuf definition of the first backend.
   */
  get proto() {
    return this.backends[0].proto;
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
   * @param {String} addr Address of the gRPC server.
   * @param {String} authority HTTP/2 authority pseudo-header for requests to the gRPC server.
   * @param {String} service Name of the gRPC service as defined in the proto file.
   * @param {String} protoPath Path to the protobuf service definition on disk.
   * @param {MetricsClient} metrics Metrics client instance for automatic instrumentation of RPC
   *                                stub invocations.
   * @param {Object} options Optional object of customized options to apply.
   */
  constructor(addr, authority, service, protoPath, metrics, options = {}) {
    const packageDefinition = protoLoader.loadSync(protoPath, {
      longs: String,
      defaults: true,
      oneofs: true,
    });
    const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
    const opts = { 'grpc.default_authority': authority };
    const Service = service && protoDescriptor[service];

    this._proto = protobufjs.loadSync(protoPath);
    this._stub = addr && service && new Service(addr, grpc.credentials.createInsecure(), opts);
    this._service = service;
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
      service: this._service.toLowerCase(),
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

  /**
   * Object representation of the protobuf's defined data structures.
   *
   * @returns {Object} Protobuf data structures.
   */
  get proto() {
    return this._proto;
  }
}
