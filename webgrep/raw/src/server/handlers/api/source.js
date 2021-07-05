import { route, withSchema } from 'supercharged/server';
import { HTTPHandler } from 'server/handlers/base';
import { withEndpointInstrumentation } from 'server/util/instrumentation';

@route('/api/source')
export default class SourceHandler extends HTTPHandler {
  @withEndpointInstrumentation
  @withSchema({
    properties: {
      repo: {
        type: 'string',
        minLength: 1,
      },
      version: {
        type: 'string',
        minLength: 1,
      },
      path: {
        type: 'string',
        minLength: 1,
      },
    },
    required: ['repo', 'version', 'path'],
    additionalProperties: false,
  })
  get({ repo, version, path }) {
    this.ctx.logic.source.read(repo, version, path, (err, resp) =>
      (err ? this.error(err) : this.success(resp)));
  }
}
