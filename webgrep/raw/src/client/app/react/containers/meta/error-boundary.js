import PropTypes from 'prop-types';
import { Component } from 'react';

/**
 * Isolates a tree of children with an error boundary against uncaught exceptions.
 */
export default class ErrorBoundaryContainer extends Component {
  static propTypes = {
    fallback: PropTypes.func,
    children: PropTypes.node.isRequired,
  };

  static defaultProps = {
    fallback: () => null,
  };

  state = {
    error: null,
  };

  static getDerivedStateFromError(error) {
    return { error };
  }

  render() {
    const { fallback, children } = this.props;
    const { error } = this.state;

    return error ? fallback(error) : children;
  }
}
