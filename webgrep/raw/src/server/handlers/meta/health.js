import { route } from 'supercharged/server';
import { HTTPHandler } from 'server/handlers/base';
import { withEndpointInstrumentation } from 'server/util/instrumentation';

@route('/_/health')
export default class HealthHandler extends HTTPHandler {
  @withEndpointInstrumentation
  get() {
    return this.success({ message: 'OK' });
  }
}
