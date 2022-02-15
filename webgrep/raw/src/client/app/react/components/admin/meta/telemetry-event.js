import PropTypes from 'prop-types';
import React, { Fragment } from 'react';
import { Spacing, Text } from 'react-elemental';
import Tooltip from 'client/app/react/components/ui/tooltip';
import { background, text } from 'client/app/util/style/color';
import { absoluteTimestamp } from 'client/app/util/format';

/**
 * Representation of a telemetry event log entry.
 */
const TelemetryEvent = ({ timestamp, action, value, tags }) => {
  const event = (
    <Fragment>
      <Spacing size="micro" style={{ display: 'flex' }} bottom>
        <Spacing right>
          <Spacing size="tiny" top>
            <div
              style={{
                backgroundColor: background.primary.BETA,
                borderRadius: '50%',
                height: '8px',
                width: '8px',
              }}
            />
          </Spacing>
        </Spacing>

        <div style={{ flexGrow: 1 }}>
          <div>
            <Spacing size="small" right inline>
              <Text size="kilo" secondary bold inline>
                {action}
              </Text>
            </Spacing>

            <Text color={text.dark.BETA} size="kilo" secondary inline>
              {value}
            </Text>
          </div>

          <Text color={text.dark.BETA} size="kilo">
            {absoluteTimestamp(timestamp / 1000)}
          </Text>
        </div>
      </Spacing>
    </Fragment>
  );

  if (!Object.keys(tags).length) {
    return event;
  }

  const tooltip = (
    <Fragment>
      {Object.entries(tags).map(([tagKey, tagValue]) => (
        <div key={tagKey}>
          <Spacing size="small" right inline>
            <Text color={text.light.ALPHA} size="kilo" bold inline>
              {tagKey}
            </Text>
          </Spacing>

          <Text
            color={text.light.BETA}
            size="kilo"
            style={{ whiteSpace: 'nowrap' }}
            secondary
            inline
          >
            {tagValue}
          </Text>
        </div>
      ))}
    </Fragment>
  );

  return (
    <Tooltip
      description=""
      contents={tooltip}
      style={{ minWidth: '100px' }}
      offset={30}
      bottom={false}
      top
    >
      {event}
    </Tooltip>
  );
};

TelemetryEvent.propTypes = {
  timestamp: PropTypes.number.isRequired,
  action: PropTypes.string.isRequired,
  value: PropTypes.number.isRequired,
  tags: PropTypes.object.isRequired,
};

export default TelemetryEvent;
