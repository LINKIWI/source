import PropTypes from 'prop-types';
import React from 'react';
import { Spacing, Text } from 'react-elemental';
import { text } from 'client/app/util/style/color';

/**
 * Key-value metadata item in index information.
 */
const MetadataItem = ({ heading, value, isCompact }) => {
  const containerStyle = {
    alignItems: 'top',
    display: 'flex',
    justifyContent: 'space-between',
    ...isCompact && {
      flexDirection: 'column',
    },
  };

  return (
    <div style={containerStyle}>
      <Spacing size="tiny" bottom={isCompact}>
        <Text color={text.dark.BETA} size="kilo">
          {heading}
        </Text>
      </Spacing>

      <Spacing left={!isCompact}>
        <Text size="kilo" style={{ wordBreak: 'break-word' }} bold right={!isCompact}>
          {value}
        </Text>
      </Spacing>
    </div>
  );
};

MetadataItem.propTypes = {
  heading: PropTypes.string.isRequired,
  value: PropTypes.string.isRequired,
  isCompact: PropTypes.bool.isRequired,
};

export default MetadataItem;
