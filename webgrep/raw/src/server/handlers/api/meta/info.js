import { route, withRequestSchema } from 'supercharged/server';
import { HTTPHandler, WebSocketTransactionHandler } from 'server/handlers/base';
import { withEndpointInstrumentation } from 'server/util/instrumentation';

@route('/api/meta/info')
export class MetaInfoHandler extends HTTPHandler {
  @withEndpointInstrumentation
  get() {
    return this.ctx.logic.meta.info((err, resp) => (err ? this.error(err) : this.success(resp)));
  }
}

@route('/api/meta/info')
export class MetaInfoLiveHandler extends WebSocketTransactionHandler {
  @withEndpointInstrumentation
  @withRequestSchema({
    type: 'object',
    properties: {
      // Optional message transaction ID used by the websocket duplex stream for strict
      // request/response ordering
      id: {
        type: 'number',
      },
    },
    additionalProperties: false,
  })
  invoke({ id }) {
    return this.ctx.logic.meta.info((err, resp) =>
      (err ? this.error({ id, ...err }) : this.success({ id, ...resp })));
  }
}
