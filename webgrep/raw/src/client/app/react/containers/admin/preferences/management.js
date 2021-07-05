import PropTypes from 'prop-types';
import React from 'react';
import { Button } from 'react-elemental';
import { MdSettings } from 'react-icons/md';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import AdminControl from 'client/app/react/components/admin/control';
import AdminPanel from 'client/app/react/components/admin/panel';
import { clearAllPreferences } from 'client/app/redux/actions/preferences';

const AdminPreferencesManagementContainer = ({ actions }) => (
  <AdminPanel
    title="Management"
    subtitle="Manage persisted preferences"
    icon={<MdSettings />}
  >
    <AdminControl
      title="Reset preferences"
      description="Reset all preferences to defaults"
    >
      <Button
        size="gamma"
        text="Reset"
        onClick={actions.clearAllPreferences}
        style={{ borderRadius: '3px', height: '35px', width: '70px' }}
      />
    </AdminControl>
  </AdminPanel>
);

AdminPreferencesManagementContainer.propTypes = {
  // HOC props
  actions: PropTypes.shape({
    clearAllPreferences: PropTypes.func.isRequired,
  }).isRequired,
};

const mapDispatchToProps = (dispatch) => ({
  actions: bindActionCreators({ clearAllPreferences }, dispatch),
});

export default connect(null, mapDispatchToProps)(AdminPreferencesManagementContainer);
