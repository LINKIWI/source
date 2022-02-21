import { Empty } from 'livegrep/proto/livegrep_pb';
import BaseLogic from 'server/logic/base';

/**
 * Logic for admin functions.
 */
export default class AdminLogic extends BaseLogic {
  /**
   * Force a manual reload of the index loaded into the codesearch backend.
   *
   * @param {Function} cb Callback invoked with (err) on completion.
   */
  reload(cb) {
    // Invalidate all search query cache entries, since they may be stale after reloading the index
    this.ctx.cache.flush('search');

    return this.ctx.service.codesearch.broadcastRPC('reload', new Empty(), (err) =>
      (err ? cb(this.formatErr(err)) : cb()));
  }

  /**
   * Read the raw client and server config data.
   *
   * @param {Function} cb Callback invoked with (err, data) on completion.
   */
  config(cb) {
    return cb(null, {
      data: this.ctx.config.config,
    });
  }
}
