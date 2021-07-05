import PropTypes from 'prop-types';
import React from 'react';
import { Spacing, Text } from 'react-elemental';

/**
 * Key-value metadata item in index information.
 */
const MetadataItem = ({ heading, value }) => (
  <div style={{ alignItems: 'top', display: 'flex', justifyContent: 'space-between' }}>
    <Text size="kilo">
      {heading}
    </Text>

    <Spacing left>
      <Text size="kilo" bold right>
        {value}
      </Text>
    </Spacing>
  </div>
);

MetadataItem.propTypes = {
  heading: PropTypes.string.isRequired,
  value: PropTypes.string.isRequired,
};

export default MetadataItem;
