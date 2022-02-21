import { route, withRequestSchema } from 'supercharged/server';
import { WebSocketTransactionHandler } from 'server/handlers/base';
import { withEndpointInstrumentation } from 'server/util/instrumentation';
import { TELEMETRY_ACTIONS, TELEMETRY_TAGS } from 'shared/constants/telemetry';

@route('/api/meta/telemetry')
export default class MetaTelemetryHandler extends WebSocketTransactionHandler {
  @withEndpointInstrumentation
  @withRequestSchema({
    type: 'object',
    properties: {
      action: {
        enum: Object.values(TELEMETRY_ACTIONS),
      },
      value: {
        type: 'number',
      },
      tags: {
        type: 'object',
        patternProperties: {
          '.+': {
            type: ['string', 'number'],
          },
        },
        additionalProperties: false,
      },
      timestamp: {
        type: 'number',
      },
      attempt: {
        type: 'number',
      },
    },
    required: ['action', 'value', 'tags'],
    additionalProperties: false,
  })
  invoke({ action, value, tags }) {
    // Only include whitelisted tag keys for the particular action type
    const sanitizedTags = (TELEMETRY_TAGS[action] || [])
      .filter((knownTagKey) => knownTagKey in tags)
      .reduce((acc, knownTagKey) => ({
        ...acc,
        [knownTagKey]: tags[knownTagKey],
      }), {});

    this.ctx.logic.meta.reportTelemetryMetrics(action, value, sanitizedTags);
  }
}
