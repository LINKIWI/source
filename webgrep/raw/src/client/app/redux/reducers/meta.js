import { SET_INDEX_META } from 'client/app/redux/actions/meta';
import createReducer from 'client/app/redux/reducers/create-reducer';

const initialState = {
  // Injected by SSR store hydration
};

const setIndexMetaReducer = (state, action) => {
  const { info } = action.payload;

  return {
    ...state,
    name: info.name,
    timestamp: info.timestamp,
    version: info.version,
    // Key repositories by name to be consistent with SSR-hydrated schema
    repositories: (info.repositories || []).reduce((acc, repo) => ({
      ...acc,
      [repo.name]: repo,
    }), {}),
  };
};

const reducerMapping = {
  [SET_INDEX_META]: setIndexMetaReducer,
};

export default createReducer(reducerMapping, initialState);
