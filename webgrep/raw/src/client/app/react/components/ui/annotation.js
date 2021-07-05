import PropTypes from 'prop-types';
import React from 'react';
import { Text } from 'react-elemental';
import { background } from 'client/app/util/style/color';

/**
 * Static, informational annotation tag.
 */
const Annotation = ({ children }) => (
  <button
    style={{
      backgroundColor: background.dark.GAMMA,
      border: 0,
      borderRadius: '3px',
      display: 'inline-flex',
      padding: '5px 10px',
    }}
  >
    <Text
      color="rgba(0, 0, 0, 0.2)"
      size="11px"
      style={{ display: 'inherit', whiteSpace: 'nowrap' }}
      uppercase
      bold
      inline
    >
      {children}
    </Text>
  </button>
);

Annotation.propTypes = {
  children: PropTypes.node.isRequired,
};

export default Annotation;
