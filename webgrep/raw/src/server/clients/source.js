import fs from 'fs';
import { Gitlab } from 'gitlab';
import git from 'isomorphic-git';
import { stopwatch } from 'shared/util/instrumentation';

/**
 * Interface for a backend providing access to raw source code.
 */
class SourceBackend {
  /**
   * Read the base64-encoded contents of a file.
   *
   * @param {String} repo Repository identifier.
   * @param {String} version Git ref identifying the commit-ish at which the source should be read.
   * @param {String} path Path to the file.
   * @param {Function} cb Callback function invoked with (err, data) on completion.
   */
  read(repo, version, path, cb) {}  // eslint-disable-line
}

/**
 * Source code viewer backend that errors on all reads.
 */
export class NoopSourceBackend extends SourceBackend {
  read(repo, version, path, cb) {  // eslint-disable-line class-methods-use-this
    return cb(new Error('source code viewer is disabled'));
  }
}

/**
 * Source code viewer backed by a repository accessible on the local disk.
 */
export class LocalSourceBackend extends SourceBackend {
  /**
   * Create a client instance.
   *
   * @param {Array} repos Description of known repositories as an array of objects, each of which
   *                      has property `name` (matching the livegrep repository name) and `path`
   *                      (pointing to the associated git repository path on disk).
   */
  constructor(repos) {
    super();

    // Reshape to a mapping of repository names to filesystem paths
    this.repos = repos
      .reduce((repoPaths, repo) => ({ ...repoPaths, [repo.name]: repo.path }), {});
  }

  read(repo, version, path, cb) {
    if (!(repo in this.repos)) {
      return cb(new Error('no known filesystem path exists for this repository'));
    }

    return git.readBlob({ fs, dir: this.repos[repo], oid: version, filepath: path })
      .then((blob) => cb(null, Buffer.from(blob.blob).toString('base64')))
      .catch((err) => cb(err));
  }
}

/**
 * Source code viewer backed by a Gitlab instance.
 */
export class GitlabSourceBackend extends SourceBackend {
  /**
   * Create a client instance.
   *
   * @param {String} host Base URL to the Gitlab instance.
   * @param {String} accessToken Gitlab API access token.
   */
  constructor(host, accessToken) {
    super();

    this.client = new Gitlab({
      host,
      token: accessToken,
    });
  }

  read(repo, version, path, cb) {
    return this.client.RepositoryFiles.showRaw(
      encodeURIComponent(repo),
      encodeURIComponent(path),
      version,
    )
      .then((raw) => cb(null, Buffer.from(raw).toString('base64')))
      .catch((err) => cb(err));
  }
}

/**
 * Client for viewing raw source code, backed by a source viewer implementation.
 */
export default class SourceClient {
  /**
   * Create a client instance.
   *
   * @param {SourceBackend} backend Instance of a source viewer backend implementation.
   * @param {MetricsClient} metrics Metrics client instance to provide transparent instrumentation
   *                                of source code read calls.
   */
  constructor(backend, metrics) {
    this.backend = backend;
    this.metrics = metrics;
  }

  /**
   * Read the base64-encoded contents of a file.
   *
   * @param {String} repo Repository identifier.
   * @param {String} version Git ref identifying the commit-ish at which the source should be read.
   * @param {String} path Path to the file.
   * @param {Function} cb Callback function invoked with (err, base64 contents) on completion.
   */
  read(repo, version, path, cb) {
    const duration = stopwatch();

    this.metrics.increment('event.source.read.request', { repo });

    return this.backend.read(repo, version, path, (err, contents) => {
      if (err) {
        this.metrics.increment('event.source.read.error', { repo });
        return cb(err);
      }

      this.metrics.timing('latency.source.read', duration.elapsed(), { repo });
      this.metrics.gauge('gauge.source.file_size', contents.length, { repo });

      return cb(null, contents);
    });
  }
}
