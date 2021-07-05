import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { withToggleState } from '@linkiwi/hoc';
import { transition } from 'client/app/util/style/transition';

/**
 * Wrapper component to delay visibility of its children after initial mount.
 */
class Delayed extends Component {
  static propTypes = {
    delayMs: PropTypes.number,
    children: PropTypes.node.isRequired,
    // HOC props
    isVisible: PropTypes.bool.isRequired,
    show: PropTypes.func.isRequired,
  };

  static defaultProps = {
    delayMs: 150,
  };

  componentDidMount() {
    const { delayMs, show } = this.props;

    this.timeout = setTimeout(show, delayMs);
  }

  componentWillUnmount() {
    clearTimeout(this.timeout);
  }

  render() {
    const { isVisible, children } = this.props;

    const style = {
      opacity: 0,
      transition: transition.opacity.DEFAULT,
      ...isVisible && {
        opacity: 1,
      },
    };

    return (
      <div style={style}>
        {children}
      </div>
    );
  }
}

export default withToggleState({
  key: 'isVisible',
  enable: 'show',
})(Delayed);
