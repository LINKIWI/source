import fs from 'fs';
import path from 'path';
import BaseLogic from 'server/logic/base';

const CLIENT_BUNDLE_PATH = path.join(__dirname, '../../../dist/client/js/main.js');

const clientTemplate = (injectedGlobals, bundle) => `
  <!DOCTYPE html>
  <html lang="en">
    <head>
      <meta charset="UTF-8" />
      <meta
        name="viewport"
        content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=0"
      />
    </head>
    <body>
      <div id="app"></div>
      <script>window.__SSR_INJECTED_GLOBALS__ = ${JSON.stringify(injectedGlobals)}</script>
      <script>${bundle}</script>
    </body>
  </html>
`;

const renderTemplate = (injectedGlobals, cb) => fs.readFile(CLIENT_BUNDLE_PATH, (err, bundle) =>
  (err ? cb(err) : cb(null, clientTemplate(injectedGlobals, bundle.toString()))));

/**
 * Logic related to rendering client templates.
 */
export default class ViewLogic extends BaseLogic {
  /**
   * Render the universal shared frontend.
   *
   * @param {Function} cb Callback invoked with (err, html) on completion. The function may error if
   *                      the client bundle does not exist on disk.
   */
  renderFrontend(cb) {
    this.ctx.metrics.increment('event.frontend.render');

    return this.ctx.logic.meta.info((err, { data = {} } = {}) => {
      const meta = {
        name: data.name || '',
        timestamp: data.timestamp || 0,
        version: data.version || '',
        // Key repositories by name for fast client-side lookup
        repositories: (data.repositories || []).reduce((acc, repo) => ({
          ...acc,
          [repo.name]: repo,
        }), {}),
      };

      const hydratedStore = {
        config: this.ctx.config.config,
        meta,
      };

      return renderTemplate({ store: hydratedStore }, cb);
    });
  }
}
