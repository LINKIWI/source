import PropTypes from 'prop-types';
import React, { Fragment } from 'react';
import { Spacing } from 'react-elemental';
import { connect } from 'react-redux';
import AdminPreferencesInterfaceContainer from 'client/app/react/containers/admin/preferences/interface';
import AdminPreferencesManagementContainer from 'client/app/react/containers/admin/preferences/management';
import AdminPreferencesSearchContainer from 'client/app/react/containers/admin/preferences/search';

// Window widths at which to consider the layout to be horizontally compact.
const breakpoints = [
  { threshold: 500, compact: true },
];

/**
 * Admin page "Preferences" tab.
 */
const AdminPreferencesContainer = ({ width }) => {
  const isCompact = breakpoints.reduce((acc, breakpoint) =>
    (width < breakpoint.threshold ? breakpoint.compact : acc), false);

  return (
    <Fragment>
      <Spacing size="large" bottom>
        <AdminPreferencesInterfaceContainer isCompact={isCompact} />
      </Spacing>

      <Spacing size="large" bottom>
        <AdminPreferencesSearchContainer isCompact={isCompact} />
      </Spacing>

      <Spacing size="large" bottom>
        <AdminPreferencesManagementContainer isCompact={isCompact} />
      </Spacing>
    </Fragment>
  );
};

AdminPreferencesContainer.propTypes = {
  // HOC props
  width: PropTypes.number,
};

AdminPreferencesContainer.defaultProps = {
  width: null,
};

const mapStateToProps = ({ context }) => ({
  width: context.window.width,
});

export default connect(mapStateToProps)(AdminPreferencesContainer);
