export const SET_PREFERENCE = 'SET_PREFERENCE';
export const CLEAR_PREFERENCE = 'CLEAR_PREFERENCE';
export const CLEAR_ALL_PREFERENCES = 'CLEAR_ALL_PREFERENCES';

/**
 * Set a preference value.
 *
 * @param {String} key Preference key.
 * @param {String} value Preference value.
 * @return {Object} Action object for setting a preference.
 */
export const setPreference = (key, value) => ({
  type: SET_PREFERENCE,
  payload: { key, value },
});

/**
 * Clear a preference by key.
 *
 * @param {String} key Preference key.
 * @return {Object} Action object for clearing a single preference.
 */
export const clearPreference = (key) => ({
  type: CLEAR_PREFERENCE,
  payload: { key },
});

/**
 * Reset all preferences to default.
 *
 * @return {Object} Action object for clearing all preferences.
 */
export const clearAllPreferences = () => ({
  type: CLEAR_ALL_PREFERENCES,
});
