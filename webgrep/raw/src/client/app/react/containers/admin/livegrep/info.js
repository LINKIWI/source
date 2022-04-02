import PropTypes from 'prop-types';
import React, { Fragment } from 'react';
import { Spacing } from 'react-elemental';
import { MdInfo } from 'react-icons/md';
import { connect } from 'react-redux';
import { withResource } from 'supercharged/client';
import { compose } from '@linkiwi/hoc';
import AdminPanel from 'client/app/react/components/admin/panel';
import MetadataItem from 'client/app/react/components/admin/livegrep/metadata-item';
import ErrorAlert from 'client/app/react/components/ui/error-alert';
import Spinner from 'client/app/react/components/ui/spinner';
import { absoluteTimestamp, relativeTimestamp } from 'client/app/util/format';

// Window widths at which to consider the layout to be horizontally compact.
const breakpoints = [
  { threshold: 550, compact: true },
];

/**
 * Index key-value metadata section.
 */
const AdminIndexInfoContainer = ({ info: { err, data }, width }) => {
  const isCompact = breakpoints.reduce((acc, breakpoint) =>
    (width < breakpoint.threshold ? breakpoint.compact : acc), false);

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
              isCompact={isCompact}
            />
          </Spacing>

          <Spacing size="small" bottom>
            <MetadataItem
              heading="Timestamp"
              value={`${absoluteTimestamp(data.timestamp)} (${relativeTimestamp(data.timestamp)})`}
              isCompact={isCompact}
            />
          </Spacing>

          <Spacing size="small" bottom>
            <MetadataItem
              heading="Repositories"
              value={data.repositories.length.toString()}
              isCompact={isCompact}
            />
          </Spacing>

          <MetadataItem
            heading="Version"
            value={process.env.VERSION}
            isCompact={isCompact}
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
  width: PropTypes.number,
};

AdminIndexInfoContainer.defaultProps = {
  width: null,
};

const mapStateToProps = ({ context }) => ({
  width: context.window.width,
});

export default compose(
  connect(mapStateToProps),
  withResource({
    key: 'info',
    cacheKey: () => 'admin:meta',
    endpoint: '/api/meta/info',
  }),
)(AdminIndexInfoContainer);
