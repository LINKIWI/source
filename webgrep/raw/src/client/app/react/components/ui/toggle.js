import PropTypes from 'prop-types';
import React from 'react';
import { colors, Button, Text } from 'react-elemental';
import { compose, withForwardedRef, withToggleState } from '@linkiwi/hoc';
import Tooltip from 'client/app/react/components/ui/tooltip';
import { background, text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

/**
 * Button-like toggle between active and inactive states.
 */
const Toggle = ({
  forwardedRef,
  description,
  secondary,
  onClick,
  isActive,
  isHover,
  isFocus,
  handleMouseEnter,
  handleMouseLeave,
  handleFocus,
  handleBlur,
  children,
}) => {
  const buttonStyle = {
    backgroundColor: background.primary.EPSILON,
    borderRadius: '3px',
    display: 'inline-flex',
    padding: '5px 10px',
    ...secondary && {
      backgroundColor: 'transparent',
      borderColor: background.primary.BETA,
      borderStyle: 'solid',
      borderWidth: '1px',
      padding: '4px 9px',
    },
    ...(isFocus || isHover) && {
      backgroundColor: background.primary.GAMMA,
    },
    ...isActive && {
      backgroundColor: background.primary.ALPHA,
    },
  };

  const textStyle = {
    color: colors.primary,
    display: 'inherit',
    transition: transition.all.DEFAULT,
    ...isActive && {
      color: text.light.ALPHA,
    },
  };

  const button = (
    <Button
      ref={forwardedRef}
      style={buttonStyle}
      onClick={onClick}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={compose(handleMouseLeave, handleBlur)}
      onFocus={handleFocus}
      onBlur={handleBlur}
    >
      <Text size="11px" style={textStyle} uppercase={!secondary} bold inline>
        {children}
      </Text>
    </Button>
  );

  if (!description) {
    return button;
  }

  return (
    <Tooltip
      description={description}
      style={{ minWidth: '80px' }}
    >
      {button}
    </Tooltip>
  );
};

Toggle.propTypes = {
  description: PropTypes.string,
  secondary: PropTypes.bool,
  onClick: PropTypes.func.isRequired,
  isActive: PropTypes.bool,
  children: PropTypes.node.isRequired,
  // HOC props
  forwardedRef: PropTypes.oneOfType([
    PropTypes.shape({ current: PropTypes.instanceOf(Element) }),
    PropTypes.func,
  ]),
  isHover: PropTypes.bool.isRequired,
  isFocus: PropTypes.bool.isRequired,
  handleMouseEnter: PropTypes.func.isRequired,
  handleMouseLeave: PropTypes.func.isRequired,
  handleFocus: PropTypes.func.isRequired,
  handleBlur: PropTypes.func.isRequired,
};

Toggle.defaultProps = {
  forwardedRef: null,
  description: null,
  secondary: false,
  isActive: false,
};

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
)(Toggle);
