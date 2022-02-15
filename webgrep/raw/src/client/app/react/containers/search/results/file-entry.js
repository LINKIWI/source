import PropTypes from 'prop-types';
import React, { createRef, Component, Fragment } from 'react';
import { Spacing } from 'react-elemental';
import { MdCode, MdContentCopy, MdSearch } from 'react-icons/md';
import { compose, withToggleState } from '@linkiwi/hoc';
import withClipboard from 'client/app/react/containers/hoc/with-clipboard';
import withTelemetry from 'client/app/react/containers/hoc/with-telemetry';
import withToast from 'client/app/react/containers/hoc/with-toast';
import SourcePreviewContainer from 'client/app/react/containers/source/preview';
import FileHeader from 'client/app/react/components/search/results/file-header';
import OverflowMenu from 'client/app/react/components/ui/overflow-menu';
import ButtonMenuItem from 'client/app/react/components/ui/overflow-menu/button-menu-item';
import { string } from 'client/app/util/format';
import { scroll } from 'client/app/util/navigation';
import { transition } from 'client/app/util/style/transition';
import { TELEMETRY_ACTIONS } from 'shared/constants/telemetry';
import { objLookup } from 'shared/util/data';

class FileEntry extends Component {
  static propTypes = {
    file: PropTypes.shape({
      repo: PropTypes.string.isRequired,
      version: PropTypes.string.isRequired,
      path: PropTypes.string.isRequired,
      bounds: PropTypes.array.isRequired,
    }).isRequired,
    repositories: PropTypes.object.isRequired,
    position: PropTypes.number.isRequired,
    windowWidth: PropTypes.number.isRequired,
    onSearchQueryChange: PropTypes.func.isRequired,
    onSearchRepositoryFilePathChange: PropTypes.func.isRequired,
    // HOC props
    recordTelemetryEvent: PropTypes.func.isRequired,
    clipboard: PropTypes.shape({
      read: PropTypes.func.isRequired,
      write: PropTypes.func.isRequired,
    }).isRequired,
    toast: PropTypes.func.isRequired,
    isOverflowMenuVisible: PropTypes.bool.isRequired,
    showOverflowMenu: PropTypes.func.isRequired,
    hideOverflowMenu: PropTypes.func.isRequired,
    isSourcePreviewVisible: PropTypes.bool.isRequired,
    showSourcePreview: PropTypes.func.isRequired,
    hideSourcePreview: PropTypes.func.isRequired,
    isHover: PropTypes.bool.isRequired,
    handleMouseEnter: PropTypes.func.isRequired,
    handleMouseLeave: PropTypes.func.isRequired,
  };

  overflow = createRef();

  handleFilePathClick = this._handleFilePathClick.bind(this);

  handleFilePathCopy = this._handleFilePathCopy.bind(this);

  handleFileDownload = this._handleFileDownload.bind(this);

  handleSourcePreview = this._handleSourcePreview.bind(this);

  handleSingleFileSearch = this._handleSingleFileSearch.bind(this);

  _handleFilePathClick() {
    const { file, position, recordTelemetryEvent } = this.props;

    recordTelemetryEvent(TELEMETRY_ACTIONS.CLICK_FILE_RESULT_PATH, 1, { repo: file.repo });
    recordTelemetryEvent(
      TELEMETRY_ACTIONS.CLICK_FILE_RESULT_POSITION,
      position,
      { repo: file.repo },
    );
  }

  _handleFilePathCopy() {
    const { file, clipboard, toast, recordTelemetryEvent } = this.props;

    recordTelemetryEvent(TELEMETRY_ACTIONS.CLIPBOARD_WRITE);
    clipboard.write(file.path, (err) => toast(err ?
      'There was an error writing to the system clipboard.' :
      `Copied to clipboard: ${file.path}`));
  }

  _handleFileDownload(file) {
    const { file: { repo }, recordTelemetryEvent } = this.props;

    return () => {
      recordTelemetryEvent(TELEMETRY_ACTIONS.SOURCE_RAW_DOWNLOAD, 1, { repo });

      const downloadURL = URL.createObjectURL(file);
      window.open(downloadURL, '_blank');
      URL.revokeObjectURL(downloadURL);
    };
  }

  _handleSourcePreview() {
    const { file, recordTelemetryEvent, showSourcePreview } = this.props;

    recordTelemetryEvent(TELEMETRY_ACTIONS.SOURCE_PREVIEW, 1, { repo: file.repo });
    showSourcePreview();
  }

  _handleSingleFileSearch() {
    const { file, toast, onSearchRepositoryFilePathChange } = this.props;

    onSearchRepositoryFilePathChange(file.repo, file.path);
    toast(`Applied search filters: repository ${file.repo} and file path ${file.path}.`);
    scroll();
  }

  render() {
    const {
      file,
      repositories,
      windowWidth,
      isOverflowMenuVisible,
      showOverflowMenu,
      hideOverflowMenu,
      isSourcePreviewVisible,
      hideSourcePreview,
      isHover,
      handleMouseEnter,
      handleMouseLeave,
      onSearchQueryChange,
    } = this.props;

    const isOverflowMenuCompact = this.overflow.current && (() => {
      const { x: overflowMenuX } = this.overflow.current.getBoundingClientRect();
      return windowWidth - overflowMenuX < 170;
    })();

    const containerStyle = {
      alignItems: 'center',
      display: 'inline-flex',
    };

    const overflowStyle = {
      opacity: 0,
      transition: transition.all.DEFAULT,
      ...(isHover || isOverflowMenuVisible) && {
        opacity: 1,
      },
    };

    const overflowMenuStyle = {
      // OverflowMenu sets and unsets position `left` and `right` properties to switch between
      // mirrored and standard modes. Unfortunately, this makes transitions between the two states
      // impossible. Since the width of the overflow menu is known in advance, the following logic
      // overrides the overflow menu positioning in a way that is compatible with transitions.
      left: 'unset',
      right: 0,
      width: '155px',
      marginRight: '-128px',
      ...isOverflowMenuCompact && {
        marginRight: 0,
      },
    };

    return (
      <Fragment>
        <div
          onMouseEnter={handleMouseEnter}
          onMouseLeave={handleMouseLeave}
          style={containerStyle}
        >
          <div>
            <FileHeader
              repo={file.repo}
              path={file.path}
              bounds={file.bounds}
              repoHref={string(objLookup(repositories, [file.repo, 'url']), {
                name: file.repo,
                version: file.version.replace(/.*\//g, ''),
                path: '',
                lno: '',
              })}
              pathHref={string(objLookup(repositories, [file.repo, 'url']), {
                name: file.repo,
                version: file.version.replace(/.*\//g, ''),
                path: file.path,
                lno: 1,
              })}
              onPathClick={this.handleFilePathClick}
            />
          </div>

          <Spacing size="small" style={overflowStyle} left>
            <div ref={this.overflow}>
              <OverflowMenu
                isMenuVisible={isOverflowMenuVisible}
                showMenu={showOverflowMenu}
                hideMenu={hideOverflowMenu}
                buttonStyle={{ height: '18px' }}
                menuStyle={overflowMenuStyle}
                mirrored={!isOverflowMenuCompact}
              >
                <ButtonMenuItem
                  label="Copy file path"
                  icon={<MdContentCopy style={{ display: 'inherit' }} />}
                  onClick={this.handleFilePathCopy}
                />
                <ButtonMenuItem
                  label="Preview source"
                  icon={<MdCode style={{ display: 'inherit' }} />}
                  onClick={this.handleSourcePreview}
                />
                <ButtonMenuItem
                  label="Search this file"
                  icon={<MdSearch style={{ display: 'inherit' }} />}
                  onClick={this.handleSingleFileSearch}
                />
              </OverflowMenu>
            </div>
          </Spacing>
        </div>

        {isSourcePreviewVisible && (
          <SourcePreviewContainer
            repo={file.repo}
            version={file.version}
            path={file.path}
            urlPattern={objLookup(repositories, [file.repo, 'url'])}
            onHide={hideSourcePreview}
            onSearchQueryChange={onSearchQueryChange}
            onSingleFileSearchClick={this.handleSingleFileSearch}
            onPathCopy={this.handleFilePathCopy}
            onFileDownload={this.handleFileDownload}
          />
        )}
      </Fragment>
    );
  }
}

export default compose(
  withClipboard,
  withTelemetry,
  withToast,
  withToggleState({
    key: 'isOverflowMenuVisible',
    enable: 'showOverflowMenu',
    disable: 'hideOverflowMenu',
  }),
  withToggleState({
    key: 'isSourcePreviewVisible',
    enable: 'showSourcePreview',
    disable: 'hideSourcePreview',
  }),
  withToggleState({
    key: 'isHover',
    enable: 'handleMouseEnter',
    disable: 'handleMouseLeave',
  }),
)(FileEntry);
