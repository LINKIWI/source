import { RECORD_SEARCH_HISTORY_ITEM } from 'client/app/redux/actions/search';
import { SET_PREFERENCE, CLEAR_PREFERENCE, CLEAR_ALL_PREFERENCES } from 'client/app/redux/actions/preferences';
import createMiddleware from 'client/app/redux/middleware/create-middleware';
import { sharedStorage } from 'client/app/util/storage';

/**
 * Snapshot the current store within the requested namespace to persistent, client-side storage.
 */
const snapshot = (namespace) => (store, next, action) => {
  const reduced = next(action);

  sharedStorage.set('snapshot', namespace, store.getState()[namespace]);

  return reduced;
};

const middlewareMapping = {
  [RECORD_SEARCH_HISTORY_ITEM]: snapshot('search'),
  [SET_PREFERENCE]: snapshot('preferences'),
  [CLEAR_PREFERENCE]: snapshot('preferences'),
  [CLEAR_ALL_PREFERENCES]: snapshot('preferences'),
};

export default createMiddleware(middlewareMapping);
