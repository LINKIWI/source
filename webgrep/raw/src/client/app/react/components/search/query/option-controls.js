import PropTypes from 'prop-types';
import React, { createRef, Component } from 'react';
import { Spacing } from 'react-elemental';
import { compose, withToggleState } from '@linkiwi/hoc';
import FilePathFilter from 'client/app/react/components/search/query/file-path-filter';
import RepositoryFilter from 'client/app/react/components/search/query/repository-filter';
import Flyout from 'client/app/react/components/ui/flyout';
import Toggle from 'client/app/react/components/ui/toggle';
import { MAX_MATCHES } from 'client/app/util/constants/search';

// Reasonably robust regular expression patterns for capturing various components of a file path
// that may be interesting as an automatically populated file path filter suggestion.
const ROOT_DIRECTORY_PATTERN = /^([^/]+\/)/;
const FILE_EXTENSION_PATTERN = /(\.[^./]+)$/;

/**
 * Controls for query options.
 */
class QueryOptionControls extends Component {
  static propTypes = {
    regex: PropTypes.bool.isRequired,
    caseSensitive: PropTypes.bool.isRequired,
    maxMatches: PropTypes.number.isRequired,
    filePath: PropTypes.string.isRequired,
    repositories: PropTypes.array.isRequired,
    filteredRepos: PropTypes.arrayOf(PropTypes.string.isRequired).isRequired,
    searchResults: PropTypes.arrayOf(PropTypes.shape({
      path: PropTypes.string.isRequired,
    })).isRequired,
    onRegexChange: PropTypes.func.isRequired,
    onCaseSensitivityChange: PropTypes.func.isRequired,
    onMaxMatchesChange: PropTypes.func.isRequired,
    onFilePathChange: PropTypes.func.isRequired,
    onRepositoryToggle: PropTypes.func.isRequired,
    isCompact: PropTypes.bool.isRequired,
    // HOC props
    isPathPatternFlyoutVisible: PropTypes.bool.isRequired,
    showPathPatternFlyout: PropTypes.func.isRequired,
    hidePathPatternFlyout: PropTypes.func.isRequired,
    isRepositoriesFlyoutVisible: PropTypes.bool.isRequired,
    showRepositoriesFlyout: PropTypes.func.isRequired,
    hideRepositoriesFlyout: PropTypes.func.isRequired,
  };

  shouldComponentUpdate(nextProps) {
    const { searchResults, isPathPatternFlyoutVisible } = this.props;

    // Search results are expected to change frequently, since they are updated on every new query.
    // However, search results are only used to render suggestions in the path pattern flyout, so
    // unnecessary re-renders can be avoided by updating the component in response to a different
    // set of search results only when the path pattern flyout is visible.
    // This logic should only run when the search results change, the detection of which is
    // unfortunately a potentially expensive operation if there are many search results. A
    // reasonably high-signal heuristic to capture most (but not all) of the cases in which search
    // results change is examining whether the length of the search results have changed, or, in the
    // case that the length is unchanged, whether the first results have different paths.
    const searchResultsLikelyChanged =
      searchResults.length !== nextProps.searchResults.length ||
      (
        searchResults.length > 0 &&
        nextProps.searchResults > 0 &&
        searchResults[0].path !== nextProps.searchResults[0].path
      );

    if (searchResultsLikelyChanged) {
      if (!isPathPatternFlyoutVisible && !nextProps.isPathPatternFlyoutVisible) {
        return false;
      }
    }

    return true;
  }

  componentDidUpdate(prevProps) {
    // Automatically focus the file path and repository filter fields when the flyout is opened.
    // It's hackily wrapped in a setTimeout to work around a race condition with the toggle stealing
    // focus on click.

    if (
      this.props.isPathPatternFlyoutVisible &&
      !prevProps.isPathPatternFlyoutVisible &&
      this.filePathFieldRef.current
    ) {
      window.setTimeout(() => this.filePathFieldRef.current.focus(), 100);
    }

    if (
      this.props.isRepositoriesFlyoutVisible &&
      !prevProps.isRepositoriesFlyoutVisible &&
      this.repositoriesFieldRef.current
    ) {
      window.setTimeout(() => this.repositoriesFieldRef.current.focus(), 100);
    }
  }

  filePathToggleRef = createRef();

  filePathFieldRef = createRef();

  repositoriesToggleRef = createRef();

  repositoriesFieldRef = createRef();

  render() {
    const {
      regex,
      caseSensitive,
      maxMatches,
      filePath,
      repositories,
      filteredRepos,
      searchResults,
      onRegexChange,
      onCaseSensitivityChange,
      onMaxMatchesChange,
      onFilePathChange,
      onRepositoryToggle,
      isRepositoriesFlyoutVisible,
      showRepositoriesFlyout,
      hideRepositoriesFlyout,
      isPathPatternFlyoutVisible,
      showPathPatternFlyout,
      hidePathPatternFlyout,
      isCompact,
    } = this.props;

    const suggestions = isPathPatternFlyoutVisible && (() => {
      // Data manipulation routine to match each result path against a regular expression, and order
      // the resulting capturing groups by descending order of occurrence frequency.
      const orderByRegexMatchFrequency = (pattern) => {
        const frequency = searchResults
          // Match all paths against the desired pattern, extracting the first capture group as long
          // as it is available.
          .map((result) => result.path.match(pattern))
          .filter((match) => match && match.length)
          .map(([_, capture]) => capture)
          // Transform it into a mapping of file extension to its occurrence frequency.
          .reduce((freq, capture) => ({
            ...freq,
            [capture]: (freq[capture] || 0) + 1,
          }), {});

        return Object.keys(frequency)
          // Rank the extensions in descending order of occurrence frequency.
          .sort((a, b) => frequency[b] - frequency[a]);
      };

      const directories = orderByRegexMatchFrequency(ROOT_DIRECTORY_PATTERN).slice(0, 5);
      const extensions = orderByRegexMatchFrequency(FILE_EXTENSION_PATTERN).slice(0, 5);

      // It's not particularly useful to provide suggestions if all current search results already
      // share the same extension.
      return [
        ...directories.length > 1 ? directories : [],
        ...extensions.length > 1 ? extensions : [],
      ];
    })();

    const containerStyle = {
      display: 'flex',
      ...isCompact && {
        flexWrap: 'wrap',
        marginBottom: '-10px',
      },
    };

    const toggleStyle = {
      marginBottom: 0,
      ...isCompact && {
        marginBottom: '10px',
      },
    };

    const flyoutToggleStyle = {
      position: 'relative',
      ...isCompact && {
        marginBottom: '10px',
        position: 'unset',
      },
    };

    return (
      <div style={containerStyle}>
        <Spacing size="small" style={toggleStyle} right>
          <Toggle
            isActive={regex}
            description="Search using regular expressions"
            onClick={() => onRegexChange(!regex)}
          >
            Regex
          </Toggle>
        </Spacing>

        <Spacing size="small" style={toggleStyle} right>
          <Toggle
            isActive={caseSensitive}
            description="Respect text casing"
            onClick={() => onCaseSensitivityChange(!caseSensitive)}
          >
            Case sensitive
          </Toggle>
        </Spacing>

        <Spacing size="small" style={flyoutToggleStyle} right>
          <Toggle
            ref={this.filePathToggleRef}
            isActive={filePath.length > 0}
            description="File path pattern to filter"
            onClick={isPathPatternFlyoutVisible ? hidePathPatternFlyout : showPathPatternFlyout}
          >
            File path
          </Toggle>

          <Flyout
            triggerRef={this.filePathToggleRef}
            isOpen={isPathPatternFlyoutVisible}
            isCompact={isCompact}
            onHide={hidePathPatternFlyout}
          >
            <FilePathFilter
              ref={this.filePathFieldRef}
              value={filePath}
              suggestions={suggestions || []}
              onChange={onFilePathChange}
              onHide={hidePathPatternFlyout}
            />
          </Flyout>
        </Spacing>

        <Spacing size="small" style={flyoutToggleStyle} right>
          <Toggle
            ref={this.repositoriesToggleRef}
            isActive={filteredRepos.length > 0}
            description={filteredRepos.length ?
              [
                'Including',
                filteredRepos.length,
                filteredRepos.length === 1 ? 'repository' : 'repositories',
              ].join(' ') :
              'Select repositories to filter'}
            onClick={isRepositoriesFlyoutVisible ? hideRepositoriesFlyout : showRepositoriesFlyout}
          >
            Repositories
          </Toggle>

          <Flyout
            triggerRef={this.repositoriesToggleRef}
            isOpen={isRepositoriesFlyoutVisible}
            isCompact={isCompact}
            onHide={hideRepositoriesFlyout}
          >
            <RepositoryFilter
              ref={this.repositoriesFieldRef}
              repositories={repositories}
              selectedRepos={new Set(filteredRepos)}
              onRepositoryToggle={onRepositoryToggle}
              onHide={hideRepositoriesFlyout}
            />
          </Flyout>
        </Spacing>

        <div style={toggleStyle}>
          <Toggle
            isActive={maxMatches === MAX_MATCHES.UNLIMITED}
            description="Exhaust all results (slow)"
            onClick={() => onMaxMatchesChange(maxMatches === MAX_MATCHES.UNLIMITED ?
              null : MAX_MATCHES.UNLIMITED)}
          >
            All matches
          </Toggle>
        </div>
      </div>
    );
  }
}

export default compose(
  withToggleState({
    key: 'isPathPatternFlyoutVisible',
    enable: 'showPathPatternFlyout',
    disable: 'hidePathPatternFlyout',
  }),
  withToggleState({
    key: 'isRepositoriesFlyoutVisible',
    enable: 'showRepositoriesFlyout',
    disable: 'hideRepositoriesFlyout',
  }),
)(QueryOptionControls);
