import {
  HTTPHandler as SuperchargedHTTPHandler,
  WebSocketHandler as SuperchargedWebSocketHandler,
} from 'supercharged/server';
import { CODE_SERVER_UNDEFINED } from 'supercharged/shared/constants/error';
import { parseJSON } from 'shared/util/data';

export class HTTPHandler extends SuperchargedHTTPHandler {
  constructor(ctx, ...args) {
    super(...args);

    this.ctx = ctx;
  }
}

export class WebSocketHandler extends SuperchargedWebSocketHandler {
  constructor(ctx, ...args) {
    super(...args);

    this.ctx = ctx;
  }
}

/**
 * Wrapper over a WebSocketHandler that provides abstractions for using a websocket channel to
 * serially send and receive schema-adherent data in the style of an HTTP request/response
 * transaction. The response shape is similar to that of a canonical Supercharged HTTP response.
 */
export class WebSocketTransactionHandler extends WebSocketHandler {
  // Subclasses should override this method, optionally with a withSchema decorator.
  invoke() {
    return this.success({});
  }

  onMessage(message) {
    this.req.body = parseJSON(message).data || {};

    // Parameters are automatically injected by the schema validation decorator
    return this.invoke();
  }

  error({ id = this.req.body.id, code = CODE_SERVER_UNDEFINED, message = null, data = null }) {
    return this._sendJSON({
      success: false,
      id,
      code,
      message,
      data,
    });
  }

  success({ id = this.req.body.id, data = null }) {
    return this._sendJSON({
      success: true,
      id,
      data,
    });
  }

  _sendJSON(payload) {
    return this.sendMessage(JSON.stringify(payload));
  }
}
