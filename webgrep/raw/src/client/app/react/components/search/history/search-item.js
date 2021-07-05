import PropTypes from 'prop-types';
import React from 'react';
import { MdSearch } from 'react-icons/md';
import { Link, Spacing, Text } from 'react-elemental';
import { background, text } from 'client/app/util/style/color';

/**
 * Represents a single search history item.
 */
const SearchItem = ({ query, href, onClick }) => (
  <Link
    type="plain"
    href={href}
    onClick={onClick}
    style={{ display: 'inline-block', textDecoration: 'none' }}
  >
    <div
      style={{
        alignItems: 'center',
        backgroundColor: background.dark.GAMMA,
        borderRadius: '3px',
        display: 'flex',
        padding: '10px 14px',
      }}
    >
      <Spacing size="small" style={{ display: 'inherit' }} right>
        <MdSearch style={{ opacity: 0.2 }} />
      </Spacing>

      <Text color={text.dark.BETA} size="kilo">
        {query}
      </Text>
    </div>
  </Link>
);

SearchItem.propTypes = {
  query: PropTypes.string.isRequired,
  href: PropTypes.string.isRequired,
  onClick: PropTypes.func.isRequired,
};

export default SearchItem;
