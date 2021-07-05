import PropTypes from 'prop-types';
import React, { createRef, Component, Fragment } from 'react';
import { spacing } from 'react-elemental';
import Box from 'client/app/react/components/common/box';
import ClickListener from 'client/app/react/components/passive/click-listener';
import KeyboardListener from 'client/app/react/components/passive/keyboard-listener';
import { transition } from 'client/app/util/style/transition';

/**
 * A "sticky modal" attached to an existing element.
 */
class Flyout extends Component {
  static propTypes = {
    triggerRef: PropTypes.oneOfType([
      PropTypes.shape({ current: PropTypes.instanceOf(Element) }),
      PropTypes.func,
    ]).isRequired,
    isOpen: PropTypes.bool.isRequired,
    isCompact: PropTypes.bool.isRequired,
    onHide: PropTypes.func.isRequired,
    children: PropTypes.node.isRequired,
  };

  ref = createRef();

  handleClick = this._handleClick.bind(this);

  _handleClick(evt) {
    const { triggerRef, isOpen, onHide } = this.props;

    // Close the flyout if something was clicked outside of the flyout bounds, but only if it isn't
    // the DOM element that triggers the opening of this flyout.
    const isOutsideCloseClick = isOpen &&
      this.ref.current &&
      !this.ref.current.contains(evt.target) &&
      !triggerRef.current.contains(evt.target);

    if (isOutsideCloseClick) {
      onHide();
    }
  }

  render() {
    const { isOpen, isCompact, onHide, children } = this.props;

    const containerStyle = {
      boxSizing: 'border-box',
      left: '-18px',
      opacity: 0,
      position: 'absolute',
      transform: 'perspective(1000px) rotateX(-5deg) rotateY(5deg)',
      transformOrigin: 'left top',
      transition: transition.all.DEFAULT,
      visibility: 'hidden',
      width: '440px',
      zIndex: 2,
      ...isCompact && {
        left: 0,
        transform: 'perspective(1000px) rotateX(-5deg)',
        transformOrigin: 'center top',
        width: '100%',
      },
      ...isOpen && {
        visibility: 'inherit',
        opacity: 1,
        transform: 'perspective(1000px) rotate(0deg)',
      },
    };

    const style = {
      margin: spacing.small,
    };

    return (
      <Fragment>
        <KeyboardListener
          character="Escape"
          handler={onHide}
        />
        <ClickListener
          handler={this.handleClick}
        />

        <div ref={this.ref} style={containerStyle}>
          <Box variant="alpha" style={style}>
            {children}
          </Box>
        </div>
      </Fragment>
    );
  }
}

export default Flyout;
