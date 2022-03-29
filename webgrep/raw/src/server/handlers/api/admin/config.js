import { route, withRequestLog } from 'supercharged/server';
import { HTTPHandler } from 'server/handlers/base';
import { withEndpointInstrumentation } from 'server/util/instrumentation';

@route('/api/admin/config')
export default class AdminConfigHandler extends HTTPHandler {
  @withRequestLog
  @withEndpointInstrumentation
  get() {
    return this.ctx.logic.admin.config((err, resp) => (err ? this.error(err) : this.success(resp)));
  }
}
