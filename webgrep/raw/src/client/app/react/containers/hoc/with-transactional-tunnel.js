import { withTunnel } from 'supercharged/client';
import { parseSuperchargedResponse } from 'client/app/util/data';
import UIDGenerator from 'client/app/util/uid-generator';

/**
 * Factory for an HOC that creates a websocket tunnel that automatically handles request/response
 * "transactions" for messages.
 *
 * @param {Object} opts Options object passed directly to the underlying Supercharged tunnel client.
 * @returns {Function} HOC factory providing the same API as the Supercharged tunnel client.
 */
const withTransactionalTunnel = (opts) => (WrappedComponent) => {
  const transactionIDGenerator = new UIDGenerator();
  const WithTunnelHOC = withTunnel(opts)(WrappedComponent);

  return class WithTransactionalTunnelHOC extends WithTunnelHOC {
    txID = null;

    _onMessage({ data }) {
      const scResp = parseSuperchargedResponse(data);

      // Reject incoming messages if the relayed ID does not equal that of the most recent request,
      // to effectively cancel stale responses
      if (this.txID === scResp.id) {
        super._onMessage({ data: scResp });
      }
    }

    _sendMessage(message) {
      if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
        return;
      }

      this.txID = transactionIDGenerator.next();

      super._sendMessage(JSON.stringify({
        ...message,
        id: this.txID,
      }));
    }
  };
};

export default withTransactionalTunnel;
