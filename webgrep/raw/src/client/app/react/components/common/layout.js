import PropTypes from 'prop-types';
import React from 'react';
import { pure } from 'recompose';
import { Spacing } from 'react-elemental';
import { transition } from 'client/app/util/style/transition';

// Horizontal padding values at various window width threshold breakpoints.
const breakpoints = [
  { threshold: 1500, padding: 100 },
  { threshold: 900, padding: 50 },
  { threshold: 450, padding: 25 },
];

/**
 * Common abstraction for sizing the width of arbitrary contents based on the window width.
 */
const Layout = ({ containerStyle, contentStyle, children, width }) => {
  if (!width) {
    return null;
  }

  const padding = breakpoints.reduce((acc, breakpoint) =>
    (width < breakpoint.threshold ? breakpoint.padding : acc), 150);

  return (
    <div style={containerStyle}>
      <Spacing
        size={`${padding}px`}
        style={{
          boxSizing: 'border-box',
          margin: 'auto',
          maxWidth: '1800px',
          transition: transition.all.DEFAULT,
          ...contentStyle,
        }}
        left
        right
        padding
      >
        {children}
      </Spacing>
    </div>
  );
};

Layout.propTypes = {
  containerStyle: PropTypes.object,
  contentStyle: PropTypes.object,
  children: PropTypes.node.isRequired,
  width: PropTypes.number,
};

Layout.defaultProps = {
  containerStyle: {},
  contentStyle: {},
  width: null,
};

export default pure(Layout);
