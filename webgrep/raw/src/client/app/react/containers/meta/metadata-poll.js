import PropTypes from 'prop-types';
import { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { compose } from '@linkiwi/hoc';
import withTelemetry from 'client/app/react/containers/hoc/with-telemetry';
import withTransactionalTunnel from 'client/app/react/containers/hoc/with-transactional-tunnel';
import { setIndexMeta } from 'client/app/redux/actions/meta';
import { TELEMETRY_ACTIONS } from 'shared/constants/telemetry';

// Threshold number of consecutive, consistent metadata messages before considering the metadata
// accurate.
const METADATA_COMMIT_THRESHOLD = 3;

/**
 * Dummy container that polls index metadata at a regular intervals and updates global state.
 */
class MetadataPollContainer extends Component {
  static propTypes = {
    pollInterval: PropTypes.number.isRequired,
    // HOC props
    info: PropTypes.shape({
      err: PropTypes.object,
      messages: PropTypes.array.isRequired,
      isConnected: PropTypes.bool.isRequired,
      sendMessage: PropTypes.func.isRequired,
    }).isRequired,
    actions: PropTypes.shape({
      setIndexMeta: PropTypes.func.isRequired,
    }).isRequired,
    recordTelemetryEvent: PropTypes.func.isRequired,
  };

  componentDidMount() {
    this._tick();
  }

  componentDidUpdate() {
    const { info: { messages }, actions, recordTelemetryEvent } = this.props;

    if (messages.length !== METADATA_COMMIT_THRESHOLD) {
      return;
    }

    // Consider the metadata to be valid only if all messages in the buffer have a consistent index
    // name and timestamp. This effectively requires the server-side metadata to be consistent for
    // an extended duration of time, to accommodate clustered livegrep deployments which might be
    // serving different versions of the index at the same time while the index is being updated.
    const metadata = messages.reduce(
      (acc, message) => (
        (acc &&
          !message.err &&
          message.data &&
          acc.name === message.data.name &&
          acc.timestamp === message.data.timestamp) ? acc : null
      ),
      messages[0].data || {},
    );

    if (metadata) {
      actions.setIndexMeta(metadata);
      recordTelemetryEvent(TELEMETRY_ACTIONS.COMMIT_SERVER_METADATA);
    }
  }

  componentWillUnmount() {
    window.clearInterval(this.timeout);
  }

  _tick() {
    const { info, pollInterval, recordTelemetryEvent } = this.props;

    this.timeout = window.setTimeout(() => {
      if (!document.hidden) {
        info.sendMessage();
        recordTelemetryEvent(TELEMETRY_ACTIONS.POLL_SERVER_METADATA);
      }

      this._tick();
    }, pollInterval);
  }

  render() {
    return null;
  }
}

const mapDispatchToProps = (dispatch) => ({
  actions: bindActionCreators({ setIndexMeta }, dispatch),
});

export default compose(
  connect(null, mapDispatchToProps),
  withTelemetry,
  withTransactionalTunnel({
    key: 'info',
    endpoint: '/api/meta/info',
    messageBufferSize: METADATA_COMMIT_THRESHOLD,
  }),
)(MetadataPollContainer);
