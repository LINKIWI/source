import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { Button, Spinner } from 'react-elemental';
import { MdRefresh } from 'react-icons/md';
import { withResource } from 'supercharged/client';
import { compose } from '@linkiwi/hoc';
import withToast from 'client/app/react/containers/hoc/with-toast';
import AdminControl from 'client/app/react/components/admin/control';
import AdminPanel from 'client/app/react/components/admin/panel';

/**
 * Imperative controls for changing livegrep process state.
 */
class AdminIndexServiceContainer extends Component {
  handleReload = this._handleReload.bind(this);

  _handleReload() {
    const { reload, toast } = this.props;

    reload.invoke(null, (err) => (err ?
      toast(`Error reloading index: ${err.message}`) :
      toast('Successfully reloaded index.')));
  }

  render() {
    const { reload: { isLoaded } } = this.props;

    return (
      <AdminPanel
        title="Service"
        subtitle="Manage the process serving the index"
        icon={<MdRefresh />}
      >
        <AdminControl
          title="Reload"
          description="Manually reload the contents of the index file on disk."
        >
          <Button
            size="gamma"
            onClick={this.handleReload}
            disabled={!isLoaded}
            style={{
              alignItems: 'center',
              borderRadius: '3px',
              display: 'flex',
              height: '35px',
              justifyContent: 'center',
              width: '80px',
            }}
            {...isLoaded && { text: 'Reload' }}
          >
            {!isLoaded && (
              <Spinner
                size="gamma"
                ringColor="rgba(255, 255, 255, 0.1)"
                accentColor="rgba(255, 255, 255, 0.6)"
              />
            )}
          </Button>
        </AdminControl>
      </AdminPanel>
    );
  }
}

AdminIndexServiceContainer.propTypes = {
  // HOC props
  toast: PropTypes.func.isRequired,
  reload: PropTypes.shape({
    err: PropTypes.object,
    data: PropTypes.object,
    isLoaded: PropTypes.bool.isRequired,
    invoke: PropTypes.func.isRequired,
  }).isRequired,
};

export default compose(
  withToast,
  withResource({
    key: 'reload',
    endpoint: '/api/admin/reload',
    method: 'PUT',
    invokeOnMount: false,
  }),
)(AdminIndexServiceContainer);
