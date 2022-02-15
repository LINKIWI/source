import { RECORD_TELEMETRY_EVENT } from 'client/app/redux/actions/telemetry';
import createReducer from 'client/app/redux/reducers/create-reducer';

// Retain only this many recent telemetry events in history.
const MAX_TELEMETRY_HISTORY_EVENTS = 20;

const initialState = {
  events: [],
};

const recordTelemetryEventReducer = (state, action) => {
  // Skip appending telemetry events that are retried due to initial failure.
  if (action.payload.attempt) {
    return state;
  }

  return {
    ...state,
    events: [action.payload, ...state.events].slice(0, MAX_TELEMETRY_HISTORY_EVENTS),
  };
};

const reducerMapping = {
  [RECORD_TELEMETRY_EVENT]: recordTelemetryEventReducer,
};

export default createReducer(reducerMapping, initialState);
