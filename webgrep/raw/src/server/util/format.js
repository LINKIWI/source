/**
 * Escape regex control characters from a string.
 *
 * @param {String} s Input string to escape.
 * @returns {String} Input string with all regex control characters escaped with backslashes.
 */
export const escapeRegex = (s) => s.replace(/[-/\\^$*+?.()|[\]{}]/g, '\\$&');

export default undefined;
