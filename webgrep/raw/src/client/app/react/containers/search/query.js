import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { Spacing } from 'react-elemental';
import { connect } from 'react-redux';
import ConnectionStatus from 'client/app/react/components/search/meta/connection-status';
import QueryOptionControls from 'client/app/react/components/search/query/option-controls';
import QueryField from 'client/app/react/components/search/query/field';
import { objLookup } from 'shared/util/data';

/**
 * Wrapper for search query concerns, including the query entry itself and associated query options.
 */
class SearchQueryContainer extends Component {
  static propTypes = {
    value: PropTypes.string.isRequired,
    regex: PropTypes.bool.isRequired,
    caseSensitive: PropTypes.bool.isRequired,
    maxMatches: PropTypes.number.isRequired,
    filePath: PropTypes.string.isRequired,
    repositories: PropTypes.array.isRequired,
    filteredRepos: PropTypes.arrayOf(PropTypes.string.isRequired).isRequired,
    results: PropTypes.shape({
      err: PropTypes.object,
      data: PropTypes.object,
    }),
    onQueryChange: PropTypes.func.isRequired,
    onRegexChange: PropTypes.func.isRequired,
    onCaseSensitivityChange: PropTypes.func.isRequired,
    onMaxMatchesChange: PropTypes.func.isRequired,
    onFilePathChange: PropTypes.func.isRequired,
    onRepositoriesChange: PropTypes.func.isRequired,
    onQueryRecord: PropTypes.func.isRequired,
    connectivity: PropTypes.string.isRequired,
    isCompact: PropTypes.bool.isRequired,
  };

  static defaultProps = {
    results: null,
  };

  handleQueryFieldBlur = this._handleQueryFieldBlur.bind(this);

  handleRepositoryToggle = this._handleRepositoryToggle.bind(this);

  _handleQueryFieldBlur(evt) {
    if (evt.target.value) {
      this.props.onQueryRecord(evt.target.value);
    }
  }

  _handleRepositoryToggle(toggled) {
    const { filteredRepos: existingRepos, onRepositoriesChange } = this.props;

    // Reset repository filter list entirely
    if (!toggled) {
      return onRepositoriesChange([]);
    }

    const filteredRepos = new Set(existingRepos);

    if (filteredRepos.has(toggled)) {
      filteredRepos.delete(toggled);
    } else {
      filteredRepos.add(toggled);
    }

    return onRepositoriesChange([...filteredRepos]);
  }

  render() {
    const {
      value,
      regex,
      caseSensitive,
      maxMatches,
      filePath,
      repositories,
      results,
      filteredRepos,
      onQueryChange,
      onRegexChange,
      onCaseSensitivityChange,
      onMaxMatchesChange,
      onFilePathChange,
      connectivity,
      isCompact,
    } = this.props;

    const searchResults = [
      ...objLookup(results, ['data', 'files']) || [],
      ...objLookup(results, ['data', 'code']) || [],
    ];

    return (
      <Fragment>
        <Spacing size="20px" bottom>
          <QueryField
            value={value}
            onChange={onQueryChange}
            onBlur={this.handleQueryFieldBlur}
          />
        </Spacing>

        <div style={{ alignItems: 'center', display: 'flex', justifyContent: 'space-between' }}>
          <div style={{ flexGrow: 1 }}>
            <QueryOptionControls
              regex={regex}
              caseSensitive={caseSensitive}
              maxMatches={maxMatches}
              filePath={filePath}
              repositories={repositories}
              searchResults={searchResults}
              filteredRepos={filteredRepos}
              onRegexChange={onRegexChange}
              onCaseSensitivityChange={onCaseSensitivityChange}
              onMaxMatchesChange={onMaxMatchesChange}
              onFilePathChange={onFilePathChange}
              onRepositoryToggle={this.handleRepositoryToggle}
              isCompact={isCompact}
            />
          </div>

          {!isCompact && (
            <Spacing style={{ display: 'inherit' }} left>
              <ConnectionStatus connectivity={connectivity} />
            </Spacing>
          )}
        </div>
      </Fragment>
    );
  }
}

const mapStateToProps = ({ meta }) => ({
  repositories: Object.values(meta.repositories).map((repo) => ({
    name: repo.name,
    remote: repo.remote,
    labels: repo.labels,
  })),
});

export default connect(mapStateToProps)(SearchQueryContainer);
