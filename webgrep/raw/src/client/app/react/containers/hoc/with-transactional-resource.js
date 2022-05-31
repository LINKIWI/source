import { withResource } from 'supercharged/client';
import UIDGenerator from 'client/app/util/uid-generator';

/**
 * Factory for an HOC that wraps a resource fetcher that automatically handles request/response
 * "transactions" across invocations.
 *
 * @param {Object} opts Options object passed directly to the underlying Supercharged resource
 *                      client.
 * @returns {Function} HOC factory providing the same API as the Supercharged resource client.
 */
const withTransactionalResource = (opts) => (WrappedComponent) => {
  const transactionIDGenerator = new UIDGenerator();
  const WithResourceHOC = withResource(opts)(WrappedComponent);

  return class WithTransactionalResourceHOC extends WithResourceHOC {
    txID = null;

    shouldComponentUpdate(nextProps, nextState) {
      const { respData } = this.state;
      const { respData: nextRespData } = nextState;

      // Transitioning empty data to populated data; ensure that the first response is consistent
      // with the current client transaction ID.
      if (!respData && nextRespData) {
        if (nextRespData.id === undefined) {
          return true;
        }

        if (nextRespData.id !== this.txID) {
          return false;
        }
      }

      // Transitioning from one populated data state to another; transaction consistency only needs
      // to be checked if the server-reported response has changed since the last render.
      if (respData && nextRespData) {
        if (nextRespData.id === undefined) {
          return true;
        }

        if (respData.id !== nextRespData.id) {
          if (nextRespData.id !== this.txID) {
            return false;
          }
        }
      }

      return true;
    }

    componentDidUpdate(prevProps, prevState) {
      // After a transition from a non-error state to an error state, reset the transaction ID so
      // that any potential future populated responses are rejected. This is fundamentally a
      // workaround for the fact that the transaction ID is not relayed to the client for error
      // responses.
      if (!prevState.err && this.state.err) {
        transactionIDGenerator.reset();
        this.txID = transactionIDGenerator.next();
      }
    }

    _invokeResource(data, cb) {
      this.txID = transactionIDGenerator.next();
      super._invokeResource({ id: this.txID, ...data }, cb);
    }
  };
};

export default withTransactionalResource;
