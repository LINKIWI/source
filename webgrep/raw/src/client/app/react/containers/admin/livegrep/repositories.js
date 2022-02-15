import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { Spacing } from 'react-elemental';
import { MdLibraryBooks } from 'react-icons/md';
import { withResource } from 'supercharged/client';
import { compose } from '@linkiwi/hoc';
import withClipboard from 'client/app/react/containers/hoc/with-clipboard';
import withTelemetry from 'client/app/react/containers/hoc/with-telemetry';
import withToast from 'client/app/react/containers/hoc/with-toast';
import Repository from 'client/app/react/components/admin/livegrep/repository';
import AdminPanel from 'client/app/react/components/admin/panel';
import ErrorAlert from 'client/app/react/components/ui/error-alert';
import Spinner from 'client/app/react/components/ui/spinner';
import { TELEMETRY_ACTIONS } from 'shared/constants/telemetry';

/**
 * Listing of all indexed repositories.
 */
class AdminIndexRepositoriesContainer extends Component {
  static propTypes = {
    // HOC props
    info: PropTypes.shape({
      err: PropTypes.object,
      data: PropTypes.object,
    }).isRequired,
    clipboard: PropTypes.shape({
      read: PropTypes.func.isRequired,
      write: PropTypes.func.isRequired,
    }).isRequired,
    recordTelemetryEvent: PropTypes.func.isRequired,
    toast: PropTypes.func.isRequired,
  };

  handleClipboardCopy = this._handleClipboardCopy.bind(this);

  _handleClipboardCopy(text) {
    const { clipboard, toast, recordTelemetryEvent } = this.props;

    return () => {
      recordTelemetryEvent(TELEMETRY_ACTIONS.CLIPBOARD_WRITE);
      clipboard.write(text, (err) => toast(err ?
        'There was an error writing to the system clipboard.' :
        `Copied to clipboard: ${text}`));
    };
  }

  render() {
    const { info: { err, data } } = this.props;

    const body = (() => {
      if (err) {
        return (
          <ErrorAlert
            size="beta"
            message={err.message}
          />
        );
      }

      if (data) {
        return data.repositories.map((repo, idx) => (
          <Spacing key={repo.name} bottom={idx < data.repositories.length - 1}>
            <Repository
              name={repo.name}
              version={repo.version}
              url={repo.url}
              remote={repo.remote}
              onCopyName={this.handleClipboardCopy(repo.name)}
              onCopyVersion={this.handleClipboardCopy(repo.version)}
              onCopyRemote={this.handleClipboardCopy(repo.remote)}
            />
          </Spacing>
        ));
      }

      return (
        <Spinner />
      );
    })();

    return (
      <AdminPanel
        title="Repositories"
        subtitle="Repositories served by the index"
        icon={<MdLibraryBooks />}
      >
        {body}
      </AdminPanel>
    );
  }
}

export default compose(
  withClipboard,
  withTelemetry,
  withToast,
  withResource({
    key: 'info',
    cacheKey: 'admin:meta',
    endpoint: '/api/meta/info',
  }),
)(AdminIndexRepositoriesContainer);
