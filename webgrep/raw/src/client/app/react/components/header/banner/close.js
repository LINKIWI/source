import PropTypes from 'prop-types';
import React from 'react';
import { sizes, Button } from 'react-elemental';
import { MdClear } from 'react-icons/md';
import { withToggleState } from '@linkiwi/hoc';
import { text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

const BannerCloseButton = ({ onClick, isHover, handleMouseEnter, handleMouseLeave }) => (
  <Button
    color="transparent"
    onClick={onClick}
    onMouseEnter={handleMouseEnter}
    onMouseLeave={handleMouseLeave}
    style={{
      alignItems: 'center',
      display: 'flex',
      justifyContent: 'center',
      padding: '6px',
    }}
  >
    <MdClear
      style={{
        color: isHover ? text.light.GAMMA : text.light.EPSILON,
        fontSize: sizes.kilo,
        transition: transition.all.DEFAULT,
      }}
    />
  </Button>
);

BannerCloseButton.propTypes = {
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
})(BannerCloseButton);
