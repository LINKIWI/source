import { combineReducers } from 'redux';
import configReducer from 'client/app/redux/reducers/config';
import contextReducer from 'client/app/redux/reducers/context';
import metaReducer from 'client/app/redux/reducers/meta';
import preferencesReducer from 'client/app/redux/reducers/preferences';
import searchReducer from 'client/app/redux/reducers/search';
import telemetryReducer from 'client/app/redux/reducers/telemetry';
import toastReducer from 'client/app/redux/reducers/toast';

const reducer = combineReducers({
  config: configReducer,
  context: contextReducer,
  meta: metaReducer,
  preferences: preferencesReducer,
  search: searchReducer,
  telemetry: telemetryReducer,
  toasts: toastReducer,
});

export default reducer;
