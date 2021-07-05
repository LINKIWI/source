import * as Sentry from '@sentry/browser';

/**
 * Middleware that injects context breadcrumbs detailing dispatched Redux actions and corresponding
 * state reductions, to provide context for potential reported errors.
 */
const errorContext = (store) => (next) => (action) => {
  Sentry.addBreadcrumb({
    message: 'State reduction in response to dispatched action',
    category: 'redux',
    data: {
      action,
      store: store.getState(),
    },
  });

  return next(action);
};

export default errorContext;
