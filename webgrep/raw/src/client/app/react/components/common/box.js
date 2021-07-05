import PropTypes from 'prop-types';
import React from 'react';
import { withToggleState } from '@linkiwi/hoc';
import { Spacing } from 'react-elemental';
import { transition } from 'client/app/util/style/transition';

/**
 * Presentational wrapper around a block of arbitrary content.
 */
const Box = ({
  variant,
  style: styleOverrides,
  children,
  isHover,
  handleMouseEnter,
  handleMouseLeave,
}) => {
  const baseStyle = {
    backgroundColor: 'white',
    borderRadius: '3px',
    transition: transition.all.DEFAULT,
  };

  const styleVariants = {
    alpha: {
      ...baseStyle,
      boxShadow: isHover ?
        '0px 6px 24px -2px rgba(0, 0, 0, 0.125)' :
        '0px 6px 20px -2px rgba(0, 0, 0, 0.1)',
      ...styleOverrides,
    },
    beta: {
      ...baseStyle,
      boxShadow: isHover ?
        '0px 8px 24px -2px rgba(0, 0, 0, 0.075)' :
        '0px 8px 20px -2px rgba(0, 0, 0, 0.05)',
      ...styleOverrides,
    },
  };

  return (
    <Spacing
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      style={styleVariants[variant]}
      top
      right
      bottom
      left
      padding
    >
      {children}
    </Spacing>
  );
};

Box.propTypes = {
  variant: PropTypes.oneOf(['alpha', 'beta']),
  style: PropTypes.object,
  children: PropTypes.node.isRequired,
  // HOC props
  isHover: PropTypes.bool.isRequired,
  handleMouseEnter: PropTypes.func.isRequired,
  handleMouseLeave: PropTypes.func.isRequired,
};

Box.defaultProps = {
  variant: 'beta',
  style: {},
};

export default withToggleState({
  key: 'isHover',
  enable: 'handleMouseEnter',
  disable: 'handleMouseLeave',
})(Box);
