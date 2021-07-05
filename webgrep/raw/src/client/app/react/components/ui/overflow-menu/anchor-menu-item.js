import PropTypes from 'prop-types';
import React from 'react';
import { Spacing, Text } from 'react-elemental';
import { compose, withToggleState } from '@linkiwi/hoc';
import { background, text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

/**
 * Single overflow menu item acting as a native anchor element.
 */
const AnchorMenuItem = ({
  label,
  icon,
  href,
  onClick,
  handleMouseEnter,
  handleMouseLeave,
  handleMouseDown,
  handleMouseUp,
  isHover,
  isClicked,
}) => {
  const baseStyle = {
    backgroundColor: 'white',
    border: 0,
    cursor: 'pointer',
    textAlign: 'left',
    transition: transition.all.DEFAULT,
    width: '100%',
  };

  const hoverStyle = {
    backgroundColor: background.primary.EPSILON,
  };

  const clickStyle = {
    backgroundColor: text.dark.GAMMA,
  };

  const anchorStyle = {
    alignItems: 'center',
    color: 'currentColor',
    display: 'flex',
    padding: '10px 16px',
    textDecoration: 'none',
    width: '100%',
  };

  return (
    <Text
      color={text.dark.ALPHA}
      size="kilo"
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
      <a href={href} onClick={onClick} style={anchorStyle}>
        {icon && (
          <Spacing size="small" style={{ display: 'inherit' }} right inline>
            <Text color={text.dark.BETA} size="kilo" inline>
              {icon}
            </Text>
          </Spacing>
        )}
        {label}
      </a>
    </Text>
  );
};

AnchorMenuItem.propTypes = {
  label: PropTypes.string.isRequired,
  icon: PropTypes.node,
  href: PropTypes.string.isRequired,
  onClick: PropTypes.func,
  // HOC props
  handleMouseEnter: PropTypes.func.isRequired,
  handleMouseLeave: PropTypes.func.isRequired,
  handleMouseDown: PropTypes.func.isRequired,
  handleMouseUp: PropTypes.func.isRequired,
  isHover: PropTypes.bool.isRequired,
  isClicked: PropTypes.bool.isRequired,
};

AnchorMenuItem.defaultProps = {
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
)(AnchorMenuItem);
