import BaseLogic from 'server/logic/base';

/**
 * Logic for reading raw source code.
 */
export default class SourceLogic extends BaseLogic {
  /**
   * Read a source file contents, with transparent caching.
   *
   * @param {String} repo Repository identifier.
   * @param {String} version Git ref or commit-ish at which the file should be fetched.
   * @param {String} file Path to the file to read.
   * @param {Function} cb Callback invoked with (err, resp) on completion.
   */
  read(repo, version, file, cb) {
    this.ctx.log.debug(
      'source',
      'serving file preview',
      { repo, version, file },
    );

    const transaction = this.ctx.cache.transaction(
      'source',
      'read',
      {
        repo,
        revision: version,
        file,
        version: process.env.VERSION,
      },
      JSON.stringify,
      JSON.parse,
    );

    return transaction.get((cacheErr, cached) => {
      if (!cacheErr && cached) {
        return cb(null, cached);
      }

      return this.ctx.source.read(repo, version, file, (err, contents) => {
        if (err) {
          return cb({
            status: 500,
            message: err.toString(),
          });
        }

        const resp = { data: contents };

        transaction.set(resp);

        return cb(null, resp);
      });
    });
  }
}
