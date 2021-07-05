import PropTypes from 'prop-types';
import { PureComponent } from 'react';

/**
 * Run a handler on a click anywhere within the viewport.
 */
export default class ClickListener extends PureComponent {
  static propTypes = {
    handler: PropTypes.func.isRequired,
  };

  componentDidMount() {
    window.addEventListener('click', this.handleClick);
  }

  componentWillUnmount() {
    window.removeEventListener('click', this.handleClick);
  }

  handleClick = this.props.handler.bind(this);

  render() {
    return null;
  }
}
