import PropTypes from 'prop-types';
import React from 'react';
import { Link, Spacing, Text } from 'react-elemental';
import { MdKeyboardArrowDown } from 'react-icons/md';
import { withDefaultPrevented } from 'client/app/util/navigation';
import { text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

/**
 * Conceal arbitrary children behind a visibility toggle.
 */
const Spoiler = ({ caption, isExpanded, onClick }) => (
  <Text color="primary" size="kilo" bold uppercase>
    <Link
      type="plain"
      activeColor={text.dark.ALPHA}
      onClick={withDefaultPrevented(onClick)}
      style={{ alignItems: 'center', display: 'flex' }}
    >
      <Spacing size="small" right inline>
        {caption}
      </Spacing>

      <MdKeyboardArrowDown
        style={{
          transform: isExpanded ? 'rotate(180deg)' : 'rotate(0deg)',
          transition: transition.all.DEFAULT,
        }}
      />
    </Link>
  </Text>
);

Spoiler.propTypes = {
  caption: PropTypes.string.isRequired,
  isExpanded: PropTypes.bool.isRequired,
  onClick: PropTypes.func.isRequired,
};

export default Spoiler;
