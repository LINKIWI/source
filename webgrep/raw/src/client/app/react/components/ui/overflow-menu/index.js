import PropTypes from 'prop-types';
import React, { createRef, Component, Fragment } from 'react';
import { colors, Button, Text } from 'react-elemental';
import { compose, withToggleState } from '@linkiwi/hoc';
import ClickListener from 'client/app/react/components/passive/click-listener';
import { background, text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

/**
 * Overflow menu exposing a list of additional action items.
 */
class OverflowMenu extends Component {
  static propTypes = {
    buttonStyle: PropTypes.object,
    menuStyle: PropTypes.object,
    mirrored: PropTypes.bool,
    children: PropTypes.node.isRequired,
    // HOC props
    isMenuVisible: PropTypes.bool.isRequired,
    showMenu: PropTypes.func.isRequired,
    hideMenu: PropTypes.func.isRequired,
    isOverflowHover: PropTypes.bool.isRequired,
    handleOverflowMouseEnter: PropTypes.func.isRequired,
    handleOverflowMouseLeave: PropTypes.func.isRequired,
  };

  static defaultProps = {
    buttonStyle: {},
    menuStyle: {},
    mirrored: false,
  };

  ref = createRef();

  handleClick = this._handleClick.bind(this);

  _handleClick(evt) {
    const { isMenuVisible, hideMenu } = this.props;

    if (isMenuVisible && this.ref.current && !this.ref.current.contains(evt.target)) {
      hideMenu();
    }
  }

  render() {
    const {
      isMenuVisible,
      showMenu,
      hideMenu,
      isOverflowHover,
      handleOverflowMouseEnter,
      handleOverflowMouseLeave,
      buttonStyle: buttonStyleOverrides,
      menuStyle: menuStyleOverrides,
      mirrored,
      children,
    } = this.props;

    const buttonStyle = {
      borderRadius: '3px',
      display: 'block',
      filter: 'none',
      height: '22px',
      padding: '0 6px',
      ...isOverflowHover && {
        backgroundColor: background.dark.EPSILON,
      },
      ...buttonStyleOverrides,
    };

    const ellipsisStyle = {
      fontSize: '16px',
      lineHeight: 0,
      opacity: 0.4,
      transition: transition.all.DEFAULT,
      ...isOverflowHover && {
        opacity: 0.6,
      },
      ...isMenuVisible && {
        opacity: 1,
      },
    };

    const menuStyle = {
      backgroundColor: 'white',
      borderRadius: '3px',
      boxShadow: '0px 8px 16px -2px rgba(0, 0, 0, 0.1)',
      marginTop: '2px',
      minWidth: '80px',
      opacity: 0,
      overflow: 'hidden',
      position: 'absolute',
      right: 0,
      transform: 'perspective(600px) rotateX(-15deg) rotateY(-5deg)',
      transformOrigin: 'right top',
      transition: transition.all.DEFAULT,
      visibility: 'hidden',
      zIndex: 1,
      ...mirrored && {
        left: 0,
        right: 'unset',
        transform: 'perspective(600px) rotateX(-15deg) rotateY(5deg)',
        transformOrigin: 'left top',
      },
      ...isMenuVisible && {
        marginTop: '6px',
        opacity: 1,
        visibility: 'visible',
        transform: 'perspective(600px) rotate(0deg)',
      },
      ...menuStyleOverrides,
    };

    return (
      <Fragment>
        <ClickListener handler={this.handleClick} />

        <div style={{ position: 'relative' }}>
          <Button
            ref={this.ref}
            color="transparent"
            style={buttonStyle}
            onClick={isMenuVisible ? hideMenu : showMenu}
            onMouseEnter={handleOverflowMouseEnter}
            onMouseLeave={handleOverflowMouseLeave}
          >
            <Text color={isMenuVisible ? colors.primary : text.dark.ALPHA} style={ellipsisStyle}>
              ···
            </Text>
          </Button>

          <div style={menuStyle}>
            {children}
          </div>
        </div>
      </Fragment>
    );
  }
}

export default compose(
  withToggleState({
    key: 'isMenuVisible',
    enable: 'showMenu',
    disable: 'hideMenu',
    allowOverride: true,
  }),
  withToggleState({
    key: 'isOverflowHover',
    enable: 'handleOverflowMouseEnter',
    disable: 'handleOverflowMouseLeave',
  }),
)(OverflowMenu);
