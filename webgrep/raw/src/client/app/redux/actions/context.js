export const SET_WINDOW_DIMENSIONS = 'SET_WINDOW_DIMENSIONS';

/**
 * Record the current window dimensions (in pixels).
 *
 * @param {Number} width Current width of the window.
 * @param {Number} height Current height of the window.
 * @returns {Object} Action for setting window dimensions.
 */
export const setWindowDimensions = (width, height) => ({
  type: SET_WINDOW_DIMENSIONS,
  payload: { width, height },
});
