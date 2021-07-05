import PropTypes from 'prop-types';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { recordTelemetryEvent } from 'client/app/redux/actions/telemetry';

/**
 * HOC factory for injecting a Redux dispatcher for recording a client-side telemetry event.
 *
 * @param {Component} WrappedComponent Component to wrap.
 */
const withTelemetry = (WrappedComponent) => {
  const WithTelemetryHOC = ({ hocActions, ...props }) => (
    <WrappedComponent
      {...props}
      recordTelemetryEvent={hocActions.recordTelemetryEvent}
    />
  );

  WithTelemetryHOC.propTypes = {
    // HOC props
    hocActions: PropTypes.shape({
      recordTelemetryEvent: PropTypes.func.isRequired,
    }).isRequired,
  };

  const mapDispatchToProps = (dispatch) => ({
    hocActions: bindActionCreators({ recordTelemetryEvent }, dispatch),
  });

  return connect(null, mapDispatchToProps)(WithTelemetryHOC);
};

export default withTelemetry;
