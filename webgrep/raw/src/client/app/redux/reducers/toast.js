import { SHOW_TOAST, HIDE_TOAST } from 'client/app/redux/actions/toast';
import createReducer from 'client/app/redux/reducers/create-reducer';

const initialState = {
  active: [],
};

const showToastReducer = (state, action) => ({
  ...state,
  active: [...state.active, action.payload]
    // Deduplicate by toast ID; allows updating a toast for an existing ID
    .filter(({ toastID }, idx, arr) =>
      arr.map((toast) => toast.toastID).lastIndexOf(toastID) === idx),
});

const hideToastReducer = (state, action) => ({
  ...state,
  active: state.active.filter(({ toastID }) => toastID !== action.payload.toastID),
});

const reducerMapping = {
  [SHOW_TOAST]: showToastReducer,
  [HIDE_TOAST]: hideToastReducer,
};

export default createReducer(reducerMapping, initialState);
