import Statsd from 'node-statsd';

/**
 * Interface for a backend providing push-based metrics reporting APIs.
 */
class MetricsBackend {
  /**
   * Emit a count metric.
   *
   * @param {String} metric Stat name.
   * @param {Number} value Value by which to increment the stat.
   */
  count(metric, value) {}  // eslint-disable-line

  /**
   * Emit a timing metric.
   *
   * @param {String} metric Stat name.
   * @param {Number} duration Duration of the event, in milliseconds..
   */
  timing(metric, duration) {}  // eslint-disable-line

  /**
   * Emit a gauge metric.
   *
   * @param {String} metric Stat name.
   * @param {Number} value Gauge value.
   */
  gauge(metric, value) {}  // eslint-disable-line
}

/**
 * Metrics backend that noops on all emissions.
 */
export class NoopMetricsClient extends MetricsBackend {}

/**
 * Statsd metrics backend.
 */
export class StatsdMetricsClient extends MetricsBackend {
  /**
   * Create a client instance.
   *
   * @param {String} addr Address of the statsd listener.
   * @param {String} prefix Prefix to apply to every metric.
   * @param {Number} sampleRate Metric emission sample rate.
   */
  constructor(addr, prefix, sampleRate) {
    super();

    const [host, port] = addr.split(':', 2);

    this.sampleRate = sampleRate;
    this.client = new Statsd({
      host,
      port,
      ...prefix && { prefix: `${prefix}.` },
    });
  }

  count(metric, value) {
    this.client.increment(metric, value, this.sampleRate);
  }

  timing(metric, duration) {
    this.client.timing(metric, duration, this.sampleRate);
  }

  gauge(metric, value) {
    this.client.gauge(metric, value, this.sampleRate);
  }
}

/**
 * Generic metrics client backed by a MetricsBackend implementation, with support for InfluxDB-style
 * stat tags.
 */
export default class MetricsClient {
  /**
   * Create a client instance.
   *
   * @param {MetricsBackend} backend Instance of a metrics backend implementation.
   * @param {Object} defaultTags Object of tag to apply to every metric.
   */
  constructor(backend, defaultTags) {
    this.backend = backend;
    this.defaultTags = defaultTags;
  }

  /**
   * Emit a count metric with value 1.
   *
   * @param {String} metric Stat name.
   * @param {Object} tags Optional object of tags to include.
   */
  increment(metric, tags = {}) {
    this.count(metric, 1, tags);
  }

  /**
   * Emit a count metric.
   *
   * @param {String} metric Stat name.
   * @param {Number} value Value by which to increment the stat.
   * @param {Object} tags Optional object of tags to include.
   */
  count(metric, value, tags = {}) {
    this.backend.count(this._formatMetric(metric, tags), value);
  }

  /**
   * Emit a timing metric.
   *
   * @param {String} metric Stat name.
   * @param {Number} duration Duration of the timing, in milliseconds.
   * @param {Object} tags Optional object of tags to include.
   */
  timing(metric, duration, tags = {}) {
    this.backend.timing(this._formatMetric(metric, tags), duration);
  }

  /**
   * Emit a gauge metric.
   *
   * @param {String} metric Stat name.
   * @param {Number} value Current gauge value.
   * @param {Object} tags Optional object of tags to include.
   */
  gauge(metric, value, tags = {}) {
    this.backend.gauge(this._formatMetric(metric, tags), value);
  }

  /**
   * Format a metric with tags, assuming InfluxDB-style tag delimiters.
   *
   * @param {String} metric Stat name.
   * @param {Object} tags Optional object of tags to include.
   * @returns {String} Formatted metric name.
   * @private
   */
  _formatMetric(metric, tags) {
    const combinedTags = { ...this.defaultTags, ...tags };

    if (!Object.keys(combinedTags).length) {
      return metric;
    }

    const serializedTags = Object.entries(combinedTags)
      .filter(([_, value]) => value !== undefined)
      .map(([key, value]) => `${key}=${value}`)
      .join(',');

    return `${metric},${serializedTags}`;
  }
}
