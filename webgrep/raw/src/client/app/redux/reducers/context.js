import { SET_WINDOW_DIMENSIONS } from 'client/app/redux/actions/context';
import createReducer from 'client/app/redux/reducers/create-reducer';

const initialState = {
  // Injected by SSR store hydration
};

const setWindowDimensionsReducer = (state, action) => {
  const { width, height } = action.payload;

  return {
    ...state,
    window: { width, height },
  };
};

const reducerMapping = {
  [SET_WINDOW_DIMENSIONS]: setWindowDimensionsReducer,
};

export default createReducer(reducerMapping, initialState);
