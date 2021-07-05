import React, { Fragment } from 'react';
import { Spacing } from 'react-elemental';
import AdminPreferencesInterfaceContainer from 'client/app/react/containers/admin/preferences/interface';
import AdminPreferencesManagementContainer from 'client/app/react/containers/admin/preferences/management';
import AdminPreferencesSearchContainer from 'client/app/react/containers/admin/preferences/search';

/**
 * Admin page "Preferences" tab.
 */
const AdminPreferencesContainer = () => (
  <Fragment>
    <Spacing size="large" bottom>
      <AdminPreferencesInterfaceContainer />
    </Spacing>

    <Spacing size="large" bottom>
      <AdminPreferencesSearchContainer />
    </Spacing>

    <Spacing size="large" bottom>
      <AdminPreferencesManagementContainer />
    </Spacing>
  </Fragment>
);

export default AdminPreferencesContainer;
