import { URLStateSerializer, URLStateDeserializer } from 'client/app/util/data';

/**
 * Higher-order function for injecting preventDefault() to an event callback.
 *
 * @param {Function} cb Event callback.
 * @return {Function} Function exposing the same API as an event callback, but calls
 *                    evt.preventDefault() before invoking the wrapped logic.
 */
export const withDefaultPrevented = (cb) => (evt, ...args) => {
  evt.preventDefault();
  return cb(evt, ...args);
};

/**
 * Encode arbitrary single-level state into the URL.
 *
 * @param {Object} state Arbitrary JSON-serializable description of the state.
 * @param {Object} filters Map of state keys to filter functions, to describe whether the key-value
 *                         pair should be included in URL state.
 * @param {Object} serializers Map of state keys to serialization functions.
 */
export const encodeURLState = (state, filters, serializers) => {
  // Explicitly clear the URL if the state is empty
  if (!state || !Object.keys(state).length) {
    window.history.replaceState('', '', window.location.pathname);
    return;
  }

  const qs = Object.entries(state)
    .filter(([key, value]) => (key in filters ? filters[key](value) : true))
    .map(([key, value]) => `${key}=${(serializers[key] || URLStateSerializer.identity)(value)}`)
    .join('&');

  window.history.replaceState('', '', `?${qs}`);

  // Encoding the state into the URL will remove the hash if present.
  // However, this does not trigger a hashchange event; manually trigger it here.
  window.dispatchEvent(new HashChangeEvent('hashchange'));
};

/**
 * Parse the state object encoded into the URL.
 *
 * @param {Object} defaults Object of defaults. Only keys in this object may be overridden by values
 *                          parsed from the URL querystring.
 * @param {Object} deserializers Map of state keys to deserialization functions.
 * @return {Object} Object describing the JSON-parsed state from the URL.
 */
export const decodeURLState = (defaults, deserializers) =>
  [...new URLSearchParams(window.location.search).entries()]
    // Consider only those keys that are present in the defaults object; this prevents the decoded
    // state from being polluted with extraneous query parameters
    .filter(([key]) => key in defaults)
    .reduce((acc, [key, value]) => {
      const data = (deserializers[key] || URLStateDeserializer.identity)(value);

      return {
        ...acc,
        ...value && {
          [key]: data,
        },
      };
    }, defaults);

/**
 * Scroll the viewpoint to a target coordinate.
 *
 * @param {Number} top Top pixel coordinate.
 * @param {Number} left Left pixel coordinate.
 */
export const scroll = (top = 0, left = 0) => window.scrollTo({
  top,
  left,
  behavior: 'smooth',
});
