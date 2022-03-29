import { route, withRequestLog, withRequestSchema } from 'supercharged/server';
import { HTTPHandler, WebSocketTransactionHandler } from 'server/handlers/base';
import { withEndpointInstrumentation } from 'server/util/instrumentation';

const schema = {
  type: 'object',
  properties: {
    // Search parameters
    query: {
      type: 'string',
      minLength: 1,
    },
    file: {
      type: 'string',
      default: '',
    },
    repos: {
      type: 'array',
      items: {
        type: 'string',
        minLength: 1,
      },
      default: [],
    },
    regex: {
      type: 'boolean',
      default: false,
    },
    caseSensitive: {
      type: 'boolean',
      default: false,
    },
    maxMatches: {
      type: 'number',
      minimum: -1,
      default: 50,
    },
    context: {
      type: 'number',
      minimum: 0,
      default: 4,
    },
    // Optional message transaction ID used by the websocket duplex stream for strict
    // request/response ordering
    id: {
      type: 'number',
    },
    // Optional unique index identity descriptor used for slightly more robust consistency between
    // server and client index state
    indexIdentity: {
      type: 'object',
      properties: {
        name: {
          type: 'string',
        },
        timestamp: {
          type: 'number',
        },
      },
      required: ['name', 'timestamp'],
      additionalProperties: false,
    },
  },
  required: ['query'],
  additionalProperties: false,
};

@route('/api/search')
export class SearchHandler extends HTTPHandler {
  @withRequestLog
  @withEndpointInstrumentation
  @withRequestSchema(schema)
  get(params) {
    return this.ctx.logic.search.executeSearch(params, (err, resp) =>
      (err ? this.error(err) : this.success(resp)));
  }
}

@route('/api/search')
export class SearchLiveHandler extends WebSocketTransactionHandler {
  @withRequestLog
  @withEndpointInstrumentation
  @withRequestSchema(schema)
  invoke({ id, ...params }) {
    return this.ctx.logic.search.executeSearch(params, (err, resp) =>
      (err ? this.error({ id, ...err }) : this.success({ id, ...resp })));
  }
}
