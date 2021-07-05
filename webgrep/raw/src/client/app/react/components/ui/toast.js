import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { Spacing, Text } from 'react-elemental';
import { withToggleState } from '@linkiwi/hoc';
import { background, text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

/**
 * UI element for a single, plain-text toast.
 */
class Toast extends Component {
  static propTypes = {
    duration: PropTypes.number.isRequired,
    children: PropTypes.node.isRequired,
    // HOC props
    isRevealed: PropTypes.bool.isRequired,
    reveal: PropTypes.func.isRequired,
    hide: PropTypes.func.isRequired,
  };

  componentDidMount() {
    const { duration, reveal, hide } = this.props;

    this.timeout = window.setTimeout(reveal, 10);
    this.timeout = window.setTimeout(hide, duration - 120);
  }

  componentWillUnmount() {
    window.clearTimeout(this.timeout);
  }

  render() {
    const { isRevealed, children } = this.props;

    const baseStyle = {
      backgroundColor: background.dark.BETA,
      borderRadius: '3px',
      boxShadow: '0px 4px 8px -2px rgba(0, 0, 0, 0.3)',
      boxSizing: 'border-box',
      position: 'relative',
      transition: transition.all.DEFAULT,
      width: '100%',
    };

    const initialStyle = {
      opacity: 0,
    };

    const mountedStyle = {
      opacity: 0.85,
    };

    return (
      <Spacing
        size="18px"
        style={{
          ...baseStyle,
          ...isRevealed ? mountedStyle : initialStyle,
        }}
        padding
        top
        right
        left
        bottom
      >
        <Text size="kilo" color={text.light.BETA} style={{ wordBreak: 'break-word' }} center>
          {children}
        </Text>
      </Spacing>
    );
  }
}

export default withToggleState({
  key: 'isRevealed',
  enable: 'reveal',
  disable: 'hide',
})(Toast);
