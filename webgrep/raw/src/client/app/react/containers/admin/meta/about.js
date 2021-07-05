import PropTypes from 'prop-types';
import React from 'react';
import { spacing, Text } from 'react-elemental';
import { connect } from 'react-redux';
import { MdCode } from 'react-icons/md';
import AdminPanel from 'client/app/react/components/admin/panel';
import { objLookup } from 'shared/util/data';

/**
 * About section text, read from config.
 */
const AdminMetaAboutContainer = ({ description }) => (
  <AdminPanel
    title="About"
    subtitle="Deployment information"
    icon={<MdCode />}
  >
    <Text size="kilo" style={{ lineHeight: spacing.default }}>
      {description}
    </Text>
  </AdminPanel>
);

AdminMetaAboutContainer.propTypes = {
  description: PropTypes.string,
};

AdminMetaAboutContainer.defaultProps = {
  description: null,
};

const mapStateToProps = ({ config }) => ({
  description: objLookup(config, ['client', 'site', 'about']) ||
    'webgrep is a web client for livegrep.',
});

export default connect(mapStateToProps)(AdminMetaAboutContainer);
