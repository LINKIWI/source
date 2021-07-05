import PropTypes from 'prop-types';
import React, { Fragment } from 'react';
import { Spacing } from 'react-elemental';
import { MdInfo } from 'react-icons/md';
import { withResource } from 'supercharged/client';
import AdminPanel from 'client/app/react/components/admin/panel';
import MetadataItem from 'client/app/react/components/admin/livegrep/metadata-item';
import ErrorAlert from 'client/app/react/components/ui/error-alert';
import Spinner from 'client/app/react/components/ui/spinner';
import { absoluteTimestamp, relativeTimestamp } from 'client/app/util/format';

/**
 * Index key-value metadata section.
 */
const AdminIndexInfoContainer = ({ info: { err, data } }) => {
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
      return (
        <Fragment>
          <Spacing size="small" bottom>
            <MetadataItem
              heading="Name"
              value={data.name}
            />
          </Spacing>

          <Spacing size="small" bottom>
            <MetadataItem
              heading="Timestamp"
              value={`${absoluteTimestamp(data.timestamp)} (${relativeTimestamp(data.timestamp)})`}
            />
          </Spacing>

          <Spacing size="small" bottom>
            <MetadataItem
              heading="Repositories"
              value={data.repositories.length.toString()}
            />
          </Spacing>

          <MetadataItem
            heading="Version"
            value={process.env.VERSION}
          />
        </Fragment>
      );
    }

    return (
      <Spinner />
    );
  })();

  return (
    <AdminPanel
      title="Information"
      subtitle="Livegrep index metadata"
      icon={<MdInfo />}
    >
      {body}
    </AdminPanel>
  );
};

AdminIndexInfoContainer.propTypes = {
  // HOC props
  info: PropTypes.shape({
    err: PropTypes.object,
    data: PropTypes.object,
  }).isRequired,
};

export default withResource({
  key: 'info',
  cacheKey: 'admin:meta',
  endpoint: '/api/meta/info',
})(AdminIndexInfoContainer);
