import PropTypes from 'prop-types';
import React from 'react';
import { Spacing, Text } from 'react-elemental';
import { compose, withToggleState } from '@linkiwi/hoc';
import { background, text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

/**
 * Floating button with text and an icon.
 */
const FloatingButton = ({
  value,
  icon,
  onClick,
  isHover,
  isClicked,
  handleMouseEnter,
  handleMouseLeave,
  handleMouseDown,
  handleMouseUp,
}) => {
  const baseStyle = {
    alignItems: 'center',
    backgroundColor: 'rgba(0, 0, 0, 0.7)',
    border: 'none',
    borderRadius: '3px',
    boxShadow: '0px 4px 6px -2px rgba(0, 0, 0, 0.2)',
    cursor: 'pointer',
    display: 'inline-flex',
    padding: '8px 12px',
    transition: transition.all.DEFAULT,
    whiteSpace: 'nowrap',
  };

  const hoverStyle = {
    backgroundColor: background.dark.BETA,
    boxShadow: '0px 4px 8px -2px rgba(0, 0, 0, 0.3)',
  };

  const clickStyle = {
    backgroundColor: background.dark.ALPHA,
    boxShadow: '0px 4px 10px -2px rgba(0, 0, 0, 0.5)',
  };

  const style = {
    ...baseStyle,
    ...isHover && hoverStyle,
    ...isHover && isClicked && clickStyle,
  };

  return (
    <button
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      onMouseDown={handleMouseDown}
      onMouseUp={handleMouseUp}
      style={style}
      onClick={onClick}
    >
      <Spacing size="small" right>
        <Text color={text.light.BETA} size="11px" uppercase bold>
          {value}
        </Text>
      </Spacing>

      <Text color={text.light.BETA} size="kilo" style={{ display: 'inherit' }}>
        {icon}
      </Text>
    </button>
  );
};

FloatingButton.propTypes = {
  value: PropTypes.string.isRequired,
  icon: PropTypes.node.isRequired,
  onClick: PropTypes.func.isRequired,
  // HOC props
  isHover: PropTypes.bool.isRequired,
  isClicked: PropTypes.bool.isRequired,
  handleMouseEnter: PropTypes.func.isRequired,
  handleMouseLeave: PropTypes.func.isRequired,
  handleMouseDown: PropTypes.func.isRequired,
  handleMouseUp: PropTypes.func.isRequired,
};

export default compose(
  withToggleState({
    key: 'isHover',
    enable: 'handleMouseEnter',
    disable: 'handleMouseLeave',
  }),
  withToggleState({
    key: 'isClicked',
    enable: 'handleMouseDown',
    disable: 'handleMouseUp',
  }),
)(FloatingButton);
