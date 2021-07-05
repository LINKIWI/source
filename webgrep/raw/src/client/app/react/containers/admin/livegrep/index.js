import React, { Fragment } from 'react';
import { Spacing } from 'react-elemental';
import AdminIndexInfoContainer from 'client/app/react/containers/admin/livegrep/info';
import AdminIndexRepositoriesContainer from 'client/app/react/containers/admin/livegrep/repositories';
import AdminIndexServiceContainer from 'client/app/react/containers/admin/livegrep/service';

/**
 * Admin page "Index" tab.
 */
const AdminIndexContainer = () => (
  <Fragment>
    <Spacing size="large" bottom>
      <AdminIndexInfoContainer />
    </Spacing>

    <Spacing size="large" bottom>
      <AdminIndexRepositoriesContainer />
    </Spacing>

    <Spacing size="large" bottom>
      <AdminIndexServiceContainer />
    </Spacing>
  </Fragment>
);

export default AdminIndexContainer;
