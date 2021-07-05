import PropTypes from 'prop-types';
import React from 'react';
import { sizes, Button } from 'react-elemental';
import { withToggleState } from '@linkiwi/hoc';
import { background, text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

/**
 * Button wrapping an arbitrary icon class.
 */
const IconButton = ({
  icon: Icon,
  onClick,
  isHover,
  handleMouseEnter,
  handleMouseLeave,
}) => {
  const buttonStyle = {
    alignItems: 'center',
    borderRadius: '3px',
    display: 'flex',
    justifyContent: 'center',
  };

  const iconStyle = {
    color: text.dark.GAMMA,
    fontSize: sizes.kilo,
    transition: transition.all.DEFAULT,
    ...isHover && {
      color: text.dark.BETA,
    },
  };

  return (
    <Button
      size="gamma"
      color={background.dark.GAMMA}
      style={buttonStyle}
      onClick={onClick}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
    >
      <Icon style={iconStyle} />
    </Button>
  );
};

IconButton.propTypes = {
  icon: PropTypes.elementType.isRequired,
  onClick: PropTypes.func.isRequired,
  // HOC props
  isHover: PropTypes.bool.isRequired,
  handleMouseEnter: PropTypes.func.isRequired,
  handleMouseLeave: PropTypes.func.isRequired,
};

export default withToggleState({
  key: 'isHover',
  enable: 'handleMouseEnter',
  disable: 'handleMouseLeave',
})(IconButton);
