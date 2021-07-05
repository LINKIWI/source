import humanize from 'humanize';

/**
 * Format a string with a map of identifier values. It is assumed that the format string specifies
 * identifiers enclosed with curly braces, e.g.:
 *
 *   Base: 'Hello, {entity}'
 *   Formatters: { entity: 'world!' }
 *   Output: 'Hello, world!'
 *
 * @param {String} base Base string to format.
 * @param {Object} formatters Map of identifier names to values.
 * @returns {String} Formatted string.
 */
export const string = (base = '', formatters = {}) => Object.entries(formatters)
  .reduce((formatted, [specifier, content]) => formatted.replace(`{${specifier}}`, content), base);

/**
 * Format an absolute Unix timestamp in the local timezone.
 *
 * @param {Number} timestamp Unix timestamp, in seconds.
 * @returns {String} Human-consumable absolute description of the timestamp.
 */
export const absoluteTimestamp = (timestamp) => humanize.date('F j, Y, g:i A', timestamp);

/**
 * Format a relative Unix timestamp.
 *
 * @param {Number} timestamp Unix timestamp, in seconds.
 * @returns {String} Human-consumable relative description of the timestamp.
 */
export const relativeTimestamp = (timestamp) => humanize.relativeTime(timestamp);
