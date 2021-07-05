import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { withForm } from '@linkiwi/hoc';

/**
 * HOC factory for declaratively exposing the current URL hash state as component props.
 *
 * @param {Object} opts Options object describing semantic serialization/deserialization functions.
 * @returns {Function} HOC exposing the current hash value and a function to set the hash.
 */
const withHashState = (opts = {}) => (WrappedComponent) => {
  const {
    serializer = (value) => value,
    deserializer = (value) => value,
  } = opts;

  class WithHashStateHOC extends Component {
    static propTypes = {
      form: PropTypes.shape({
        hash: PropTypes.any,
      }).isRequired,
      handleFormChange: PropTypes.func.isRequired,
    };

    componentDidMount() {
      window.addEventListener('hashchange', this.handleHashChange);
    }

    componentWillUnmount() {
      window.removeEventListener('hashchange', this.handleHashChange);
    }

    setHash = this._setHash.bind(this);

    handleHashChange = this._handleHashChange.bind(this);

    _setHash(value) {  // eslint-disable-line class-methods-use-this
      window.location.hash = serializer(value);
    }

    _handleHashChange() {
      const hash = window.location.hash ? deserializer(window.location.hash) : null;
      this.props.handleFormChange('hash')({ target: { value: hash } });
    }

    render() {
      // eslint-disable-next-line no-unused-vars
      const { form: { hash }, handleFormChange, ...proxyProps } = this.props;

      const props = {
        ...proxyProps,
        hash,
        setHash: this.setHash,
      };

      return (
        <WrappedComponent {...props} />
      );
    }
  }

  return withForm({
    initial: () => ({
      hash: window.location.hash ? deserializer(window.location.hash) : null,
    }),
  })(WithHashStateHOC);
};

export default withHashState;
