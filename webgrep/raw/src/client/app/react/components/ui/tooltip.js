import PropTypes from 'prop-types';
import React from 'react';
import { Text, Tooltip as ElementalTooltip } from 'react-elemental';
import { text } from 'client/app/util/style/color';

/**
 * Abstraction over an Elemental Tooltip with some predefined structure and styling.
 */
const Tooltip = ({ description, style, children, ...props }) => (
  <div>
    <ElementalTooltip
      contents={
        <Text color={text.light.BETA} size="kilo">
          {description}
        </Text>
      }
      style={{
        borderRadius: '3px',
        boxShadow: '0px 4px 8px -2px rgba(0, 0, 0, 0.3)',
        marginTop: '10px',
        padding: '16px',
        zIndex: 1,
        ...style,
      }}
      bottom
      {...props}
    >
      {children}
    </ElementalTooltip>
  </div>
);

Tooltip.propTypes = {
  description: PropTypes.string.isRequired,
  style: PropTypes.object,
  children: PropTypes.node.isRequired,
};

Tooltip.defaultProps = {
  style: {},
};

export default Tooltip;
