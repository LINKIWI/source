import PropTypes from 'prop-types';
import React, { memo, Fragment } from 'react';
import { Text } from 'react-elemental';
import { lifecycle, onlyUpdateForKeys } from 'recompose';
import { compose, withToggleState } from '@linkiwi/hoc';
import { text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

/**
 * One line of source with appropriate highlighting applied based on index bounds.
 */
const HighlightedSource = lifecycle({
  shouldComponentUpdate: false,
})(({ line, bounds }) => {
  if (!bounds) {
    return (
      <Text size="lambda" secondary inline>
        {line}
      </Text>
    );
  }

  const [start, end] = bounds;
  const prefix = line.substring(0, start);
  const highlight = line.substring(start, end);
  const suffix = line.substring(end);

  return (
    <Fragment>
      <Text size="lambda" secondary inline>
        {prefix}
      </Text>
      <Text
        size="lambda"
        style={{ backgroundColor: text.highlight.MATCH }}
        secondary
        inline
        bold
      >
        {highlight}
      </Text>
      <Text size="lambda" secondary inline>
        {suffix}
      </Text>
    </Fragment>
  );
});

/**
 * Describes one line of code, with its line number, content, and external link.
 */
const CodeLine = ({
  number,
  line,
  bounds,
  href,
  onClick,
  isHover,
  handleHoverStart,
  handleHoverEnd,
}) => (
  <a
    rel="noreferrer noopener"
    target="_blank"
    href={href}
    onFocus={handleHoverStart}
    onBlur={handleHoverEnd}
    onMouseEnter={handleHoverStart}
    onMouseLeave={handleHoverEnd}
    onClick={onClick}
    draggable={false}
    style={{
      alignItems: 'center',
      backgroundColor: isHover ? text.highlight.CONTEXT : 'unset',
      display: 'flex',
      height: '18px',
      padding: '1px 2px',
      textDecoration: 'none',
      transition: transition.all.FAST,
      MozUserSelect: 'text',
    }}
  >
    <div style={{ flexShrink: 0, MozUserSelect: 'none', userSelect: 'none', width: '40px' }}>
      <Text color={text.dark.GAMMA} size="lambda" bold={!!bounds} secondary>
        {number}
      </Text>
    </div>

    <div style={{ display: 'inherit', whiteSpace: 'pre' }}>
      <HighlightedSource line={line} bounds={bounds} />
    </div>
  </a>
);

CodeLine.propTypes = {
  number: PropTypes.number.isRequired,
  line: PropTypes.string.isRequired,
  bounds: PropTypes.arrayOf(PropTypes.number.isRequired),
  href: PropTypes.string.isRequired,
  onClick: PropTypes.func,
  // HOC props
  isHover: PropTypes.bool.isRequired,
  handleHoverStart: PropTypes.func.isRequired,
  handleHoverEnd: PropTypes.func.isRequired,
};

CodeLine.defaultProps = {
  bounds: null,
  onClick: () => {},
};

export default compose(
  memo,
  onlyUpdateForKeys(['bounds', 'isHover']),
  withToggleState({
    key: 'isHover',
    enable: 'handleHoverStart',
    disable: 'handleHoverEnd',
  }),
)(CodeLine);
