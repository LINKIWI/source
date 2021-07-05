/**
 * Invert the key-value mapping of an object. Behavior is undefined when values are not unique.
 *
 * @param {Object} obj Arbitrary key-value map with ideally unique values.
 * @returns {Object} Input object with keys and values inverted.
 */
export const invertMap = (obj) =>
  Object.entries(obj).reduce((acc, [key, value]) => ({ ...acc, [value]: key }), {});

export default undefined;
