import { RECORD_SEARCH_HISTORY_ITEM } from 'client/app/redux/actions/search';
import { recordTelemetryEvent } from 'client/app/redux/actions/telemetry';
import createMiddleware from 'client/app/redux/middleware/create-middleware';
import { TELEMETRY_ACTIONS } from 'shared/constants/telemetry';

/**
 * Emit a telemetry event when a search history item is persisted to client-side storage.
 */
const notifyRecordedSearchHistoryItem = (store, next, action) => {
  store.dispatch(recordTelemetryEvent(TELEMETRY_ACTIONS.RECORD_SEARCH_HISTORY_ITEM));

  return next(action);
};

const middlewareMapping = {
  [RECORD_SEARCH_HISTORY_ITEM]: notifyRecordedSearchHistoryItem,
};

export default createMiddleware(middlewareMapping);
