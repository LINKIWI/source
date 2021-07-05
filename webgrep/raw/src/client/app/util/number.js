/**
 * Compute the modulus of n by m, supporting negative values of the dividend.
 *
 * @param {Number} n Dividend.
 * @param {Number} m Divisor.
 * @returns {Number} Result of n % m.
 */
export const modulo = (n, m) => ((n % m) + m) % m;

export default undefined;
