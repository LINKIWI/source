import { SET_PREFERENCE, CLEAR_PREFERENCE, CLEAR_ALL_PREFERENCES } from 'client/app/redux/actions/preferences';
import createReducer from 'client/app/redux/reducers/create-reducer';
import { PREFERENCE_KEYS, PREFERENCE_VALUES } from 'client/app/util/constants/preferences';
import { sharedStorage } from 'client/app/util/storage';

const defaultState = {
  [PREFERENCE_KEYS.CODE_SEARCH_CONTEXT_LINES]: 4,
  [PREFERENCE_KEYS.CODE_SEARCH_INITIAL_MATCH_LIMIT]: 50,
  [PREFERENCE_KEYS.CODE_SEARCH_FILE_RESULTS_LIMIT]: 10,
  [PREFERENCE_KEYS.CODE_SEARCH_RESULT_NAVIGATION]:
    PREFERENCE_VALUES[PREFERENCE_KEYS.CODE_SEARCH_RESULT_NAVIGATION].SOURCE_LINK_NEW_TAB,
  [PREFERENCE_KEYS.CODE_SEARCH_TRANSPORT_PROTOCOL]:
    PREFERENCE_VALUES[PREFERENCE_KEYS.CODE_SEARCH_TRANSPORT_PROTOCOL].AUTO,
  [PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME]:
    PREFERENCE_VALUES[PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME].COY,
};

const initialState = {
  ...defaultState,
  ...sharedStorage.get('snapshot', 'preferences'),
};

const setPreferenceReducer = (state, action) => {
  const { key, value } = action.payload;

  return {
    ...state,
    [key]: value,
  };
};

const clearPreferenceReducer = (state, action) => {
  const { key } = action.payload;

  return {
    ...state,
    [key]: defaultState[key],
  };
};

const clearAllPreferencesReducer = () => defaultState;

const reducerMapping = {
  [SET_PREFERENCE]: setPreferenceReducer,
  [CLEAR_PREFERENCE]: clearPreferenceReducer,
  [CLEAR_ALL_PREFERENCES]: clearAllPreferencesReducer,
};

export default createReducer(reducerMapping, initialState);
