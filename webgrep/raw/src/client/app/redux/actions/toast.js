import UIDGenerator from 'client/app/util/uid-generator';

export const SHOW_TOAST = 'SHOW_TOAST';
export const HIDE_TOAST = 'HIDE_TOAST';

// Default toast display time, in milliseconds.
const DEFAULT_TOAST_DURATION = 5000;
const TOAST_ID_GENERATOR = new UIDGenerator();

/**
 * Show a toast.
 *
 * @param {Number} toastID Unique toast ID.
 * @param {String} text Toast text.
 * @param {Number} duration Duration, in milliseconds, that the toast will be active.
 * @return {Object} Action object for displaying a toast.
 */
export const showToast = (toastID, text, duration) => ({
  type: SHOW_TOAST,
  payload: { toastID, text, duration },
});

/**
 * Hide a toast.
 *
 * @param {Number} toastID Unique toast ID of a currently displayed toast.
 * @return {Object} Action object for hiding a toast.
 */
export const hideToast = (toastID) => ({
  type: HIDE_TOAST,
  payload: { toastID },
});

/**
 * Convenience action to asynchronously dispatch a show and hide action with a timeout in between.
 *
 * @param {String} text Toast text.
 * @param {Number} duration Duration, in milliseconds, that the toast will be active.
 * @return {Object} Action object for displaying a toast.
 */
export const cycleToast = (text, duration = DEFAULT_TOAST_DURATION) => (dispatch) => {
  const allocatedToastID = TOAST_ID_GENERATOR.next();

  dispatch(showToast(allocatedToastID, text, duration));
  setTimeout(() => dispatch(hideToast(allocatedToastID)), duration);
};
