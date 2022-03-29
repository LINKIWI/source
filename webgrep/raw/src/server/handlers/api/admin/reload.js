import { route, withRequestLog } from 'supercharged/server';
import { HTTPHandler } from 'server/handlers/base';
import { withEndpointInstrumentation } from 'server/util/instrumentation';

@route('/api/admin/reload')
export default class AdminReloadHandler extends HTTPHandler {
  @withRequestLog
  @withEndpointInstrumentation
  put() {
    return this.ctx.logic.admin.reload((err, resp) => (err ? this.error(err) : this.success(resp)));
  }
}
