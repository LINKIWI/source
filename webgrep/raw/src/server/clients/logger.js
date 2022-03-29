import Logger from 'supercharged/server/logging';
import winston from 'winston';

/**
 * Interface for a backend for leveled logging.
 */
class LoggerBackend {
  /**
   * Log at a debug level.
   *
   * @param {String} args Variable arguments.
   */
  debug(...args) {}  // eslint-disable-line

  /**
   * Log at an info level.
   *
   * @param {String} args Variable arguments.
   */
  info(...args) {}  // eslint-disable-line

  /**
   * Log at a warn level.
   *
   * @param {String} args Variable arguments.
   */
  warn(...args) {}  // eslint-disable-line

  /**
   * Log at an error level.
   *
   * @param {String} args Variable arguments.
   */
  error(...args) {}  // eslint-disable-line
}

/**
 * Logger implementation that noops on all logs.
 */
export class NoopLogger extends LoggerBackend {}

/**
 * Logger implementation backed by the Winston logging library.
 */
export class WinstonLogger extends LoggerBackend {
  /**
   * Create a logger backed by Winston.
   *
   * @param {Object} opts Object of Winston logger options.
   */
  constructor(opts) {
    super();

    this.logger = winston.createLogger(opts);
  }

  debug(...args) {
    return this.logger.debug(...args);
  }

  info(...args) {
    return this.logger.info(...args);
  }

  warn(...args) {
    return this.logger.warn(...args);
  }

  error(...args) {
    return this.logger.error(...args);
  }
}

/**
 * Top-level abstraction for a logger with levels, namespaces, and tags.
 */
export default class LoggerClient extends Logger {
  /**
   * Create a logger client instance.
   *
   * @param {LoggerBackend} backend Instance of an implementation of a logger backend.
   */
  constructor(backend) {
    super();

    this.backend = backend;
  }

  debug(namespace, message, tags = {}) {
    return this._log('debug', namespace, message, tags);
  }

  info(namespace, message, tags = {}) {
    return this._log('info', namespace, message, tags);
  }

  warn(namespace, message, tags = {}) {
    return this._log('warn', namespace, message, tags);
  }

  error(namespace, message, tags = {}) {
    return this._log('error', namespace, message, tags);
  }

  _log(level, namespace, message, tags) {
    const logger = this.backend[level].bind(this.backend);

    if (!Object.keys(tags || {})) {
      return logger(`${namespace}: ${message}`);
    }

    const formattedTags = Object.entries(tags)
      .map(([key, value]) => ((typeof value === 'object' || Array.isArray(value)) ?
        `${key}=${JSON.stringify(value)}` :
        `${key}=${value}`))
      .join(' ');

    return logger(`${namespace}: ${message}: ${formattedTags}`);
  }
}
