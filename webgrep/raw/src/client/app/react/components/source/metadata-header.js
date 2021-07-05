import PropTypes from 'prop-types';
import React from 'react';
import { Spacing } from 'react-elemental';
import FileHeader from 'client/app/react/components/search/results/file-header';
import Annotation from 'client/app/react/components/ui/annotation';
import { string } from 'client/app/util/format';

/**
 * File metadata header in the source preview modal.
 */
const SourceMetadataHeader = ({ repo, version, path, urlPattern, annotations }) => (
  <div>
    <Spacing size="small" bottom>
      <FileHeader
        repo={repo}
        path={path}
        repoHref={string(urlPattern, {
          name: repo,
          version: version.replace(/.*\//g, ''),
          path: '',
          lno: '',
        })}
        pathHref={string(urlPattern, {
          name: repo,
          version: version.replace(/.*\//g, ''),
          path,
          lno: 1,
        })}
        isCollapsible
      />
    </Spacing>

    <div style={{ alignItems: 'center', display: 'flex', flexWrap: 'wrap' }}>
      {annotations.map((annotation, idx) => (
        <Spacing key={annotation} size="small" right={idx < annotations.length - 1}>
          <Spacing size="8px" bottom>
            <Annotation>
              {annotation}
            </Annotation>
          </Spacing>
        </Spacing>
      ))}
    </div>
  </div>
);

SourceMetadataHeader.propTypes = {
  repo: PropTypes.string.isRequired,
  version: PropTypes.string.isRequired,
  path: PropTypes.string.isRequired,
  urlPattern: PropTypes.string.isRequired,
  annotations: PropTypes.arrayOf(PropTypes.string.isRequired).isRequired,
};

export default SourceMetadataHeader;
