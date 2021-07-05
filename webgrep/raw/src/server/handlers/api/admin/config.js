import { route } from 'supercharged/server';
import { HTTPHandler } from 'server/handlers/base';
import { withEndpointInstrumentation } from 'server/util/instrumentation';

@route('/api/admin/config')
export default class AdminConfigHandler extends HTTPHandler {
  @withEndpointInstrumentation
  get() {
    return this.ctx.logic.admin.config((err, resp) => (err ? this.error(err) : this.success(resp)));
  }
}
