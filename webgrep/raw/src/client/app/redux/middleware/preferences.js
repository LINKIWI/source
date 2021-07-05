import { SET_PREFERENCE, CLEAR_ALL_PREFERENCES } from 'client/app/redux/actions/preferences';
import { recordTelemetryEvent } from 'client/app/redux/actions/telemetry';
import { cycleToast } from 'client/app/redux/actions/toast';
import createMiddleware from 'client/app/redux/middleware/create-middleware';
import { TELEMETRY_ACTIONS } from 'shared/constants/telemetry';

/**
 * Show a toast notification and emit a telemetry event when the preference value has changed.
 */
const notifyPreferenceSave = (store, next, action) => {
  const { key, value } = action.payload;

  if (store.getState().preferences[key] !== value) {
    store.dispatch(cycleToast('Preferences saved successfully.'));
    store.dispatch(recordTelemetryEvent(TELEMETRY_ACTIONS.SET_PREFERENCE, 1, { key }));
  }

  return next(action);
};

/**
 * Show a toast notification and emit a telemetry event when all preferences are reset to defaults.
 */
const notifyPreferencesReset = (store, next, action) => {
  store.dispatch(cycleToast('All preferences reset to defaults successfully.'));
  store.dispatch(recordTelemetryEvent(TELEMETRY_ACTIONS.RESET_DEFAULT_PREFERENCES));

  return next(action);
};

const middlewareMapping = {
  [SET_PREFERENCE]: notifyPreferenceSave,
  [CLEAR_ALL_PREFERENCES]: notifyPreferencesReset,
};

export default createMiddleware(middlewareMapping);
