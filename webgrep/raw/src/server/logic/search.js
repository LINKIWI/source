import { Query, SearchStats } from 'livegrep/proto/livegrep_pb';
import BaseLogic from 'server/logic/base';
import { invertMap } from 'server/util/data';
import { escapeRegex } from 'server/util/format';
import { ERROR_CODES } from 'shared/constants/error';

/**
 * Logic for executing index searches.
 */
export default class SearchLogic extends BaseLogic {
  /**
   * Execute a search against the livegrep codesearch backend.
   *
   * @param {Object} params Object of search parameters.
   * @param {Function} cb Callback invoked with (err, resp) on completion.
   */
  executeSearch(params, cb) {
    const req = {
      line: params.regex ? params.query : escapeRegex(params.query),
      file: params.regex ? params.file : escapeRegex(params.file),
      // Manually construct a regular expression for a logical disjunction among all repositories
      repo: params.repos.map((repo) => `^${escapeRegex(repo)}$`).join('|'),
      foldCase: !params.caseSensitive,
      maxMatches: params.maxMatches,
      contextLines: params.context,
    };

    this.ctx.log.debug(
      'search',
      'serving query',
      {
        line: req.line,
        file: params.line || '(none)',
        repos: params.repos.length ? params.repos : '(all)',
        case: params.caseSensitive,
        matches: params.maxMatches,
        context: params.context,
      },
    );

    const transaction = this.ctx.cache.transaction(
      'search',
      'query-results',
      {
        ...req,
        // Encode the current server-side application version into the cache key for robustness
        // against potential key schema changes across different versions
        version: process.env.VERSION,
        // Encode the unique index identity into the cache key for robustness against stale values
        // (within TTL) after the server index is rolled
        ...params.indexIdentity && {
          indexName: params.indexIdentity.name,
          indexTimestamp: params.indexIdentity.timestamp,
        },
      },
      JSON.stringify,
      JSON.parse,
    );

    return transaction.get((cacheErr, cached) => {
      if (!cacheErr && cached) {
        return cb(null, cached);
      }

      const query = new Query();
      query.setLine(req.line);
      query.setFile(req.file);
      query.setRepo(req.repo);
      query.setFoldCase(req.foldCase);
      query.setMaxMatches(req.maxMatches);
      query.setContextLines(req.contextLines);

      return this.ctx.service.codesearch.rpc('search', query, (rpcErr, rpcResp) => {
        if (rpcErr) {
          this.ctx.log.error(
            'search',
            'encountered RPC error',
            { method: 'search', code: rpcErr.code, details: rpcErr.details, request: req },
          );
          return cb(this.formatErr(rpcErr));
        }

        const data = rpcResp.toObject();

        // Treat early terminations due to server-side timeout as fatal errors
        if (data.stats.exitReason === SearchStats.ExitReason.TIMEOUT) {
          this.ctx.log.error(
            'search: exceeded livegrep search timeout: request=%s',
            JSON.stringify(req),
          );
          return cb({
            code: ERROR_CODES.SEARCH_TIMEOUT,
            message: 'Server-side search timeout exceeded',
          });
        }

        const { code, stats, files } = this._reshapeResults(data);

        this.ctx.metrics.gauge('gauge.search.code_results', data.resultsList.length);
        this.ctx.metrics.gauge('gauge.search.path_results', data.fileResultsList.length);
        this.ctx.metrics.gauge('gauge.search.file_results', code.length);
        this.ctx.metrics.timing(
          'latency.search.exec',
          stats.time,
          { exit: invertMap(SearchStats.ExitReason)[stats.exitReason] },
        );

        const resp = {
          data: {
            stats,
            code,
            files,
          },
        };

        // Only write to the cache if the search result was produced by the same index identity as
        // that specified in the request. Otherwise, this would pollute the cache with data from the
        // wrong index for this request.
        // Note that a side effect of this logic is that requests without a supplied index identity
        // will never write to cache, in order to avoid potential index inconsistency.
        const isCacheWriteSafe = params.indexIdentity &&
          data.indexName === params.indexIdentity.name &&
          parseInt(data.indexTime, 10) === params.indexIdentity.timestamp;

        if (isCacheWriteSafe) {
          transaction.set(resp);
        }

        return cb(null, resp);
      });
    });
  }

  /**
   * Massage the response from livegrep into a form that can be more easily interpreted by the
   * webgrep frontend.
   *
   * @param {Object} data Raw response payload from livegrep.
   * @return {Object} Object of code and file results, and search stats.
   * @private
   */
  _reshapeResults(data) {  // eslint-disable-line class-methods-use-this
    const code = Object.values(data.resultsList
      // Aggregate lines by repo and path, so that each unique (repo, path) combination is
      // described by an array of all matching lines and the left/right bounds for each line.
      .reduce((aggregated, result) => {
        // Aggregation key: combine all results for the same file in the same repo
        const key = `${result.tree}-${result.path}`;
        // Line number of the matching line
        const lineNumber = parseInt(result.lineNumber, 10);
        // The existing entry for this repo/path combination, if it exists
        const existing = aggregated[key] || {};

        // Create a map of line numbers -> { bounds description, line }, being careful not to
        // override the bounds if they have already been specified. Since context lines are
        // overlapping, it's possible that a context line does not have a bounds description,
        // but it has one from an earlier result.
        const contextLines = [
          ...result.contextBeforeList.reverse(),
          result.line,
          ...result.contextAfterList,
        ].reduce((lines, line, idx) => {
          const contextLno = idx + lineNumber - result.contextBeforeList.length;
          const bounds = (() => {
            // Examining the matching line, for which bounds information is available
            if (contextLno === lineNumber) {
              return [result.bounds.left, result.bounds.right];
            }

            // Defer to existing bounds information
            return existing.lines &&
              existing.lines[contextLno] &&
              existing.lines[contextLno].bounds;
          })();

          return {
            ...lines,
            [contextLno]: { bounds, line },
          };
        }, {});

        return {
          ...aggregated,
          [key]: {
            repo: result.tree,
            version: result.version,
            path: result.path,
            lines: {
              ...existing.lines || {},
              ...contextLines,
            },
          },
        };
      }, {}))
      .map((resultGroup) => ({
        ...resultGroup,
        lines: Object.entries(resultGroup.lines)
          .map(([number, details]) => ({
            ...details,
            number: parseInt(number, 10),
          }))
          .sort((a, b) => a.number - b.number),
      }));

    const files = Object.values(data.fileResultsList.reduce((acc, file) => ({
      ...acc,
      // Deduplicate results keyed by its repository and file path
      [`${file.tree}-${file.path}`]: {
        repo: file.tree,
        version: file.version,
        path: file.path,
        bounds: [file.bounds.left, file.bounds.right],
      },
    }), {}));

    const stats = {
      exitReason: data.stats.exitReason,
      time: parseInt(data.stats.totalTime, 10),
    };

    return { code, files, stats };
  }
}
