import PropTypes from 'prop-types';
import React, { Fragment } from 'react';
import { colors, Spacing } from 'react-elemental';
import { connect } from 'react-redux';
import { compose, withToggleState } from '@linkiwi/hoc';
import Banner from 'client/app/react/components/header/banner';
import Logo from 'client/app/react/components/header/logo';
import Layout from 'client/app/react/components/common/layout';
import NavigationTabs, { NAVIGATION_TABS } from 'client/app/react/components/header/navigation-tabs';
import { background } from 'client/app/util/style/color';
import { objLookup } from 'shared/util/data';

const HeaderContainer = ({
  pathname,
  width,
  indexErr,
  logo,
  banner,
  isInfoBannerVisible,
  isErrorBannerVisible,
  hideInfoBanner,
  hideErrorBanner,
}) => {
  const tab = NAVIGATION_TABS.find((descriptor) => pathname.startsWith(descriptor.href)) || {};

  return (
    <Fragment>
      {banner.title && banner.description && isInfoBannerVisible && (
        <Layout width={width} containerStyle={{ backgroundColor: colors.primary }}>
          <Banner
            title={banner.title}
            description={banner.description}
            link={banner.link}
            href={banner.href}
            onHide={hideInfoBanner}
          />
        </Layout>
      )}

      {indexErr && isErrorBannerVisible && (
        <Layout width={width} containerStyle={{ backgroundColor: colors.red }}>
          <Banner
            title="Index error"
            description={[
              'There was an error communicating with the livegrep backend.',
              'Is the codesearch process running?',
            ].join(' ')}
            onHide={hideErrorBanner}
          />
        </Layout>
      )}

      <Layout width={width} containerStyle={{ backgroundColor: background.dark.ALPHA }}>
        <div style={{ alignItems: 'center', display: 'flex', justifyContent: 'space-between' }}>
          <Spacing size="large" right>
            <Logo src={logo} />
          </Spacing>

          <NavigationTabs tab={tab.value} />
        </div>
      </Layout>
    </Fragment>
  );
};

HeaderContainer.propTypes = {
  pathname: PropTypes.string.isRequired,
  width: PropTypes.number.isRequired,
  // HOC props
  indexErr: PropTypes.bool.isRequired,
  logo: PropTypes.string,
  banner: PropTypes.shape({
    title: PropTypes.string,
    description: PropTypes.string,
    link: PropTypes.string,
    href: PropTypes.string,
  }).isRequired,
  isInfoBannerVisible: PropTypes.bool.isRequired,
  isErrorBannerVisible: PropTypes.bool.isRequired,
  hideInfoBanner: PropTypes.func.isRequired,
  hideErrorBanner: PropTypes.func.isRequired,
};

HeaderContainer.defaultProps = {
  logo: null,
};

const mapStateToProps = ({ config, meta }) => ({
  indexErr: meta.timestamp === 0,
  logo: objLookup(config, ['client', 'site', 'logo']),
  banner: objLookup(config, ['client', 'site', 'banner']) || {},
});

export default compose(
  connect(mapStateToProps),
  withToggleState({
    key: 'isInfoBannerVisible',
    disable: 'hideInfoBanner',
    initial: true,
  }),
  withToggleState({
    key: 'isErrorBannerVisible',
    disable: 'hideErrorBanner',
    initial: true,
  }),
)(HeaderContainer);
