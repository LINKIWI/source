import { SET_INDEX_META } from 'client/app/redux/actions/meta';
import { cycleToast } from 'client/app/redux/actions/toast';
import createMiddleware from 'client/app/redux/middleware/create-middleware';
import { relativeTimestamp } from 'client/app/util/format';

/**
 * Show a toast notification if the updated index metadata indicates that the remote search index
 * has changed.
 */
const notifyIndexUpdate = (store, next, action) => {
  const { name, timestamp } = action.payload.info;
  const { name: lastName, timestamp: lastTimestamp } = store.getState().meta;

  // Toast a notification if the server index identity has changed
  if (lastName !== name || lastTimestamp !== timestamp) {
    const notification = [
      'Remote search index updated!',
      `Now searching ${name}, built ${relativeTimestamp(timestamp)}.`,
    ].join(' ');
    store.dispatch(cycleToast(notification), 10 * 1000);
  }

  return next(action);
};

const middlewareMapping = {
  [SET_INDEX_META]: notifyIndexUpdate,
};

export default createMiddleware(middlewareMapping);
