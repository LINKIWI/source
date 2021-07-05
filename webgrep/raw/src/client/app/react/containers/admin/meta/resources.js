import PropTypes from 'prop-types';
import React from 'react';
import { Spacing } from 'react-elemental';
import { connect } from 'react-redux';
import { MdLink } from 'react-icons/md';
import AdminPanel from 'client/app/react/components/admin/panel';
import ExternalLink from 'client/app/react/components/admin/meta/external-link';

/**
 * Resources section, with config-defined external links.
 */
const AdminMetaResourcesContainer = ({ resources }) => (
  <AdminPanel
    title="Resources"
    subtitle="External links"
    icon={<MdLink />}
  >
    {resources.map((resource, idx) => (
      <Spacing key={resource.title} size="small" bottom={idx < resources.length - 1}>
        <ExternalLink title={resource.title} href={resource.href} />
      </Spacing>
    ))}
  </AdminPanel>
);

AdminMetaResourcesContainer.propTypes = {
  resources: PropTypes.array.isRequired,
};

const mapStateToProps = ({ config }) => ({
  resources: config.client.resources || [],
});

export default connect(mapStateToProps)(AdminMetaResourcesContainer);
