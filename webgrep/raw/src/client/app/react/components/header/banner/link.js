import PropTypes from 'prop-types';
import React from 'react';
import { Link, Spacing, Text } from 'react-elemental';
import { MdNavigateNext } from 'react-icons/md';
import { withToggleState } from '@linkiwi/hoc';
import { text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

const BannerLink = ({ title, href, isHover, handleMouseEnter, handleMouseLeave }) => (
  <Link
    type="plain"
    activeColor={text.light.ALPHA}
    href={href}
    onMouseEnter={handleMouseEnter}
    onMouseLeave={handleMouseLeave}
    style={{ alignItems: 'center', display: 'inline-flex', justifyContent: 'space-between' }}
  >
    <Text color={text.light.BETA} inline>
      {title}
    </Text>

    <Spacing size="micro" style={{ display: 'inherit' }} left>
      <MdNavigateNext
        style={{
          color: text.light.ALPHA,
          marginLeft: isHover ? '4px' : 0,
          transition: transition.all.DEFAULT,
        }}
      />
    </Spacing>
  </Link>
);

BannerLink.propTypes = {
  title: PropTypes.string.isRequired,
  href: PropTypes.string.isRequired,
  // HOC props
  isHover: PropTypes.bool.isRequired,
  handleMouseEnter: PropTypes.func.isRequired,
  handleMouseLeave: PropTypes.func.isRequired,
};

export default withToggleState({
  key: 'isHover',
  enable: 'handleMouseEnter',
  disable: 'handleMouseLeave',
})(BannerLink);
