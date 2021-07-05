/**
 * Safely parse a JSON string, without throwing an exception on failure.
 *
 * @param {String} stringified String JSON to parse.
 * @return {Object} Object with property `data` for the parsed object and boolean `ok` indicating
 *                  whether the parse was successful.
 */
export const parseJSON = (stringified) => {
  try {
    return { data: JSON.parse(stringified), ok: true };
  } catch (e) {
    return { data: undefined, ok: false };
  }
};

/**
 * Decode a base64-encoded string payload.
 *
 * @param {String} encoded Base64-encoded data.
 * @returns {Object} Object with property `data` for the decoded payload and boolean `ok` indicating
 *                   whether the decode was successful.
 */
export const decodeBase64 = (encoded) => {
  try {
    return { data: window.atob(encoded), ok: true };
  } catch (e) {
    return { data: undefined, ok: false };
  }
};

/**
 * Safely retrieve the value associated with a path in an object.
 *
 * @param {Object} obj Root object.
 * @param {Array} path Array of sequential paths to look up in the object.
 * @returns {*} Item at the requested lookup path; undefined if the path points to no item.
 */
export const objLookup = (obj, path) => {
  if (!path.length) {
    return obj;
  }

  if (obj === undefined || obj === null) {
    return undefined;
  }

  const [current, ...rest] = path;

  return objLookup(obj[current], rest);
};

/**
 * Omit entries from an object with certain keys.
 *
 * @param {Object} obj Input object.
 * @param {Array} paths Array of strings denoting keys in the input object to remove.
 * @return {Object} New object instance with the requested keys removed.
 */
export const omit = (obj, paths = []) => Object.entries(obj)
  .filter(([key]) => !paths.includes(key))
  .reduce((acc, [key, value]) => ({ ...acc, [key]: value }), {});
