import { CODE_SERVER_UNDEFINED } from 'supercharged/shared/constants/error';
import { objLookup } from 'shared/util/data';
import { stopwatch } from 'shared/util/instrumentation';

/**
 * Wrap a callback-style function with an operation timeout. This higher-order function will invoke
 * the wrapped function. If the original callback is not invoked within a timeout, the callback is
 * prematurely invoked with a timeout error, and the original callback is prevented from future
 * invocation.
 *
 * @param {Number} timeout Timeout in milliseconds; 0 to disable (offer unlimited timeout).
 * @return {Function} Function Function accepting a function, an array of arguments to pass to the
 *                    function, and its callback.
 */
export const withTimeout = (timeout) => (func, args, cb) => {
  let exceededTimeout = false;

  if (timeout <= 0) {
    return func(...args, cb);
  }

  const timeoutID = setTimeout(() => {
    exceededTimeout = true;
    cb(new Error('operation exceeded timeout'));
  }, timeout);

  return func(...args, (...cbArgs) => {
    // Only invoke the original callback if the timeout was not exceeded.
    if (!exceededTimeout) {
      // Clear the timeout timer so that the timeout error callback is not invoked after the
      // original request completes successfully.
      clearTimeout(timeoutID);
      cb(...cbArgs);
    }
  });
};

/**
 * Transparently decorator a handler method with basic metrics instrumentation.
 *
 * @param {Object} target Method target.
 * @param {String} key Name of the decorated method.
 * @param {Object} descriptor Method descriptor.
 * @return {Object} Wrapped method descriptor.
 */
export const withEndpointInstrumentation = (target, key, descriptor) => {
  const wrappedFunc = descriptor.value;

  // eslint-disable-next-line no-param-reassign
  descriptor.value = function reportEndpointMetrics(...args) {
    const duration = stopwatch();
    const metricTags = {
      handler: target.constructor.name,
      method: key,
    };

    this.ctx.metrics.increment('event.api.request', metricTags);

    /* Intercept request completion callbacks in order to record end-to-end latency timers */

    const success = this.success.bind(this);
    const error = this.error.bind(this);

    this.success = ((data) => {
      success(data);
      this.success = success;
      this.error = error;
      this.ctx.metrics.timing('latency.api.request', duration.elapsed(), {
        ...metricTags,
        status: objLookup(data, ['status']) || 200,
      });
    });

    this.error = ((data) => {
      error(data);
      this.success = success;
      this.error = error;
      this.ctx.metrics.timing('latency.api.request', duration.elapsed(), {
        ...metricTags,
        status: objLookup(data, ['status']) || 500,
        code: objLookup(data, ['code']) || CODE_SERVER_UNDEFINED,
      });
    });

    return wrappedFunc.call(this, ...args);
  };

  return descriptor;
};
