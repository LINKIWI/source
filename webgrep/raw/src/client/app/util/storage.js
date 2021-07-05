/**
 * Abstract superclass for a persistent key-value store.
 */
class KeyValueStore {
  /**
   * Create a key-value store.
   *
   * @param {String} prefix Optional key prefix.
   * @param {Function} serializer Unary function to use as the value serializer.
   * @param {Function} deserializer Unary function to use as the value deserializer.
   */
  constructor(prefix = 'webgrep', serializer = JSON.stringify, deserializer = JSON.parse) {
    this.prefix = prefix;
    this.serializer = serializer;
    this.deserializer = deserializer;
  }

  /**
   * Format a string key to use for insertion into the store.
   *
   * @param {String} namespace Key namespace; user-defined.
   * @param {String} key Key itself; user-defined.
   * @return {string} Formatted single string to use as the identifying key in the backing store.
   * @private
   */
  _formatKey(namespace, key) {
    return `${this.prefix}:${namespace}:${key}`;
  }
}

/**
 * Key-value store backed by a browser's local storage engine.
 */
class LocalStorageClient extends KeyValueStore {
  /**
   * Get a value from the store. Throws on error.
   *
   * @param {String} namespace Key namespace.
   * @param {String} key Key itself.
   * @return {Object|null} Corresponding value, or null if the key does not exist.
   */
  get(namespace, key) {
    const value = window.localStorage.getItem(this._formatKey(namespace, key));
    if (value === null) {
      return null;
    }

    return this.deserializer(value);
  }

  /**
   * Set a value to the store. Throws on error.
   *
   * @param {String} namespace Key namespace.
   * @param {String} key Key itself.
   * @param {Object} data Payload to persist.
   */
  set(namespace, key, data) {
    window.localStorage.setItem(this._formatKey(namespace, key), this.serializer(data));
  }

  /**
   * Delete a value from the store by key. Throws on error.
   *
   * @param {String} namespace Key namespace.
   * @param {String} key Key itself.
   */
  delete(namespace, key) {
    window.localStorage.removeItem(this._formatKey(namespace, key));
  }
}

/**
 * Key-value store backend by an in-memory (transient) map.
 */
class MemoryStorageClient extends KeyValueStore {
  store = {};

  /**
   * Get a value from the store.
   *
   * @param {String} namespace Key namespace.
   * @param {String} key Key itself.
   * @return {Object|null} Corresponding value, or null if the key does not exist.
   */
  get(namespace, key) {
    const value = this.store[this._formatKey(namespace, key)];
    if (value === undefined) {
      return undefined;
    }

    return this.deserializer(value);
  }

  /**
   * Set a value to the store.
   *
   * @param {String} namespace Key namespace.
   * @param {String} key Key itself.
   * @param {Object} data Payload to persist.
   */
  set(namespace, key, data) {
    this.store[this._formatKey(namespace, key)] = this.serializer(data);
  }

  /**
   * Delete a value from the store by key. Throws on error.
   *
   * @param {String} namespace Key namespace.
   * @param {String} key Key itself.
   */
  delete(namespace, key) {
    delete this.store[this._formatKey(namespace, key)];
  }
}

/**
 * Storage client that abstracts out the backing storage implementation. Uses the browser's local
 * storage API if available, but otherwise falls back to a transient, in-memory store. Assumes JSON
 * as the serialization/deserialization mechanism.
 */
export default class StorageClient {
  constructor(...args) {
    this.memory = new MemoryStorageClient(...args);
    this.persistent = new LocalStorageClient(...args);
  }

  /**
   * Get a value from the store.
   *
   * @param {String} namespace Key namespace.
   * @param {String} key Key itself.
   * @return {Object||null} Corresponding value, or null if the key does not exist.
   */
  get(namespace, key) {
    try {
      return this.persistent.get(namespace, key);
    } catch (e) {
      return this.memory.get(namespace, key);
    }
  }

  /**
   * Set a value to the store.
   *
   * @param {String} namespace Key namespace.
   * @param {String} key Key itself.
   * @param {Object} data Payload to persist.
   */
  set(namespace, key, data) {
    try {
      return this.persistent.set(namespace, key, data);
    } catch (e) {
      return this.memory.set(namespace, key, data);
    }
  }

  /**
   * Delete a value from the store by key. Throws on error.
   *
   * @param {String} namespace Key namespace.
   * @param {String} key Key itself.
   */
  delete(namespace, key) {
    try {
      return this.persistent.delete(namespace, key);
    } catch (e) {
      return this.memory.delete(namespace, key);
    }
  }
}

// Static singleton instance for general-purpose, shared use cases.
export const sharedStorage = new StorageClient();
