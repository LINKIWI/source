import PropTypes from 'prop-types';
import { PureComponent } from 'react';

/**
 * Run a handler on keydown for a specific character.
 */
export default class KeyboardListener extends PureComponent {
  static propTypes = {
    character: PropTypes.string.isRequired,
    handler: PropTypes.func.isRequired,
  };

  componentDidMount() {
    document.addEventListener('keydown', this.handleKeypress);
  }

  componentWillUnmount() {
    document.removeEventListener('keydown', this.handleKeypress);
  }

  handleKeypress = this._handleKeypress.bind(this);

  _handleKeypress(evt) {
    const { character, handler } = this.props;

    if (evt.key === character) {
      handler(evt);
    }
  }

  render() {
    return null;
  }
}
