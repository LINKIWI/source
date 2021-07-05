import { route } from 'supercharged/server';
import { HTTPHandler } from 'server/handlers/base';
import { withEndpointInstrumentation } from 'server/util/instrumentation';

@route('/api/admin/reload')
export default class AdminReloadHandler extends HTTPHandler {
  @withEndpointInstrumentation
  put() {
    return this.ctx.logic.admin.reload((err, resp) => (err ? this.error(err) : this.success(resp)));
  }
}
