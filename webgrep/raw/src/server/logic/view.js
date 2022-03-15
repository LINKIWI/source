import { CODE_NOT_FOUND, CODE_SERVER_UNDEFINED } from 'supercharged/shared/constants/error';
import BaseLogic from 'server/logic/base';
import { clientResource, serverResource } from 'server/util/resource';
import { render } from 'server/util/template';

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
      const context = {
        timestamp: Date.now(),
        window: { width: null, height: null },
      };

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
        context,
        meta,
      };

      return serverResource('resources/template/frontend.html', (frontendTemplateErr, index) => {
        if (frontendTemplateErr) {
          return cb(frontendTemplateErr);
        }

        return clientResource('js/main.js', (clientBundleErr, bundle) => {
          if (clientBundleErr) {
            return cb(clientBundleErr);
          }

          const variables = {
            SEARCH_ENGINE_NAME: this.ctx.config.get('client.search.name') || 'webgrep',
            CLIENT_BUNDLE: bundle.toString(),
            SSR_INJECTED_GLOBALS: JSON.stringify({ store: hydratedStore }),
          };

          return cb(null, render(index.toString(), variables));
        });
      });
    });
  }

  /**
   * Render a resource.
   *
   * @param {String} path Original request path.
   * @param {Function} cb Function invoked with (err, resp) on completion. The response is an object
   *                      with keys `header` and `body` describing HTTP response parameters.
   */
  renderResource(path, cb) {
    switch (path) {
      case '/resource/opensearch.xml':
        return serverResource('resources/template/opensearch.xml', (specErr, spec) => {
          if (specErr) {
            return cb({
              status: 500,
              code: CODE_SERVER_UNDEFINED,
              message: specErr.toString(),
            });
          }

          return clientResource('img/logo.png.base64', (logoErr, logo) => {
            if (logoErr) {
              return cb({
                status: 500,
                code: CODE_SERVER_UNDEFINED,
                message: logoErr.toString(),
              });
            }

            const variables = {
              NAME: this.ctx.config.get('client.search.name') || 'webgrep',
              DESCRIPTION: this.ctx.config.get('client.search.description') ||
                'webgrep code search',
              LOGO: logo.toString(),
              BASE_URL: this.ctx.config.get('client.search.base_url') || '',
            };

            return cb(null, {
              headers: { 'Content-Type': 'application/opensearchdescription+xml' },
              body: render(spec.toString(), variables),
            });
          });
        });

      default:
        return cb({
          status: 404,
          code: CODE_NOT_FOUND,
          message: 'no such resource',
          data: { path },
        });
    }
  }
}
