import PropTypes from 'prop-types';
import React, { Fragment } from 'react';
import { colors, Spacing } from 'react-elemental';
import { MdCode, MdContentCopy, MdLink, MdSearch } from 'react-icons/md';
import Box from 'client/app/react/components/common/box';
import CodeLine from 'client/app/react/components/search/results/code-line';
import FileHeader from 'client/app/react/components/search/results/file-header';
import Annotation from 'client/app/react/components/ui/annotation';
import OverflowMenu from 'client/app/react/components/ui/overflow-menu';
import AnchorMenuItem from 'client/app/react/components/ui/overflow-menu/anchor-menu-item';
import ButtonMenuItem from 'client/app/react/components/ui/overflow-menu/button-menu-item';
import { string } from 'client/app/util/format';
import { background } from 'client/app/util/style/color';

/**
 * Represents a single code snippet (for one file) in code search results.
 */
const CodeSnippet = ({
  repo,
  path,
  version,
  urlPattern,
  lines,
  numMatches,
  permalink,
  onCodePathCopy,
  onCodePathClick,
  onCodeLineClick,
  onSourcePreviewClick,
  onSingleFileSearchClick,
}) => (
  <Box>
    <Spacing
      style={{ alignItems: 'center', display: 'flex', justifyContent: 'space-between' }}
      bottom
    >
      <Spacing style={{ alignItems: 'center', display: 'flex' }} right>
        <div
          style={{
            backgroundColor: background.primary.BETA,
            borderRadius: '50%',
            flexShrink: 0,
            height: '8px',
            marginLeft: '3px',
            marginRight: '30px',
            width: '8px',
          }}
        />

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
          onPathClick={onCodePathClick}
          isCollapsible
        />
      </Spacing>

      <div style={{ alignItems: 'center', display: 'flex', flexShrink: 0 }}>
        <Spacing size="small" right>
          <Annotation>
            {numMatches} {numMatches === 1 ? 'match' : 'matches'}
          </Annotation>
        </Spacing>

        <OverflowMenu menuStyle={{ width: '155px' }}>
          <ButtonMenuItem
            label="Copy file path"
            icon={<MdContentCopy style={{ display: 'inherit' }} />}
            onClick={onCodePathCopy}
          />
          <ButtonMenuItem
            label="Preview source"
            icon={<MdCode style={{ display: 'inherit' }} />}
            onClick={onSourcePreviewClick}
          />
          <ButtonMenuItem
            label="Search this file"
            icon={<MdSearch style={{ display: 'inherit' }} />}
            onClick={onSingleFileSearchClick}
          />
          <AnchorMenuItem
            label="Permalink"
            icon={<MdLink style={{ display: 'inherit' }} />}
            href={permalink}
          />
        </OverflowMenu>
      </div>
    </Spacing>

    <div style={{ overflow: 'auto' }}>
      {lines.map((line, lineIdx, arr) => (
        <Fragment key={line.number}>
          {lineIdx > 0 && line.number - arr[lineIdx - 1].number > 1 && (
            <div
              style={{
                backgroundColor: colors.gray15,
                height: '0.5px',
                width: '100%',
              }}
            />
          )}

          <CodeLine
            number={line.number}
            line={line.line}
            bounds={line.bounds}
            onClick={onCodeLineClick(line)}
            href={string(urlPattern, {
              name: repo,
              version: version.replace(/.*\//g, ''),
              path,
              lno: line.number,
            })}
          />
        </Fragment>
      ))}
    </div>
  </Box>
);

CodeSnippet.propTypes = {
  repo: PropTypes.string.isRequired,
  path: PropTypes.string.isRequired,
  version: PropTypes.string.isRequired,
  urlPattern: PropTypes.string.isRequired,
  lines: PropTypes.arrayOf(PropTypes.shape({
    number: PropTypes.number.isRequired,
    line: PropTypes.string.isRequired,
    bounds: PropTypes.array,
  })).isRequired,
  numMatches: PropTypes.number.isRequired,
  permalink: PropTypes.string.isRequired,
  onCodePathCopy: PropTypes.func.isRequired,
  onCodePathClick: PropTypes.func.isRequired,
  onCodeLineClick: PropTypes.func.isRequired,
  onSourcePreviewClick: PropTypes.func.isRequired,
  onSingleFileSearchClick: PropTypes.func.isRequired,
};

export default CodeSnippet;
