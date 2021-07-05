import PropTypes from 'prop-types';
import React, { Fragment } from 'react';
import { Spacing, Text } from 'react-elemental';
import SearchItem from 'client/app/react/components/search/history/search-item';
import { withDefaultPrevented } from 'client/app/util/navigation';

/**
 * Display of a list of recent search queries.
 */
const RecentSearches = ({ history, onHistoricalSearch }) => (
  <Fragment>
    <Spacing bottom>
      <Text size="kilo" bold>
        Recent searches
      </Text>
    </Spacing>

    {history.map((query, idx) => (
      <Spacing key={query} size="small" bottom={idx < history.length - 1}>
        <SearchItem
          query={query}
          href={`/search?query=${query}`}
          onClick={withDefaultPrevented(() => onHistoricalSearch(query))}
        />
      </Spacing>
    ))}
  </Fragment>
);

RecentSearches.propTypes = {
  history: PropTypes.arrayOf(PropTypes.string.isRequired).isRequired,
  onHistoricalSearch: PropTypes.func.isRequired,
};

export default RecentSearches;
