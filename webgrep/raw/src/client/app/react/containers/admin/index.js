import PropTypes from 'prop-types';
import React, { Fragment } from 'react';
import { Spacing, Text } from 'react-elemental';
import { Helmet } from 'react-helmet';
import { connect } from 'react-redux';
import { Redirect, Route, Switch } from 'react-router';
import AdminIndexContainer from 'client/app/react/containers/admin/livegrep';
import AdminMetaContainer from 'client/app/react/containers/admin/meta';
import AdminPreferencesContainer from 'client/app/react/containers/admin/preferences';
import AdminNavigationTab from 'client/app/react/components/admin/navigation-tab';
import Layout from 'client/app/react/components/common/layout';
import { background, text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';
import { string } from 'client/app/util/format';
import { objLookup } from 'shared/util/data';

// Window widths at which to consider the layout to be horizontally compact.
const breakpoints = [
  { threshold: 800, compact: true },
];

/**
 * Container for all admin pages.
 */
const AdminContainer = ({ title, width, location: { pathname } }) => {
  const isCompact = breakpoints.reduce((acc, breakpoint) =>
    (width < breakpoint.threshold ? breakpoint.compact : acc), false);

  const layoutStyle = {
    alignItems: 'top',
    display: 'flex',
    ...isCompact && {
      flexWrap: 'wrap',
    },
  };

  const sidebarStyle = {
    transition: transition.all.DEFAULT,
    width: '250px',
    ...isCompact && {
      width: '100%',
    },
  };

  return (
    <Fragment>
      <Helmet>
        <title>
          {title}
        </title>
      </Helmet>

      <Spacing size="large" style={{ backgroundColor: background.dark.ALPHA }} top bottom padding>
        <Layout width={width}>
          <Spacing size="tiny" bottom>
            <Text color={text.light.ALPHA} size="delta" bold>
              Admin
            </Text>
          </Spacing>

          <Text color={text.light.EPSILON}>
            Manage webgrep and livegrep behavior
          </Text>
        </Layout>
      </Spacing>

      <Spacing size="large" top>
        <Layout width={width} contentStyle={layoutStyle}>
          <Spacing
            size="large"
            style={sidebarStyle}
            right={!isCompact}
            bottom
          >
            <AdminNavigationTab
              label="Index"
              href="/admin/index"
              isSelected={pathname.startsWith('/admin/index')}
            />
            <AdminNavigationTab
              label="Preferences"
              href="/admin/preferences"
              isSelected={pathname.startsWith('/admin/preferences')}
            />
            <AdminNavigationTab
              label="Meta"
              href="/admin/meta"
              isSelected={pathname.startsWith('/admin/meta')}
            />
          </Spacing>

          <div style={{ flexGrow: 1 }}>
            <Switch>
              <Route path="/admin/index" component={AdminIndexContainer} />
              <Route path="/admin/preferences" component={AdminPreferencesContainer} />
              <Route path="/admin/meta" component={AdminMetaContainer} />
              <Redirect from="*" to="/admin/index" exact />
            </Switch>
          </div>
        </Layout>
      </Spacing>
    </Fragment>
  );
};

AdminContainer.propTypes = {
  location: PropTypes.shape({
    pathname: PropTypes.string.isRequired,
  }).isRequired,
  // HOC props
  title: PropTypes.string.isRequired,
  width: PropTypes.number,
};

AdminContainer.defaultProps = {
  width: null,
};

const mapStateToProps = ({ config, context, meta }) => ({
  title: string(objLookup(config, ['client', 'site', 'title']) || 'webgrep', {
    page: 'Admin',
    name: meta.name,
  }),
  width: context.window.width,
});

export default connect(mapStateToProps)(AdminContainer);
