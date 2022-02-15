import PropTypes from 'prop-types';
import React from 'react';
import { Spacing } from 'react-elemental';
import { MdStackedBarChart } from 'react-icons/md';
import { connect } from 'react-redux';
import AdminPanel from 'client/app/react/components/admin/panel';
import TelemetryEvent from 'client/app/react/components/admin/meta/telemetry-event';

/**
 * Telemetry event logs history section.
 */
const AdminMetaTelemetryContainer = ({ events }) => (
  <AdminPanel
    title="Telemetry"
    subtitle="Client event logs from the current session"
    icon={<MdStackedBarChart />}
  >
    {events.map((event, idx) => (
      <Spacing
        key={`${event.timestamp}-${event.action}-${event.value}`}
        size="small"
        bottom={idx < events.length - 1}
      >
        <TelemetryEvent
          timestamp={event.timestamp}
          action={event.action}
          value={event.value}
          tags={event.tags}
        />
      </Spacing>
    ))}
  </AdminPanel>
);

AdminMetaTelemetryContainer.propTypes = {
  events: PropTypes.array.isRequired,
};

const mapStateToProps = ({ telemetry }) => ({
  events: telemetry.events,
});

export default connect(mapStateToProps)(AdminMetaTelemetryContainer);
