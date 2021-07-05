import async from 'async';
import crypto from 'crypto';
import Redis from 'ioredis';
import snappy from 'snappyjs';
import { withTimeout } from 'server/util/instrumentation';
import { stopwatch } from 'shared/util/instrumentation';

/**
 * Interface for a backend providing R/W APIs to a (potentially transient) key-value store.
 */
class CacheBackend {
  /**
   * Perform a key read operation.
   *
   * @param {String} key Cache key.
   * @param {Function} cb Callback invoked with (null, data) on completion.
   */
  get(key, cb) {}  // eslint-disable-line

  /**
   * Write a key to memory.
   *
   * @param {String} key Cache key.
   * @param {String} value Associated cache value.
   * @param {Number} ttl Time to live, in milliseconds.
   * @param {Function} cb Callback invoked with (null) on completion.
   */
  set(key, value, ttl, cb) {}  // eslint-disable-line

  /**
   * Delete a key. This is a noop if the key does not exist.
   *
   * @param {String} key Cache key.
   * @param {Function} cb Function invoked with (null) on completion.
   */
  delete(key, cb) {}  // eslint-disable-line

  /**
   * Read a list of keys matching the specified prefix.
   *
   * @param {String} prefix Key prefix.
   * @param {Function} cb Callback invoked with (null, keys) on completion.
   */
  keys(prefix, cb) {}  // eslint-disable-line

  /**
   * Close the underlying client connection, as necessary.
   */
  close() {}  // eslint-disable-line
}

/**
 * Cache backend implementation that effectively disables cache by reporting all reads as misses.
 */
export class NoopCache extends CacheBackend {
  get(key, cb) {  // eslint-disable-line class-methods-use-this
    return cb(null, null);
  }

  set(key, value, ttl, cb) {  // eslint-disable-line class-methods-use-this
    return cb(null);
  }

  delete(key, cb) {  // eslint-disable-line class-methods-use-this
    return cb(null);
  }

  keys(prefix, cb) {  // eslint-disable-line class-methods-use-this
    return cb(null, []);
  }
}

/**
 * Simple in-memory key-value store with TTL support.
 */
export class MemoryCache extends CacheBackend {
  /**
   * Create an in-memory cache instance. Data is not persisted.
   */
  constructor() {
    super();

    this.store = {};
  }

  get(key, cb) {
    if (!(key in this.store)) {
      return cb(null, null);
    }

    const { value, expiry } = this.store[key];

    if (expiry < 0 || expiry > Date.now()) {
      return cb(null, value);
    }

    // Key expired: remove from backing store
    delete this.store[key];
    return cb(null, null);
  }

  set(key, value, ttl, cb) {
    this.store[key] = {
      value,
      expiry: ttl > 0 ? Date.now() + ttl : -1,
    };

    return cb(null);
  }

  delete(key, cb) {
    if (key in this.store) {
      delete this.store[key];
    }

    return cb(null);
  }

  keys(prefix, cb) {
    return cb(null, Object.entries(this.store)
      .filter(([key]) => key.startsWith(prefix))
      .filter(([_, { expiry }]) => expiry > Date.now())
      .map(([key]) => key));
  }
}

/**
 * Cache backend implementation backed by Redis.
 */
export class RedisCache extends CacheBackend {
  /**
   * Create a Redis cache client instance.
   *
   * @param {String} sock Path to the Redis Unix socket on disk.
   * @param {Number} timeout Timeout to use for all operations; 0 to disable.
   */
  constructor(sock, timeout) {
    super();

    this.timeoutWrap = withTimeout(timeout);
    this.client = new Redis({
      path: sock,
      // Enable fail-fast behavior; return immediately with an error when Redis is unavailable
      enableOfflineQueue: false,
    });
  }

  get(key, cb) {
    return this.timeoutWrap(this.client.getBuffer.bind(this.client), [key], cb);
  }

  set(key, value, ttl, cb) {
    return this.timeoutWrap(this.client.set.bind(this.client), [key, value, 'EX', ttl / 1000], cb);
  }

  delete(key, cb) {
    return this.timeoutWrap(this.client.del.bind(this.client), [key], cb);
  }

  keys(prefix, cb) {
    return this.timeoutWrap(this.client.keys.bind(this.client), [`${prefix}*`], cb);
  }

  close() {
    return this.client.quit();
  }
}

/**
 * Client abstraction for a key-value cache.
 */
export default class CacheClient {
  /**
   * Create a client instance.
   *
   * @param {CacheBackend} backend Instance of an implementation of a cache backend.
   * @param {String} prefix Prefix to apply to all cache keys.
   * @param {Number} defaultTTL Default TTL to apply to set requests, in milliseconds.
   * @param {MetricsClient} metrics Metrics client instance for automatic instrumentation of cache
   *                                requests.
   */
  constructor(backend, prefix, defaultTTL, metrics) {
    this.backend = backend;
    this.prefix = prefix;
    this.defaultTTL = defaultTTL;
    this.metrics = metrics;
  }

  /**
   * Get the value associated with a key.
   *
   * @param {String} namespace Namespace of the key.
   * @param {String} key The unique identifier of the key itself.
   * @param {Object} tags Object of key-value pairs describing tags for the key.
   * @param {Function} cb Callback invoked with (err, data) on completion.
   */
  get(namespace, key, tags, cb) {
    const cacheKey = this._formatKey(namespace, key, tags);
    const duration = stopwatch();

    this.metrics.increment('event.cache.get.request', { namespace });

    return this.backend.get(cacheKey, (err, cached, ...args) => {
      if (err) {
        this.metrics.increment('event.cache.get.error', { namespace });
        return cb(err);
      }

      this.metrics.timing('latency.cache.get', duration.elapsed(), { namespace });

      if (cached === null || cached === undefined) {
        this.metrics.increment('event.cache.get.miss', { namespace });
        return cb(err, cached, ...args);
      }

      this.metrics.increment('event.cache.get.hit', { namespace });
      return cb(err, snappy.uncompress(cached), ...args);
    });
  }

  /**
   * Set a value associated with a key.
   *
   * @param {String} namespace Namespace of the key.
   * @param {String} key The unique identifier of the key itself.
   * @param {String} value String value to set.
   * @param {Object} tags Object of key-value pairs describing tags for the key.
   * @param {Number} ttl Time to live, in milliseconds.
   * @param {Function} cb Callback invoked with (err, data) on completion.
   */
  set(namespace, key, value, tags, ttl, cb = () => {}) {
    const cacheKey = this._formatKey(namespace, key, tags);
    const buf = snappy.compress(Buffer.from(value));
    const duration = stopwatch();

    this.metrics.increment('event.cache.set.request', { namespace });
    this.metrics.gauge('gauge.cache.value_size', buf.length, { namespace });

    return this.backend.set(cacheKey, buf, this.defaultTTL || ttl, (err, ...args) => {
      if (err) {
        this.metrics.increment('event.cache.set.error', { namespace });
        return cb(err, ...args);
      }

      this.metrics.timing('latency.cache.set', duration.elapsed(), { namespace });

      return cb(err, ...args);
    });
  }

  /**
   * Delete a key.
   *
   * @param {String} namespace Namespace of the key.
   * @param {String} key The unique identifier of the key itself.
   * @param {Object} tags Object of key-value pairs describing tags for the key.
   * @param {Function} cb Callback invoked with (err, data) on completion.
   */
  delete(namespace, key, tags, cb = () => {}) {
    const cacheKey = this._formatKey(namespace, key, tags);

    this.metrics.increment('event.cache.delete.request', { namespace });

    return this.backend.delete(cacheKey, cb);
  }

  /**
   * Flush all keys under a namespace.
   *
   * @param {String} namespace Namespace of the keys.
   * @param {Function} cb Callback invoked with (err) on completion.
   */
  flush(namespace, cb = () => {}) {
    return this.backend.keys(`${this.prefix}:${namespace}:`, (err, keys) => {
      if (err) {
        return cb(err);
      }

      return async.parallel(keys.map((key) => (done) => {
        this.metrics.increment('event.cache.delete.request', { namespace });
        this.backend.delete(key, done);
      }), cb);
    });
  }

  /**
   * Create a R/W transaction closure for the same namespace, key, and tags.
   *
   * @param {String} namespace Namespace of the key.
   * @param {String} key The unique identifier of the key itself.
   * @param {Object} tags Object of key-value pairs describing tags for the key.
   * @param {Function} serializer Write-path value serializer.
   * @param {Function} deserializer Read-path value deserializer.
   * @return {Object} Object with get, set, and delete methods for the computed key.
   */
  transaction(namespace, key, tags, serializer, deserializer) {
    return {
      get: (cb) => this.get(namespace, key, tags, ((err, cached, ...args) =>
        (err ? cb(err, cached, ...args) : cb(null, cached && deserializer(cached), ...args)))),
      set: (value, ttl, cb) => this.set(namespace, key, value && serializer(value), tags, ttl, cb),
      delete: (cb) => this.delete(namespace, key, tags, cb),
    };
  }

  /**
   * End the client connection.
   */
  close() {
    return this.backend.close();
  }

  /**
   * Format a cache key for insertion.
   *
   * @param {String} namespace Namespace of the key.
   * @param {String} key The unique identifier of the key itself.
   * @param {Object} tags Object of key-value pairs describing tags for the key.
   * @return {String} Formatted cache key sent to Redis.
   * @private
   */
  _formatKey(namespace, key, tags) {
    const serializedTags = Object.entries(tags)
      .sort(([a], [b]) => a.localeCompare(b))
      .map(([tagKey, tagValue]) => `${tagKey}=${tagValue}`)
      .join('&');

    const hashedKey = crypto
      .createHash('sha1')
      .update(`${namespace}:${key}:${serializedTags}`)
      .digest('hex');

    // Explicitly skip hashing of the prefix, since its value may be used at routing layer
    return `${this.prefix}:${hashedKey}`;
  }
}
