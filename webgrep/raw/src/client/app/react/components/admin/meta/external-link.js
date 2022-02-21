import PropTypes from 'prop-types';
import React from 'react';
import { colors, Link, Spacing, Text } from 'react-elemental';
import { MdNavigateNext } from 'react-icons/md';
import { withToggleState } from '@linkiwi/hoc';
import { text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

/**
 * Representation of an external link resource.
 */
const ExternalLink = ({ title, href, isHover, handleMouseEnter, handleMouseLeave }) => (
  <Link
    type="plain"
    href={href}
    onMouseEnter={handleMouseEnter}
    onMouseLeave={handleMouseLeave}
    style={{ alignItems: 'center', display: 'flex', justifyContent: 'space-between' }}
  >
    <div>
      <Spacing size="micro" bottom>
        <Text size="kilo" bold>
          {title}
        </Text>
      </Spacing>

      <Text color={colors.primary} style={{ wordBreak: 'break-word' }} size="kilo">
        {href}
      </Text>
    </div>

    <Spacing left>
      <MdNavigateNext
        style={{
          color: isHover ? colors.primary : text.dark.GAMMA,
          marginRight: isHover ? '-6px' : 0,
          transition: transition.all.DEFAULT,
        }}
      />
    </Spacing>
  </Link>
);

ExternalLink.propTypes = {
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
})(ExternalLink);
