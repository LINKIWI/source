import PropTypes from 'prop-types';
import React, { Fragment } from 'react';
import { Spacing } from 'react-elemental';
import { MdContentCopy } from 'react-icons/md';
import AdminControl from 'client/app/react/components/admin/control';
import OverflowMenu from 'client/app/react/components/ui/overflow-menu';
import ButtonMenuItem from 'client/app/react/components/ui/overflow-menu/button-menu-item';
import { string } from 'client/app/util/format';
import { background } from 'client/app/util/style/color';

/**
 * Informational representation of an indexed repository and its config.
 */
const Repository = ({
  name,
  version,
  url,
  remote,
  onCopyVersion,
  onCopyRemote,
}) => (
  <Fragment>
    <Spacing size="micro" style={{ display: 'flex' }} bottom>
      <Spacing right>
        <Spacing size="micro" top>
          <div
            style={{
              backgroundColor: background.primary.BETA,
              borderRadius: '50%',
              height: '8px',
              width: '8px',
            }}
          />
        </Spacing>
      </Spacing>

      <Spacing style={{ flexGrow: 1 }} right>
        <AdminControl
          title={name}
          description={remote}
          linkTitle="View repository"
          linkHref={string(url, {
            name,
            version: version.replace(/.*\//g, ''),
            path: '',
            lno: '',
          })}
        />
      </Spacing>

      <OverflowMenu menuStyle={{ width: '150px' }}>
        <ButtonMenuItem
          label="Copy version"
          icon={<MdContentCopy style={{ display: 'inherit' }} />}
          onClick={onCopyVersion}
        />
        <ButtonMenuItem
          label="Copy remote"
          icon={<MdContentCopy style={{ display: 'inherit' }} />}
          onClick={onCopyRemote}
        />
      </OverflowMenu>
    </Spacing>
  </Fragment>
);

Repository.propTypes = {
  name: PropTypes.string.isRequired,
  version: PropTypes.string.isRequired,
  url: PropTypes.string.isRequired,
  remote: PropTypes.string.isRequired,
  onCopyVersion: PropTypes.func.isRequired,
  onCopyRemote: PropTypes.func.isRequired,
};

export default Repository;
