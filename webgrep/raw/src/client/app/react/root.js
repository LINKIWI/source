import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import Favicon from 'react-favicon';
import { connect } from 'react-redux';
import { Redirect, Route, Switch } from 'react-router';
import { withRouter } from 'react-router-dom';
import { bindActionCreators } from 'redux';
import { compose, withCSS, withWindowDimensions } from '@linkiwi/hoc';
import AdminContainer from 'client/app/react/containers/admin';
import HeaderContainer from 'client/app/react/containers/header';
import withTelemetry from 'client/app/react/containers/hoc/with-telemetry';
import ErrorBoundaryContainer from 'client/app/react/containers/meta/error-boundary';
import ToastContainer from 'client/app/react/containers/meta/toast';
import MetadataPollContainer from 'client/app/react/containers/meta/metadata-poll';
import SearchContainer from 'client/app/react/containers/search';
import PresenceListener from 'client/app/react/components/passive/presence-listener';
import { setWindowDimensions } from 'client/app/redux/actions/context';
import { text } from 'client/app/util/style/color';
import logoResource from 'client/resources/img/logo.png';
import { TELEMETRY_ACTIONS } from 'shared/constants/telemetry';
import { objLookup } from 'shared/util/data';
import { stopwatch } from 'shared/util/instrumentation';

/**
 * Application root.
 */
class Root extends Component {
  static propTypes = {
    // HOC props
    metadataPollInterval: PropTypes.number.isRequired,
    location: PropTypes.shape({
      pathname: PropTypes.string.isRequired,
    }).isRequired,
    history: PropTypes.shape({
      listen: PropTypes.func.isRequired,
      location: PropTypes.shape({
        pathname: PropTypes.string.isRequired,
      }).isRequired,
    }).isRequired,
    width: PropTypes.number.isRequired,
    height: PropTypes.number.isRequired,
    recordTelemetryEvent: PropTypes.func.isRequired,
    serverTimestamp: PropTypes.number.isRequired,
    actions: PropTypes.shape({
      setWindowDimensions: PropTypes.func.isRequired,
    }).isRequired,
  };

  constructor(...args) {
    super(...args);

    this.totalSessionStopwatch = stopwatch();
    this.activeSessionStopwatch = stopwatch();
  }

  componentDidMount() {
    const { width, height, actions, history } = this.props;

    actions.setWindowDimensions(width, height);

    this._recordRootRender();
    this._recordRouteRender(history.location.pathname);

    this.unlistenHistory = history.listen((location, action) => {
      if (['PUSH', 'REPLACE'].includes(action)) {
        this._recordRouteRender(location.pathname);
      }
    });
  }

  componentDidUpdate(prevProps) {
    const { width, height, actions } = this.props;

    if (width !== prevProps.width || height !== prevProps.height) {
      actions.setWindowDimensions(width, height);
    }
  }

  componentWillUnmount() {
    this.unlistenHistory();
  }

  handlePageReveal = this._handlePageReveal.bind(this);

  handlePageHide = this._handlePageHide.bind(this);

  handlePageUnload = this._handlePageUnload.bind(this);

  _recordRootRender() {
    const { recordTelemetryEvent, serverTimestamp } = this.props;

    recordTelemetryEvent(TELEMETRY_ACTIONS.RENDER_DELAY, Date.now() - serverTimestamp);
  }

  _recordRouteRender(path) {
    this.props.recordTelemetryEvent(TELEMETRY_ACTIONS.RENDER_VIEW_ROUTE, 1, { path });
  }

  _handlePageReveal() {
    // A new active session is started when the page is revealed; reset the current stopwatch.
    this.activeSessionStopwatch = stopwatch();
  }

  _handlePageHide() {
    this.props.recordTelemetryEvent(
      TELEMETRY_ACTIONS.ACTIVE_SESSION_LENGTH,
      this.activeSessionStopwatch.elapsed(),
    );
  }

  _handlePageUnload() {
    this.handlePageHide();
    this.props.recordTelemetryEvent(
      TELEMETRY_ACTIONS.TOTAL_SESSION_LENGTH,
      this.totalSessionStopwatch.elapsed(),
    );
  }

  render() {
    const { location, width, metadataPollInterval } = this.props;

    return (
      <Fragment>
        <Favicon url={logoResource} />

        <HeaderContainer
          pathname={location.pathname}
          width={width}
        />

        <Switch>
          <Route path="/search" component={SearchContainer} />
          <Route path="/admin" component={AdminContainer} />
          <Redirect from="*" to="/search" />
        </Switch>

        <ToastContainer />

        <ErrorBoundaryContainer>
          {metadataPollInterval > 0 && (
            <MetadataPollContainer pollInterval={metadataPollInterval} />
          )}
        </ErrorBoundaryContainer>

        <PresenceListener
          onReveal={this.handlePageReveal}
          onHide={this.handlePageHide}
          onUnload={this.handlePageUnload}
        />
      </Fragment>
    );
  }
}

const mapStateToProps = ({ config, context }) => ({
  metadataPollInterval: objLookup(config, ['client', 'options', 'metadata_poll_interval']) || 0,
  serverTimestamp: objLookup(context, ['timestamp']) || Date.now(),
});

const mapDispatchToProps = (dispatch) => ({
  actions: bindActionCreators({ setWindowDimensions }, dispatch),
});

export default compose(
  withRouter,
  connect(mapStateToProps, mapDispatchToProps),
  withWindowDimensions,
  withTelemetry,
  withCSS({
    key: () => 'selection',
    css: () => `::selection {background: ${text.highlight.SELECT}}`,
  }),
)(Root);
