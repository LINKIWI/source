/**
 * Do a one-shot render against a fully materialized blob of text.
 *
 * Template variables in the source blob are expected to be of the form {NAME}, i.e. surrounded by
 * curly braces. All such instances are replaced with the value associated with the NAME key in the
 * template variables.
 *
 * @param {String} source Original template text from which output should be rendered.
 * @param {Object} variables Mapping of template variable names to values to render.
 * @returns {String} Rendered output.
 */
export const render = (source, variables) => Object.entries(variables)
  // The String.prototype.replaceAll function recognizes special symbols in the replacement value
  // to change the replacement behavior. To force replacement of the string literal, use the
  // alternate API which invokes a callback function to supply the literal replacement value.
  .reduce((rendered, [key, value]) => rendered.replaceAll(`{${key}}`, () => value), source);

export default undefined;
