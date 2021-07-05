import PropTypes from 'prop-types';
import React, { Fragment } from 'react';
import { Spacing } from 'react-elemental';
import { withToggleState } from '@linkiwi/hoc';
import FileEntryContainer from 'client/app/react/containers/search/results/file-entry';
import Spoiler from 'client/app/react/components/ui/spoiler';

/**
 * List of results for matching files.
 */
const FileResultsContainer = ({
  files,
  repositories,
  fileResultsLimit,
  windowWidth,
  isExpanded,
  expand,
  contract,
  onSearchQueryChange,
  onSearchRepositoryFilePathChange,
}) => (
  <Fragment>
    <Spacing size="small" bottom={files.length > fileResultsLimit}>
      {files.slice(0, !isExpanded ? fileResultsLimit : undefined).map((file, idx) => (
        <Spacing
          key={`${file.repo}-${file.path}-${file.bounds[0]}-${file.bounds[1]}`}
          size="tiny"
          bottom={idx < files.length - 1}
          style={{ display: 'flex' }}
        >
          <FileEntryContainer
            file={file}
            repositories={repositories}
            position={idx + 1}
            windowWidth={windowWidth}
            onSearchQueryChange={onSearchQueryChange}
            onSearchRepositoryFilePathChange={onSearchRepositoryFilePathChange}
          />
        </Spacing>
      ))}
    </Spacing>

    {files.length > fileResultsLimit && (
      <div style={{ display: 'inline-block' }}>
        <Spoiler
          caption={isExpanded ? 'Show less' : 'Show all'}
          onClick={isExpanded ? contract : expand}
          isExpanded={isExpanded}
        />
      </div>
    )}
  </Fragment>
);

FileResultsContainer.propTypes = {
  files: PropTypes.array.isRequired,
  repositories: PropTypes.object.isRequired,
  fileResultsLimit: PropTypes.number.isRequired,
  windowWidth: PropTypes.number.isRequired,
  onSearchQueryChange: PropTypes.func.isRequired,
  onSearchRepositoryFilePathChange: PropTypes.func.isRequired,
  // HOC props
  isExpanded: PropTypes.bool.isRequired,
  expand: PropTypes.func.isRequired,
  contract: PropTypes.func.isRequired,
};

export default withToggleState({
  key: 'isExpanded',
  enable: 'expand',
  disable: 'contract',
})(FileResultsContainer);
