export const SET_INDEX_META = 'SET_INDEX_META';

/**
 * Update the store's server-reported index metadata.
 *
 * @param {Object} info Server index metadata object.
 * @return {Object} Action object for setting the index metadata.
 */
export const setIndexMeta = (info) => ({
  type: SET_INDEX_META,
  payload: { info },
});
