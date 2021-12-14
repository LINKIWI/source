import PropTypes from 'prop-types';
import React, { createRef, Component, Fragment } from 'react';
import { Modal, Spacing } from 'react-elemental';
import { MdClear, MdContentCopy, MdDownload, MdKeyboardReturn, MdSearch } from 'react-icons/md';
import { connect } from 'react-redux';
import { withResource } from 'supercharged/client';
import { compose } from '@linkiwi/hoc';
import withSelection from 'client/app/react/containers/hoc/with-selection';
import KeyboardListener from 'client/app/react/components/passive/keyboard-listener';
import SourceCodeBlock from 'client/app/react/components/source/source-code-block';
import SourceMetadataHeader from 'client/app/react/components/source/metadata-header';
import ErrorAlert from 'client/app/react/components/ui/error-alert';
import IconButton from 'client/app/react/components/ui/icon-button';
import FloatingButton from 'client/app/react/components/ui/floating-button';
import Spinner from 'client/app/react/components/ui/spinner';
import Tooltip from 'client/app/react/components/ui/tooltip';
import { PREFERENCE_KEYS } from 'client/app/util/constants/preferences';
import { fileSize } from 'client/app/util/format';
import { scroll } from 'client/app/util/navigation';
import { transition } from 'client/app/util/style/transition';
import { CANONICAL_FILE_EXTENSION_LANGUAGES, SYNTAX_HIGHLIGHT_THEME_NAMES } from 'client/resources/data/source';
import { decodeBase64 } from 'shared/util/data';

/**
 * Modal container for fetching and displaying a source file preview.
 */
class SourcePreviewContainer extends Component {
  static propTypes = {
    source: PropTypes.shape({
      err: PropTypes.shape({
        message: PropTypes.string,
      }),
      data: PropTypes.string,
      isLoaded: PropTypes.bool.isRequired,
    }).isRequired,
    repo: PropTypes.string.isRequired,
    version: PropTypes.string.isRequired,
    path: PropTypes.string.isRequired,
    focusLine: PropTypes.number,
    urlPattern: PropTypes.string.isRequired,
    onHide: PropTypes.func.isRequired,
    onSearchQueryChange: PropTypes.func.isRequired,
    onPathCopy: PropTypes.func.isRequired,
    onFileDownload: PropTypes.func.isRequired,
    onSingleFileSearchClick: PropTypes.func.isRequired,
    // HOC props
    width: PropTypes.number,
    themePreference: PropTypes.string.isRequired,
    selection: PropTypes.shape({
      anchor: PropTypes.oneOfType([
        PropTypes.instanceOf(Element),
        PropTypes.instanceOf(window.Text),
      ]),
      text: PropTypes.string.isRequired,
      bounds: PropTypes.object.isRequired,
    }).isRequired,
  };

  static defaultProps = {
    focusLine: null,
    width: null,
  };

  code = createRef();

  handleSearchPromptClick = this._handleSearchPromptClick.bind(this);

  _handleSearchPromptClick() {
    const { selection, onSearchQueryChange, onHide } = this.props;

    if (selection.text) {
      onSearchQueryChange(selection.text);
      onHide();
      // Also scroll to the top of the document, since the search results have changed
      scroll();
    }
  }

  render() {
    const {
      source: { err, data, isLoaded },
      repo,
      version,
      path,
      focusLine,
      urlPattern,
      onHide,
      width,
      themePreference,
      selection,
      onPathCopy,
      onFileDownload,
      onSingleFileSearchClick,
    } = this.props;

    const { data: contents, ok } = decodeBase64(data);
    const language = Object.entries(CANONICAL_FILE_EXTENSION_LANGUAGES)
      .reduce((acc, [extension, syntax]) => (path.endsWith(extension) ? syntax : acc), 'text');
    const annotations = [
      language,
      `${(isLoaded && ok) ? (contents.match(/\n/g) || ['']).length : '-'} lines`,
      (isLoaded && ok) ? fileSize(contents.length) : '- KB',
      version,
    ];
    const download = new File([contents], path.split(/[\\/]/).pop(), { type: 'text/plain' });

    const body = (() => {
      if (err) {
        return (
          <ErrorAlert
            size="beta"
            message={err.message}
          />
        );
      }

      if (!isLoaded) {
        return (
          <Spinner />
        );
      }

      if (!ok) {
        return (
          <ErrorAlert
            size="beta"
            message="Failure decoding base64-encoded file payload; it is likely corrupted."
          />
        );
      }

      const scrollOffset = this.code.current ? this.code.current.scrollLeft : 0;
      const codeBounds = this.code.current ?
        this.code.current.getBoundingClientRect() :
        { x: 0, y: 0 };
      const isSearchPromptVisible = selection.text &&
        selection.text !== '\n' &&
        selection.bounds.height <= 20 &&
        this.code.current &&
        this.code.current.contains(selection.anchor);

      return (
        <Fragment>
          <div
            ref={this.code}
            style={{
              marginTop: '-50px',
              paddingTop: '50px',
              overflowX: 'auto',
              position: 'relative',
            }}
          >
            <SourceCodeBlock
              language={language}
              theme={SYNTAX_HIGHLIGHT_THEME_NAMES[themePreference]}
              focusLine={focusLine}
            >
              {contents}
            </SourceCodeBlock>

            <div
              style={{
                display: 'inline-block',
                left: (selection.bounds.left + selection.bounds.right) / 2 -
                  codeBounds.x +
                  scrollOffset,
                marginTop: '2px',
                opacity: 0,
                position: 'absolute',
                top: selection.bounds.y - codeBounds.y - 25,
                transform: 'translateX(-50%) translateY(-50%)',
                transition: [transition.margin.DEFAULT, transition.opacity.DEFAULT].join(','),
                zIndex: -1,
                ...isSearchPromptVisible && {
                  marginTop: 0,
                  opacity: 1,
                  zIndex: 100,
                },
              }}
            >
              <FloatingButton
                value="New search"
                icon={<MdKeyboardReturn />}
                onClick={this.handleSearchPromptClick}
              />
            </div>
          </div>

          <KeyboardListener
            character="Enter"
            handler={this.handleSearchPromptClick}
          />
        </Fragment>
      );
    })();

    return width && (
      <Modal
        size="alpha"
        onHide={onHide}
        style={{
          borderRadius: '3px',
          boxShadow: '0px 12px 24px -2px rgba(0, 0, 0, 0.1)',
          maxHeight: '100%',
          overflow: 'auto',
          width: width < 1325 ? '100%' : '1325px',
        }}
      >
        <div>
          <Spacing size="large" top right left bottom padding>
            <Spacing bottom>
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  position: 'relative',
                  zIndex: 100,
                }}
              >
                <SourceMetadataHeader
                  repo={repo}
                  version={version}
                  path={path}
                  urlPattern={urlPattern}
                  annotations={annotations}
                />

                <Spacing style={{ display: 'flex' }} left>
                  <Spacing size="tiny" right>
                    <Tooltip
                      description="Copy path to clipboard"
                      style={{ left: 'unset', right: 0, width: '85px' }}
                    >
                      <IconButton icon={MdContentCopy} onClick={onPathCopy} />
                    </Tooltip>
                  </Spacing>

                  <Spacing size="tiny" right>
                    <Tooltip
                      description="Download raw file"
                      style={{ left: 'unset', right: 0, width: '85px' }}
                    >
                      <IconButton icon={MdDownload} onClick={onFileDownload(download)} />
                    </Tooltip>
                  </Spacing>

                  <Spacing size="tiny" right>
                    <Tooltip
                      description="Search within this file"
                      style={{ left: 'unset', right: 0, width: '85px' }}
                    >
                      <IconButton
                        icon={MdSearch}
                        onClick={compose(onHide, onSingleFileSearchClick)}
                      />
                    </Tooltip>
                  </Spacing>

                  <div>
                    <Tooltip
                      description="Close"
                      style={{ left: 'unset', right: 0 }}
                    >
                      <IconButton icon={MdClear} onClick={onHide} />
                    </Tooltip>
                  </div>
                </Spacing>
              </div>
            </Spacing>

            {body}
          </Spacing>
        </div>
      </Modal>
    );
  }
}

const mapStateToProps = ({ context, preferences }) => ({
  width: context.window.width,
  themePreference: preferences[PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME],
});

export default compose(
  connect(mapStateToProps),
  withSelection,
  withResource({
    key: 'source',
    endpoint: '/api/source',
    data: ({ repo, version, path }) => ({ repo, version, path }),
  }),
)(SourcePreviewContainer);
