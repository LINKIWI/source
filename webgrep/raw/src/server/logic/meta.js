import BaseLogic from 'server/logic/base';
import { TELEMETRY_ACTIONS } from 'shared/constants/telemetry';

/**
 * Logic for functions related to the running livegrep codesearch server itself.
 */
export default class MetaLogic extends BaseLogic {
  /**
   * Obtain metadata about the index currently loaded into the codesearch backend.
   *
   * @param {Function} cb Callback invoked with (err, data) on completion.
   */
  info(cb) {
    return this.ctx.service.codesearch.rpc('info', {}, (err, data = {}) => {
      if (err) {
        this.ctx.log.error(
          'meta: encountered RPC error: method=info code=%d details=%s',
          err.code,
          err.details,
        );
        return cb(this.formatErr(err));
      }

      const tags = { name: data.name };
      this.ctx.metrics.gauge('gauge.index.repositories', data.trees.length, tags);
      this.ctx.metrics.gauge('gauge.index.timestamp', data.indexTime, tags);

      return cb(null, {
        data: {
          name: data.name,
          timestamp: parseInt(data.indexTime, 10),
          version: process.env.VERSION,
          repositories: data.trees.map((tree) => ({
            name: tree.name,
            version: tree.version,
            url: tree.metadata.urlPattern,
            remote: tree.metadata.remote,
          })),
        },
      });
    });
  }

  /**
   * Report telemetry events as metrics.
   *
   * @param {String} action Telemetry action enum constant.
   * @param {Number} value Value associated with the action.
   * @param {Object} tags Optional object of tags.
   */
  reportTelemetryMetrics(action, value, tags) {
    this.ctx.log.debug(
      'meta: reporting client telemetry event: action=%s value=%s tags=%o',
      action,
      value,
      tags,
    );

    switch (action) {
      // Counters
      case TELEMETRY_ACTIONS.RENDER_VIEW_ROUTE:
      case TELEMETRY_ACTIONS.CLICK_FILE_RESULT_PATH:
      case TELEMETRY_ACTIONS.CLICK_CODE_RESULT_LINE:
      case TELEMETRY_ACTIONS.CLICK_CODE_RESULT_PATH:
      case TELEMETRY_ACTIONS.CLIPBOARD_WRITE:
      case TELEMETRY_ACTIONS.RECORD_SEARCH_HISTORY_ITEM:
      case TELEMETRY_ACTIONS.SOURCE_PREVIEW:
      case TELEMETRY_ACTIONS.SOURCE_RAW_DOWNLOAD:
      case TELEMETRY_ACTIONS.EXECUTE_SEARCH_RECENT:
      case TELEMETRY_ACTIONS.EXECUTE_SEARCH_SOURCE_PREVIEW:
      case TELEMETRY_ACTIONS.EXECUTE_SEARCH_SINGLE_FILE:
      case TELEMETRY_ACTIONS.LOAD_MORE_MATCH_RESULTS:
      case TELEMETRY_ACTIONS.POLL_SERVER_METADATA:
      case TELEMETRY_ACTIONS.COMMIT_SERVER_METADATA:
      case TELEMETRY_ACTIONS.SET_PREFERENCE:
      case TELEMETRY_ACTIONS.RESET_DEFAULT_PREFERENCES:
        return this.ctx.metrics.count('event.telemetry.record', value, { ...tags, action });
      // Gauges (with bucketed aggregation)
      case TELEMETRY_ACTIONS.CLICK_FILE_RESULT_POSITION:
      case TELEMETRY_ACTIONS.CLICK_CODE_RESULT_POSITION:
      case TELEMETRY_ACTIONS.ACTIVE_SESSION_LENGTH:
      case TELEMETRY_ACTIONS.TOTAL_SESSION_LENGTH:
        return this.ctx.metrics.timing('gauge.telemetry.value', value, { ...tags, action });
      default:
        return null;
    }
  }
}
