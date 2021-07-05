import PropTypes from 'prop-types';
import React from 'react';
import { pure } from 'recompose';
import { colors, Pulsator, Spacing, Text } from 'react-elemental';
import Tooltip from 'client/app/react/components/ui/tooltip';
import { text } from 'client/app/util/style/color';

export const CONNECTIVITY_STATE = {
  CONNECTED: 'connected',
  CONNECTING: 'connecting',
  DISCONNECTED: 'disconnected',
};

/**
 * Indicator of the current live connection status.
 */
const ConnectionStatus = ({ connectivity }) => {
  const states = {
    [CONNECTIVITY_STATE.DISCONNECTED]: {
      color: colors.red,
      description: 'There was an error establishing a live connection.',
      status: 'Not connected',
    },
    [CONNECTIVITY_STATE.CONNECTED]: {
      color: colors.green,
      description: 'A live connection is established.',
      status: 'Connected',
    },
    [CONNECTIVITY_STATE.CONNECTING]: {
      color: colors.gray10,
      description: 'Attempting to establish a live connection…',
      status: 'Connecting',
    },
  };

  const { color, description, status } = states[connectivity];

  return (
    <Tooltip
      description={description}
      style={{
        left: 'unset',
        right: 0,
        textAlign: 'right',
        width: '120px',
      }}
    >
      <div style={{ alignItems: 'center', display: 'flex', pointerEvents: 'none' }}>
        <Spacing size="small" style={{ display: 'inherit' }} right>
          <Pulsator
            color={color}
            size="delta"
            style={{ opacity: 0.75 }}
            inactive
          />
        </Spacing>

        <Text color={text.dark.BETA} size="kilo">
          {status}
        </Text>
      </div>
    </Tooltip>
  );
};

ConnectionStatus.propTypes = {
  connectivity: PropTypes.oneOf(Object.values(CONNECTIVITY_STATE)).isRequired,
};

export default pure(ConnectionStatus);
