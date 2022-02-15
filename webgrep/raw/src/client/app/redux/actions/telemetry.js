export const RECORD_TELEMETRY_EVENT = 'RECORD_TELEMETRY_EVENT';

/**
 * Record a telemetry event.
 *
 * @param {String} action Telemetry action enum constant.
 * @param {Number} value Value associated with the action.
 * @param {Object} tags Optional object of tags.
 * @return {Object} Action object for recording a telemetry event.
 */
export const recordTelemetryEvent = (action, value = 1, tags = {}) => ({
  type: RECORD_TELEMETRY_EVENT,
  payload: { action, value, tags, timestamp: Date.now() },
});

/**
 * Retry the reporting of a previously created telemetry event.
 *
 * @param {Object} event Previously created action of type RECORD_TELEMETRY_EVENT.
 * @returns {Object} Modified action with a flag indicating the report retry attempt.
 */
export const retryTelemetryEvent = (event) => ({
  type: RECORD_TELEMETRY_EVENT,
  payload: { ...event.payload, attempt: (event.payload.attempt || 0) + 1 },
});
