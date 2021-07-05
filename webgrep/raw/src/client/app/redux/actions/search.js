export const RECORD_SEARCH_HISTORY_ITEM = 'RECORD_SEARCH_HISTORY_ITEM';

/**
 * Record a search history item in the store.
 *
 * @param {String} query Search query, as a string.
 * @return {Object} Action object for recording a history item.
 */
export const recordSearchHistoryItem = (query) => ({
  type: RECORD_SEARCH_HISTORY_ITEM,
  payload: { query },
});
