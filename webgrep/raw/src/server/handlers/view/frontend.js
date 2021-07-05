import { route } from 'supercharged/server';
import { HTTPHandler } from 'server/handlers/base';
import { withEndpointInstrumentation } from 'server/util/instrumentation';

@route(['/', '/admin', '/admin/*', '/search'])
export default class FrontendHandler extends HTTPHandler {
  @withEndpointInstrumentation
  get() {
    return this.ctx.logic.view.renderFrontend((err, html) => {
      if (err) {
        throw err;
      }

      this.res.setHeader('Content-Type', 'text/html');
      this.res.send(html);
    });
  }
}
