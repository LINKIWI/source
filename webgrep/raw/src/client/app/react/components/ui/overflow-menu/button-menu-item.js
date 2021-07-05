import PropTypes from 'prop-types';
import React from 'react';
import { Spacing, Text } from 'react-elemental';
import { compose, withToggleState } from '@linkiwi/hoc';
import { background, text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

/**
 * Single overflow menu item as a native button element.
 */
const ButtonMenuItem = ({
  label,
  icon,
  onClick,
  handleMouseEnter,
  handleMouseLeave,
  handleMouseDown,
  handleMouseUp,
  isHover,
  isClicked,
}) => {
  const baseStyle = {
    alignItems: 'center',
    backgroundColor: 'white',
    border: 0,
    cursor: 'pointer',
    display: 'flex',
    padding: '10px 16px',
    transition: transition.all.DEFAULT,
    width: '100%',
  };

  const hoverStyle = {
    backgroundColor: background.primary.EPSILON,
  };

  const clickStyle = {
    backgroundColor: text.dark.GAMMA,
  };

  return (
    <button
      onClick={onClick}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      onMouseDown={handleMouseDown}
      onMouseUp={handleMouseUp}
      style={{
        ...baseStyle,
        ...isHover && hoverStyle,
        ...isHover && isClicked && clickStyle,
      }}
    >
      {icon && (
        <Spacing size="small" right>
          <Text color={text.dark.BETA} size="kilo">
            {icon}
          </Text>
        </Spacing>
      )}
      <Text color={text.dark.ALPHA} size="kilo" style={{ textAlign: 'left' }}>
        {label}
      </Text>
    </button>
  );
};

ButtonMenuItem.propTypes = {
  label: PropTypes.string.isRequired,
  icon: PropTypes.node,
  onClick: PropTypes.func,
  // HOC props
  handleMouseEnter: PropTypes.func.isRequired,
  handleMouseLeave: PropTypes.func.isRequired,
  handleMouseDown: PropTypes.func.isRequired,
  handleMouseUp: PropTypes.func.isRequired,
  isHover: PropTypes.bool.isRequired,
  isClicked: PropTypes.bool.isRequired,
};

ButtonMenuItem.defaultProps = {
  onClick: () => {},
  icon: null,
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
)(ButtonMenuItem);
