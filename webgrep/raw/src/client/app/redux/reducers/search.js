import { RECORD_SEARCH_HISTORY_ITEM } from 'client/app/redux/actions/search';
import createReducer from 'client/app/redux/reducers/create-reducer';
import { sharedStorage } from 'client/app/util/storage';

// Retain only this many recent search items in history.
const MAX_SEARCH_HISTORY_ITEMS = 5;

const initialState = {
  history: [],
  ...sharedStorage.get('snapshot', 'search'),
};

const recordSearchHistoryItemReducer = (state, action) => {
  const { query } = action.payload;

  return {
    ...state,
    history: [
      query,
      // Don't allow duplicating history items; instead, bump the existing one to the top.
      ...state.history.filter((saved) => saved !== query),
    ].slice(0, MAX_SEARCH_HISTORY_ITEMS),
  };
};

const reducerMapping = {
  [RECORD_SEARCH_HISTORY_ITEM]: recordSearchHistoryItemReducer,
};

export default createReducer(reducerMapping, initialState);
