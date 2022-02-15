import React, { Fragment } from 'react';
import { Spacing } from 'react-elemental';
import AdminMetaAboutContainer from 'client/app/react/containers/admin/meta/about';
import AdminMetaResourcesContainer from 'client/app/react/containers/admin/meta/resources';
import AdminMetaTelemetryContainer from 'client/app/react/containers/admin/meta/telemetry';

/**
 * Admin page "Meta" tab.
 */
const AdminMetaContainer = () => (
  <Fragment>
    <Spacing size="large" bottom>
      <AdminMetaAboutContainer />
    </Spacing>

    <Spacing size="large" bottom>
      <AdminMetaResourcesContainer />
    </Spacing>

    <Spacing size="large" bottom>
      <AdminMetaTelemetryContainer />
    </Spacing>
  </Fragment>
);

export default AdminMetaContainer;
