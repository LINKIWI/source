import PropTypes from 'prop-types';
import React from 'react';
import { Alert } from 'react-elemental';

/**
 * Abstraction over an Elemental Alert with customized defaults.
 */
const ErrorAlert = ({ title, message, ...props }) => (
  <Alert
    type="error"
    title={title}
    message={message}
    style={{ borderRadius: '3px' }}
    {...props}
  />
);

ErrorAlert.propTypes = {
  title: PropTypes.string,
  message: PropTypes.string.isRequired,
};

ErrorAlert.defaultProps = {
  title: 'Error',
};

export default ErrorAlert;
