import PropTypes from 'prop-types';
import { PureComponent } from 'react';

/**
 * Passive listener exposing callbacks when the page visibility changes.
 */
export default class PresenceListener extends PureComponent {
  static propTypes = {
    onReveal: PropTypes.func,
    onHide: PropTypes.func,
    onUnload: PropTypes.func,
  };

  static defaultProps = {
    onReveal: () => {},
    onHide: () => {},
    onUnload: () => {},
  };

  componentDidMount() {
    window.addEventListener('visibilitychange', this.handleVisibilityChange);
    window.addEventListener('beforeunload', this.handleUnload);

    // The initial state after mount is "revealed," so invoke the callback immediately.
    this.props.onReveal();
  }

  componentWillUnmount() {
    window.removeEventListener('visibilitychange', this.handleVisibilityChange);
    window.removeEventListener('beforeunload', this.handleUnload);

    // Unmounting the component implies that the page is hidden.
    this.props.onHide();
  }

  handleVisibilityChange = this._handleVisibilityChange.bind(this);

  handleUnload = this.props.onUnload.bind(this);

  _handleVisibilityChange() {
    const { onReveal, onHide } = this.props;

    if (document.hidden) {
      onHide();
    } else {
      onReveal();
    }
  }

  render() {
    return null;
  }
}
