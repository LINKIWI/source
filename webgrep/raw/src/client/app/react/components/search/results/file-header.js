import PropTypes from 'prop-types';
import React, { memo, Fragment } from 'react';
import { colors, spacing, Link, Spacing, Text } from 'react-elemental';
import { MdKeyboardArrowRight } from 'react-icons/md';
import { lifecycle, pure } from 'recompose';
import { text } from 'client/app/util/style/color';

/**
 * Segment of text for a path highlighting the matching index bounds.
 */
const HighlightedPath = lifecycle({
  shouldComponentUpdate: false,
})(({ path, bounds }) => {
  const [start, end] = bounds;
  const prefix = path.substring(0, start);
  const highlight = path.substring(start, end);
  const suffix = path.substring(end);

  return (
    <Fragment>
      <span>
        {prefix}
      </span>
      <span style={{ backgroundColor: text.highlight.MATCH }}>
        {highlight}
      </span>
      <span>
        {suffix}
      </span>
    </Fragment>
  );
});

/**
 * Representation of a file in a repository.
 */
const FileHeader = ({
  repo,
  path,
  repoHref,
  pathHref,
  bounds,
  onRepoClick,
  onPathClick,
  isCollapsible,
}) => (
  <div
    style={{
      alignItems: 'center',
      display: 'flex',
      flexWrap: isCollapsible ? 'wrap' : 'nowrap',
    }}
  >
    <Spacing size="tiny" right>
      <Text
        color="gray60"
        size="kilo"
        style={{ display: 'table-cell', whiteSpace: 'nowrap' }}
        inline
      >
        <Link
          activeColor={colors.primary}
          href={repoHref}
          onClick={onRepoClick}
          style={{ whiteSpace: isCollapsible ? 'normal' : 'unset' }}
        >
          {repo}
        </Link>

        <MdKeyboardArrowRight
          style={{ color: colors.gray30, marginLeft: spacing.tiny, verticalAlign: 'sub' }}
        />
      </Text>
    </Spacing>

    <Text size="kilo" style={{ wordBreak: 'break-word' }} bold>
      <Link activeColor={colors.primary} href={pathHref} onClick={onPathClick}>
        {bounds ? <HighlightedPath path={path} bounds={bounds} /> : path}
      </Link>
    </Text>
  </div>
);

FileHeader.propTypes = {
  repo: PropTypes.string.isRequired,
  path: PropTypes.string.isRequired,
  repoHref: PropTypes.string.isRequired,
  pathHref: PropTypes.string.isRequired,
  bounds: PropTypes.arrayOf(PropTypes.number.isRequired),
  onRepoClick: PropTypes.func,
  onPathClick: PropTypes.func,
  isCollapsible: PropTypes.bool,
};

FileHeader.defaultProps = {
  bounds: null,
  onRepoClick: () => {},
  onPathClick: () => {},
  isCollapsible: false,
};

export default memo(pure(FileHeader));
