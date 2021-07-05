import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { cycleToast } from 'client/app/redux/actions/toast';

/**
 * HOC factory for injecting an asynchronous Redux action dispatcher to show and hide a toast.
 *
 * @param {Component} WrappedComponent Component to wrap.
 */
const withToast = (WrappedComponent) => {
  class ToastProviderHOC extends Component {
    static propTypes = {
      // HOC props
      dispatch: PropTypes.func.isRequired,
    };

    toast = this._toast.bind(this);

    _toast(text, timeout) {
      this.props.dispatch(cycleToast(text, timeout));
    }

    render() {
      return (
        <WrappedComponent
          {...this.props}
          toast={this.toast}
        />
      );
    }
  }

  return connect()(ToastProviderHOC);
};

export default withToast;
