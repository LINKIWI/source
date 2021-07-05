import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { compose, withForm, withToggleState } from '@linkiwi/hoc';
import withClipboard from 'client/app/react/containers/hoc/with-clipboard';
import withTelemetry from 'client/app/react/containers/hoc/with-telemetry';
import withToast from 'client/app/react/containers/hoc/with-toast';
import SourcePreviewContainer from 'client/app/react/containers/source/preview';
import KeyboardListener from 'client/app/react/components/passive/keyboard-listener';
import CodeSnippet from 'client/app/react/components/search/results/code-snippet';
import { INPUT_FIELD_IDS } from 'client/app/util/constants/dom';
import { TELEMETRY_ACTIONS } from 'shared/constants/telemetry';
import { PREFERENCE_KEYS, PREFERENCE_VALUES } from 'client/app/util/constants/preferences';
import { scroll } from 'client/app/util/navigation';
import { objLookup } from 'shared/util/data';

/**
 * Container abstracting out stateful interactions with a single code snippet instance.
 */
class CodeSnippetContainer extends Component {
  static propTypes = {
    repositories: PropTypes.object.isRequired,
    snippet: PropTypes.shape({
      repo: PropTypes.string.isRequired,
      version: PropTypes.string.isRequired,
      path: PropTypes.string.isRequired,
      lines: PropTypes.array.isRequired,
    }).isRequired,
    position: PropTypes.number.isRequired,
    permalink: PropTypes.string.isRequired,
    isHighlighted: PropTypes.bool.isRequired,
    onSearchQueryChange: PropTypes.func.isRequired,
    onSearchRepositoryFilePathChange: PropTypes.func.isRequired,
    // HOC props
    form: PropTypes.shape({
      sourcePreviewLine: PropTypes.number,
    }).isRequired,
    handleFormChange: PropTypes.func.isRequired,
    navigationPreference: PropTypes.string.isRequired,
    clipboard: PropTypes.shape({
      read: PropTypes.func.isRequired,
      write: PropTypes.func.isRequired,
    }).isRequired,
    recordTelemetryEvent: PropTypes.func.isRequired,
    toast: PropTypes.func.isRequired,
    isSourcePreviewVisible: PropTypes.bool.isRequired,
    showSourcePreview: PropTypes.func.isRequired,
    hideSourcePreview: PropTypes.func.isRequired,
  };

  handleCodePathCopy = this._handleCodePathCopy.bind(this);

  handleSourcePreviewShow = this._handleSourcePreviewShow.bind(this);

  handleSourcePreviewHide = this._handleSourcePreviewHide.bind(this);

  handleSingleFileSearch = this._handleSingleFileSearch.bind(this);

  handleCodePathClick = this._handleCodePathClick.bind(this);

  handleCodeLineClick = this._handleCodeLineClick.bind(this);

  handleShowSourcePreviewShortcut = this._handleShowSourcePreviewShortcut.bind(this);

  handlePathCopyShortcut = this._handlePathCopyShortcut.bind(this);

  _handleCodePathCopy() {
    const { snippet, clipboard, toast, recordTelemetryEvent } = this.props;

    recordTelemetryEvent(TELEMETRY_ACTIONS.CLIPBOARD_WRITE);
    clipboard.write(snippet.path, (err) => toast(err ?
      'There was an error writing to the system clipboard.' :
      `Copied to clipboard: ${snippet.path}`));
  }

  _handleSourcePreviewShow(line) {
    const { showSourcePreview, recordTelemetryEvent, handleFormChange } = this.props;

    return () => {
      // If a line context is available, request that the source preview focus on that line after
      // the component mounts.
      if (line) {
        handleFormChange('sourcePreviewLine')(line.number);
      }

      recordTelemetryEvent(TELEMETRY_ACTIONS.SOURCE_PREVIEW);
      showSourcePreview();
    };
  }

  _handleSourcePreviewHide() {
    const { hideSourcePreview, handleFormChange } = this.props;

    // Clear the current line; it may be reset on the next source preview launch.
    hideSourcePreview();
    handleFormChange('sourcePreviewLine')({ target: { value: null } });
  }

  _handleSingleFileSearch() {
    const { snippet, onSearchRepositoryFilePathChange, toast } = this.props;

    onSearchRepositoryFilePathChange(snippet.repo, snippet.path);
    toast(`Applied search filters: repository ${snippet.repo} and file path ${snippet.path}.`);
    scroll();
  }

  _handleCodePathClick() {
    const { snippet, position, recordTelemetryEvent } = this.props;

    recordTelemetryEvent(TELEMETRY_ACTIONS.CLICK_CODE_RESULT_PATH, 1, { repo: snippet.repo });
    recordTelemetryEvent(
      TELEMETRY_ACTIONS.CLICK_CODE_RESULT_POSITION,
      position,
      { repo: snippet.repo },
    );
  }

  _handleCodeLineClick(line) {
    const { snippet, position, navigationPreference, recordTelemetryEvent } = this.props;

    return (evt) => {
      // While the anchor element already has the attributes to be self-sufficient as an outgoing
      // link, use an event handler in order to apply preference-dependent logic for the behavior.

      recordTelemetryEvent(TELEMETRY_ACTIONS.CLICK_CODE_RESULT_LINE, 1, { repo: snippet.repo });
      recordTelemetryEvent(
        TELEMETRY_ACTIONS.CLICK_CODE_RESULT_POSITION,
        position,
        { repo: snippet.repo },
      );

      // Let the browser natively handle non-left clicks and modified clicks
      if (evt.button !== 0 || evt.metaKey || evt.altKey || evt.ctrlKey || evt.shiftKey) {
        return null;
      }

      evt.preventDefault();

      switch (navigationPreference) {
        case PREFERENCE_VALUES[PREFERENCE_KEYS.CODE_SEARCH_RESULT_NAVIGATION].SOURCE_LINK_SAME_TAB:
          return window.open(evt.currentTarget.getAttribute('href'), '_self');
        case PREFERENCE_VALUES[PREFERENCE_KEYS.CODE_SEARCH_RESULT_NAVIGATION].SOURCE_PREVIEW:
          return this.handleSourcePreviewShow(line)(evt);
        default:
          // Default behavior: open source link in a new tab
          return window.open(evt.currentTarget.getAttribute('href'), '_blank');
      }
    };
  }

  _handleShowSourcePreviewShortcut(evt) {
    if (!INPUT_FIELD_IDS.includes(evt.target.id)) {
      this.handleSourcePreviewShow()(evt);
    }
  }

  _handlePathCopyShortcut(evt) {
    if (!INPUT_FIELD_IDS.includes(evt.target.id)) {
      this.handleClipboardCopy(this.props.snippet.path)(evt);
    }
  }

  render() {
    const {
      repositories,
      snippet,
      permalink,
      onSearchQueryChange,
      isHighlighted,
      isSourcePreviewVisible,
      form: { sourcePreviewLine },
    } = this.props;

    const numMatches = snippet.lines.reduce((acc, line) =>
      (line.bounds ? acc + 1 : acc), 0);

    return (
      <Fragment>
        <CodeSnippet
          repo={snippet.repo}
          path={snippet.path}
          version={snippet.version}
          lines={snippet.lines}
          urlPattern={objLookup(repositories, [snippet.repo, 'url'])}
          numMatches={numMatches}
          permalink={permalink}
          onCodePathCopy={this.handleCodePathCopy}
          onCodePathClick={this.handleCodePathClick}
          onCodeLineClick={this.handleCodeLineClick}
          onSourcePreviewClick={this.handleSourcePreviewShow()}
          onSingleFileSearchClick={this.handleSingleFileSearch}
        />

        {isSourcePreviewVisible && (
          <SourcePreviewContainer
            repo={snippet.repo}
            version={snippet.version}
            path={snippet.path}
            focusLine={sourcePreviewLine}
            urlPattern={objLookup(repositories, [snippet.repo, 'url'])}
            onHide={this.handleSourcePreviewHide}
            onSearchQueryChange={onSearchQueryChange}
            onSingleFileSearchClick={this.handleSingleFileSearch}
            onPathCopy={this.handleCodePathCopy}
          />
        )}

        {isHighlighted && (
          <KeyboardListener
            character="e"
            handler={this.handleShowSourcePreviewShortcut}
          />
        )}

        {isHighlighted && (
          <KeyboardListener
            character="y"
            handler={this.handlePathCopyShortcut}
          />
        )}
      </Fragment>
    );
  }
}

const mapStateToProps = ({ preferences }) => ({
  navigationPreference: preferences[PREFERENCE_KEYS.CODE_SEARCH_RESULT_NAVIGATION],
});

export default compose(
  connect(mapStateToProps),
  withClipboard,
  withToast,
  withTelemetry,
  withForm({
    initial: () => ({ sourcePreviewLine: null }),
  }),
  withToggleState({
    key: 'isSourcePreviewVisible',
    enable: 'showSourcePreview',
    disable: 'hideSourcePreview',
  }),
)(CodeSnippetContainer);
