import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { sizes, spacing, TextField as ElementalTextField } from 'react-elemental';
import { compose, withForwardedRef, withToggleState } from '@linkiwi/hoc';
import { background, text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';
import { omit } from 'shared/util/data';

/**
 * Styling and enhancer abstractions over an Elemental TextField.
 */
class TextField extends Component {
  static propTypes = {
    before: PropTypes.node,
    after: PropTypes.node,
    onBlur: PropTypes.func,
    beforeEnhancerStyle: PropTypes.object,
    afterEnhancerStyle: PropTypes.object,
    // HOC props
    forwardedRef: PropTypes.oneOfType([
      PropTypes.shape({ current: PropTypes.instanceOf(Element) }),
      PropTypes.func,
    ]),
    isHover: PropTypes.bool.isRequired,
    handleMouseEnter: PropTypes.func.isRequired,
    handleMouseLeave: PropTypes.func.isRequired,
    isFocus: PropTypes.bool.isRequired,
    handleFocus: PropTypes.func.isRequired,
    handleBlur: PropTypes.func.isRequired,
  };

  static defaultProps = {
    before: null,
    after: null,
    beforeEnhancerStyle: {},
    afterEnhancerStyle: {},
    onBlur: () => {},
    forwardedRef: null,
  };

  handleBlur = this._handleBlur.bind(this);

  _handleBlur(evt) {
    const { onBlur, handleBlur } = this.props;

    // Dispatch the event to both the parent-specified blur handler and toggle state handler
    onBlur(evt);
    handleBlur(evt);
  }

  render() {
    const {
      forwardedRef,
      before,
      after,
      beforeEnhancerStyle: beforeEnhancerStyleOverrides,
      afterEnhancerStyle: afterEnhancerStyleOverrides,
      isFocus,
      isHover,
      handleFocus,
      handleMouseEnter,
      handleMouseLeave,
      ...props
    } = this.props;

    const proxyProps = omit(props, ['handleBlur', 'onBlur']);

    const containerStyle = {
      position: 'relative',
    };

    const textFieldStyle = {
      backgroundColor: background.dark.EPSILON,
      border: '1px solid transparent',
      borderRadius: '3px',
      color: text.dark.BETA,
      fontSize: sizes.iota,
      height: '40px',
      padding: '12px',
      ...before && {
        paddingLeft: '46px',
      },
      ...after && {
        paddingRight: '46px',
      },
      ...isHover && {
        backgroundColor: background.dark.IOTA,
        color: text.dark.ALPHA,
      },
      ...isFocus && {
        backgroundColor: 'transparent',
        border: `1px solid ${background.primary.BETA}`,
        color: text.dark.ALPHA,
      },
    };

    const baseEnhancerStyle = {
      alignItems: 'center',
      display: 'flex',
      height: '100%',
      justifyContent: 'center',
      opacity: 0.2,
      pointerEvents: 'none',
      position: 'absolute',
      transition: transition.all.DEFAULT,
      width: spacing.default,
      ...isHover && {
        opacity: 0.4,
      },
      ...isFocus && {
        opacity: 0.7,
      },
    };

    const beforeEnhancerStyle = {
      ...baseEnhancerStyle,
      marginLeft: '14px',
      ...beforeEnhancerStyleOverrides,
    };

    const afterEnhancerStyle = {
      ...baseEnhancerStyle,
      marginRight: '14px',
      position: 'absolute',
      right: 0,
      ...afterEnhancerStyleOverrides,
    };

    return (
      <div
        onMouseEnter={handleMouseEnter}
        onMouseLeave={handleMouseLeave}
        style={containerStyle}
      >
        {before && (
          <div style={beforeEnhancerStyle}>
            {before}
          </div>
        )}

        {after && (
          <div style={afterEnhancerStyle}>
            {after}
          </div>
        )}

        <ElementalTextField
          ref={forwardedRef}
          onFocus={handleFocus}
          onBlur={this.handleBlur}
          isFocus={isFocus}
          style={textFieldStyle}
          {...proxyProps}
        />
      </div>
    );
  }
}

export default compose(
  withForwardedRef,
  withToggleState({
    key: 'isHover',
    enable: 'handleMouseEnter',
    disable: 'handleMouseLeave',
  }),
  withToggleState({
    key: 'isFocus',
    enable: 'handleFocus',
    disable: 'handleBlur',
  }),
)(TextField);
