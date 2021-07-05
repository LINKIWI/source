import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { Button, Spacing, Text } from 'react-elemental';
import { connect } from 'react-redux';
import { onlyUpdateForKeys } from 'recompose';
import { compose } from '@linkiwi/hoc';
import { CODE_INVALID_PARAMETERS } from 'supercharged/shared/constants/error';
import withTelemetry from 'client/app/react/containers/hoc/with-telemetry';
import CodeResultsContainer from 'client/app/react/containers/search/results/code-results';
import FileResultsContainer from 'client/app/react/containers/search/results/file-results';
import RecentSearches from 'client/app/react/components/search/history/recent-searches';
import SearchResultsStats from 'client/app/react/components/search/meta/stats';
import ErrorAlert from 'client/app/react/components/ui/error-alert';
import { PREFERENCE_KEYS } from 'client/app/util/constants/preferences';
import { MAX_MATCHES } from 'client/app/util/constants/search';
import { background, text } from 'client/app/util/style/color';
import { TELEMETRY_ACTIONS } from 'shared/constants/telemetry';
import { SearchStats } from 'shared/schemas';
import { objLookup } from 'shared/util/data';

/**
 * Logic container for rendering of both file and code search results.
 */
class SearchResultsContainer extends Component {
  static propTypes = {
    results: PropTypes.shape({
      err: PropTypes.object,
      data: PropTypes.object,
    }),
    maxMatches: PropTypes.number.isRequired,
    onMaxMatchesChange: PropTypes.func.isRequired,
    onFilePathChange: PropTypes.func.isRequired,
    onRepositoriesChange: PropTypes.func.isRequired,
    onQueryChange: PropTypes.func.isRequired,
    onQueryRecord: PropTypes.func.isRequired,
    // HOC props
    windowWidth: PropTypes.number.isRequired,
    repositories: PropTypes.object.isRequired,
    recordTelemetryEvent: PropTypes.func.isRequired,
    history: PropTypes.array.isRequired,
    fileResultsLimit: PropTypes.number.isRequired,
  };

  static defaultProps = {
    results: null,
  };

  executeHistoricalSearch = this._executeHistoricalSearch.bind(this);

  executeSourcePreviewSearch = this._executeSourcePreviewSearch.bind(this);

  executeSingleFileSearch = this._executeSingleFileSearch.bind(this);

  executeLoadMore = this._executeLoadMore.bind(this);

  _executeHistoricalSearch(query) {
    const { recordTelemetryEvent, onQueryChange } = this.props;

    recordTelemetryEvent(TELEMETRY_ACTIONS.EXECUTE_SEARCH_RECENT);
    onQueryChange({ target: { value: query } });
  }

  _executeSourcePreviewSearch(query) {
    const { recordTelemetryEvent, onQueryChange, onQueryRecord } = this.props;

    recordTelemetryEvent(TELEMETRY_ACTIONS.EXECUTE_SEARCH_SOURCE_PREVIEW);
    onQueryChange({ target: { value: query } });
    // Also manually persist the query to search history, since it would not otherwise be triggered
    // as when entering the query in the query field.
    onQueryRecord(query);
  }

  _executeSingleFileSearch(repo, path) {
    const { recordTelemetryEvent, onRepositoriesChange, onFilePathChange } = this.props;

    recordTelemetryEvent(TELEMETRY_ACTIONS.EXECUTE_SEARCH_SINGLE_FILE, 1, { repo });
    // Explicitly set an update completion callback to mitigate an async state update race condition
    // for different details in the same parent container.
    onRepositoriesChange([repo], () => onFilePathChange({ target: { value: path } }));
  }

  _executeLoadMore() {
    this.props.recordTelemetryEvent(TELEMETRY_ACTIONS.LOAD_MORE_MATCH_RESULTS);
    return this.props.onMaxMatchesChange(this.props.maxMatches + MAX_MATCHES.INCREMENT);
  }

  render() {
    const {
      results,
      repositories,
      history,
      fileResultsLimit,
      windowWidth,
    } = this.props;

    const data = objLookup(results, ['data']) || {
      files: [],
      code: [],
      stats: { time: 0, exitReason: SearchStats.ExitReason.NONE },
    };

    // A request/response transaction has not yet been initiated, or the client supplied an empty
    // search query
    if (!results || (results.err && results.err.code === CODE_INVALID_PARAMETERS)) {
      return history.length ? (
        <div style={{ alignItems: 'center', display: 'flex', flexDirection: 'column' }}>
          <RecentSearches
            history={history}
            onHistoricalSearch={this.executeHistoricalSearch}
          />
        </div>
      ) : null;
    }

    if (results.err) {
      return (
        <ErrorAlert message={results.err.message} />
      );
    }

    return (
      <Fragment>
        <Spacing size="large" bottom>
          <SearchResultsStats
            numPaths={data.files.length}
            numFiles={data.code.length}
            time={data.stats.time}
            exitReason={data.stats.exitReason}
          />
        </Spacing>

        {fileResultsLimit > 0 && (
          <Spacing size="large" bottom>
            <FileResultsContainer
              files={data.files}
              repositories={repositories}
              fileResultsLimit={fileResultsLimit}
              windowWidth={windowWidth}
              onSearchQueryChange={this.executeSourcePreviewSearch}
              onSearchRepositoryFilePathChange={this.executeSingleFileSearch}
            />
          </Spacing>
        )}

        <CodeResultsContainer
          snippets={data.code}
          repositories={repositories}
          onSearchQueryChange={this.executeSourcePreviewSearch}
          onSearchRepositoryFilePathChange={this.executeSingleFileSearch}
        />

        {data.stats.exitReason === SearchStats.ExitReason.MATCH_LIMIT && (
          <Spacing size="large" style={{ display: 'flex', justifyContent: 'center' }} top>
            <Button
              color={background.dark.GAMMA}
              style={{ borderRadius: '3px' }}
              onClick={this.executeLoadMore}
            >
              <Text color={text.dark.ALPHA} size="lambda" uppercase bold inline>
                Load more
              </Text>
            </Button>
          </Spacing>
        )}
      </Fragment>
    );
  }
}

const mapStateToProps = ({ context, meta, preferences, search }) => ({
  repositories: meta.repositories,
  history: search.history,
  fileResultsLimit: preferences[PREFERENCE_KEYS.CODE_SEARCH_FILE_RESULTS_LIMIT],
  windowWidth: context.window.width,
});

export default compose(
  onlyUpdateForKeys(['results']),
  connect(mapStateToProps),
  withTelemetry,
)(SearchResultsContainer);
