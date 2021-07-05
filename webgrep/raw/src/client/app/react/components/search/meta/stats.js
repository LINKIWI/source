import PropTypes from 'prop-types';
import React from 'react';
import { pure } from 'recompose';
import { Spacing, Text } from 'react-elemental';
import { SearchStats } from 'shared/schemas';

/**
 * Visualization of server-side statistics for a search request.
 */
const SearchResultsStats = ({ numPaths, numFiles, time, exitReason }) => {
  const numPathsRepr = [
    numPaths,
    exitReason === SearchStats.ExitReason.MATCH_LIMIT && numPaths > 0 ? '+' : '',
  ].join('');

  const timeRepr = `(${time} ms)`;

  return (
    <div style={{ alignItems: 'center', display: 'flex' }}>
      <Spacing size="small" right>
        <Text color="gray60" size="kilo" bold>
          {numPathsRepr} paths, {numFiles} files
        </Text>
      </Spacing>

      <Text color="gray30" size="kilo">
        {timeRepr}
      </Text>
    </div>
  );
};

SearchResultsStats.propTypes = {
  numPaths: PropTypes.number.isRequired,
  numFiles: PropTypes.number.isRequired,
  time: PropTypes.number.isRequired,
  exitReason: PropTypes.oneOf(Object.values(SearchStats.ExitReason)).isRequired,
};

export default pure(SearchResultsStats);
