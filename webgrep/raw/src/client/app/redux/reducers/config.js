import createReducer from 'client/app/redux/reducers/create-reducer';

const initialState = {
  // Injected by SSR store hydration
};

const reducerMapping = {};

export default createReducer(reducerMapping, initialState);
