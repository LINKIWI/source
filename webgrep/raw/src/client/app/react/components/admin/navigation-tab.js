import PropTypes from 'prop-types';
import React, { Fragment } from 'react';
import { colors, Spacing, Text } from 'react-elemental';
import { Link } from 'react-router-dom';
import { withToggleState } from '@linkiwi/hoc';
import { background } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

/**
 * Sidebar navigation tab for various admin pages.
 */
const AdminNavigationTab = ({
  label,
  href,
  isSelected,
  isHover,
  handleMouseEnter,
  handleMouseLeave,
}) => (
  <Fragment>
    <Link
      to={href}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      style={{
        display: 'block',
        height: '40px',
        opacity: (isHover || isSelected) ? 1 : 0.8,
        position: 'relative',
        textDecoration: 'none',
        transition: transition.all.DEFAULT,
      }}
    >
      <div
        style={{
          backgroundColor: colors.primary,
          height: '100%',
          position: 'absolute',
          transition: transition.all.DEFAULT,
          width: isSelected ? '3px' : 0,
        }}
      />

      <Spacing
        style={{
          alignItems: 'center',
          display: 'flex',
          flexGrow: 1,
          height: '100%',
          transition: transition.all.DEFAULT,
          ...isHover && {
            backgroundColor: background.dark.EPSILON,
          },
          ...isSelected && {
            backgroundColor: background.dark.GAMMA,
          },
        }}
        left
        right
        padding
      >
        <Text
          color={isSelected ? colors.primary : colors.gray80}
          size="kilo"
          bold={isSelected}
        >
          {label}
        </Text>
      </Spacing>
    </Link>
  </Fragment>
);

AdminNavigationTab.propTypes = {
  label: PropTypes.string.isRequired,
  href: PropTypes.string.isRequired,
  isSelected: PropTypes.bool.isRequired,
  // HOC props
  isHover: PropTypes.bool.isRequired,
  handleMouseEnter: PropTypes.func.isRequired,
  handleMouseLeave: PropTypes.func.isRequired,
};

export default withToggleState({
  key: 'isHover',
  enable: 'handleMouseEnter',
  disable: 'handleMouseLeave',
})(AdminNavigationTab);
