import { route } from 'supercharged/server';
import { HTTPHandler } from 'server/handlers/base';
import { withEndpointInstrumentation } from 'server/util/instrumentation';

@route(['/resource/*'])
export default class ResourcesHandler extends HTTPHandler {
  @withEndpointInstrumentation
  get() {
    return this.ctx.logic.view.renderResource(this.req.path, (err, resp) => {
      if (err) {
        return this.error(err);
      }

      Object.entries(resp.headers || {})
        .forEach(([key, value]) => this.res.setHeader(key, value));

      return this.res.send(resp.body);
    });
  }
}
