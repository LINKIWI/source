/* eslint-disable no-underscore-dangle */

import { compose, createStore, applyMiddleware } from 'redux';
import reducers from 'client/app/redux/reducers';
import middleware from 'client/app/redux/middleware';

const isDev = process.env.NODE_ENV !== 'production';
const composeEnhancers = (isDev && window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__) || compose;

export default createStore(
  reducers,
  window.__SSR_INJECTED_GLOBALS__.store,
  composeEnhancers(applyMiddleware(...middleware)),
);
