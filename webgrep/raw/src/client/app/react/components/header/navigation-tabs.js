import PropTypes from 'prop-types';
import React from 'react';
import { Spacing, Tabs, Text } from 'react-elemental';
import { Link } from 'react-router-dom';
import { text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

// Description of the values, labels, and navigation targets of all tabs, in order.
export const NAVIGATION_TABS = [
  { value: 'search', label: 'Search', href: '/search' },
  { value: 'admin', label: 'Admin', href: '/admin' },
];

/**
 * Single anchor-enabled tab.
 */
const Tab = ({ label, href, isSelected }) => (
  <Link to={href} style={{ textDecoration: 'none' }}>
    <Spacing top bottom padding>
      <Spacing size="small" left right padding>
        <Text
          size="kilo"
          bold={isSelected}
          color={isSelected ? text.light.ALPHA : text.light.BETA}
          style={{
            opacity: isSelected ? 1 : 0.7,
            transition: transition.all.DEFAULT,
          }}
          uppercase
        >
          {label}
        </Text>
      </Spacing>
    </Spacing>
  </Link>
);

Tab.propTypes = {
  label: PropTypes.string.isRequired,
  href: PropTypes.string.isRequired,
  isSelected: PropTypes.bool.isRequired,
};

/**
 * Set of switchable navigation tabs.
 */
const NavigationTabs = ({ tab }) => (
  <Tabs
    value={tab}
    options={NAVIGATION_TABS.map((descriptor) => ({
      value: descriptor.value,
      label: (
        <Tab
          label={descriptor.label}
          href={descriptor.href}
          isSelected={tab === descriptor.value}
        />
      ),
    }))}
    secondary
    invert
    fit
  />
);

NavigationTabs.propTypes = {
  tab: PropTypes.oneOf(NAVIGATION_TABS.map((descriptor) => descriptor.value)),
};

NavigationTabs.defaultProps = {
  tab: null,
};

export default NavigationTabs;
