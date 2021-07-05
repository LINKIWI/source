/**
 * When no middleware function matches the action type, simply continue the middleware chain by
 * reducing the current action.
 */
const cascade = (store, next, action) => next(action);

/**
 * Create a middleware function from a mapping of action types to middleware functions, each taking
 * as input the store and the triggering action.
 *
 * @param {Object} middlewareMapping Map of action type names to middleware functions.
 */
const createMiddleware = (middlewareMapping) => (store) => (next) => (action) =>
  (middlewareMapping[action.type] || cascade)(store, next, action);

export default createMiddleware;
