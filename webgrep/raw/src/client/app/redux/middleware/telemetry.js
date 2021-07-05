import ReconnectingWebSocket from 'reconnectingwebsocket';
import { RECORD_TELEMETRY_EVENT } from 'client/app/redux/actions/telemetry';
import createMiddleware from 'client/app/redux/middleware/create-middleware';

// Constant number of milliseconds to wait before attempting to retry reporting a telemetry event.
const FAILED_REPORT_RETRY_DELAY = 2000;

// Mechanics for instantiating a persistent websocket connection to the telemetry endpoint
const wsProtocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
const endpoint = '/api/meta/telemetry';
const connectURL = `${wsProtocol}://${window.location.host}${endpoint}`;
const socket = new ReconnectingWebSocket(connectURL, null, { reconnectDecay: 1 });

/**
 * Report a telemetry event over the websocket. If the websocket connection is unavailable, this
 * middleware will schedule a retry by dispatching the same action after a timeout.
 */
const reportTelemetryEvent = (store, next, action) => {
  if (socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(action.payload));
  } else {
    window.setTimeout(() => store.dispatch(action), FAILED_REPORT_RETRY_DELAY);
  }

  return next(action);
};

const middlewareMapping = {
  [RECORD_TELEMETRY_EVENT]: reportTelemetryEvent,
};

export default createMiddleware(middlewareMapping);
