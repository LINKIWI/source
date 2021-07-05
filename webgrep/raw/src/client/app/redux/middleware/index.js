import thunk from 'redux-thunk';
import errorContext from 'client/app/redux/middleware/error';
import meta from 'client/app/redux/middleware/meta';
import preferences from 'client/app/redux/middleware/preferences';
import search from 'client/app/redux/middleware/search';
import snapshot from 'client/app/redux/middleware/snapshot';
import telemetry from 'client/app/redux/middleware/telemetry';

const middleware = [
  thunk,
  errorContext,
  meta,
  preferences,
  search,
  snapshot,
  telemetry,
];

export default middleware;
