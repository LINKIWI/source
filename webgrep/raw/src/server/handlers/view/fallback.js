import { route } from 'supercharged/server';
import FrontendHandler from 'server/handlers/view/frontend';
import { withEndpointInstrumentation } from 'server/util/instrumentation';

@route('*')
export default class FallbackHandler extends FrontendHandler {
  @withEndpointInstrumentation
  get() {
    this.res.status(404);
    return super.get();
  }
}
