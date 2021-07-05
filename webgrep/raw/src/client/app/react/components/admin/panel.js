import PropTypes from 'prop-types';
import React from 'react';
import { Spacing, Text } from 'react-elemental';
import Box from 'client/app/react/components/common/box';
import { background, text } from 'client/app/util/style/color';

/**
 * Single admin section panel, with a standardized header style.
 */
const AdminPanel = ({ icon, title, subtitle, children }) => (
  <Box style={{ padding: 0 }}>
    <Spacing
      style={{ alignItems: 'top', backgroundColor: background.dark.EPSILON, display: 'flex' }}
      top
      right
      bottom
      left
      padding
    >
      <Spacing style={{ opacity: 0.8 }} right>
        {icon}
      </Spacing>

      <div>
        <Spacing size="micro" bottom>
          <Text size="kilo" uppercase bold>
            {title}
          </Text>
        </Spacing>

        <Text color={text.dark.BETA} size="kilo">
          {subtitle}
        </Text>
      </div>
    </Spacing>

    <Spacing top right bottom left padding>
      {children}
    </Spacing>
  </Box>
);

AdminPanel.propTypes = {
  icon: PropTypes.node.isRequired,
  title: PropTypes.string.isRequired,
  subtitle: PropTypes.string.isRequired,
  children: PropTypes.node.isRequired,
};

export default AdminPanel;
