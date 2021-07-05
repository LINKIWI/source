import PropTypes from 'prop-types';
import React, { createRef, Component, Fragment } from 'react';
import { sizes, Button, Spacing, Text } from 'react-elemental';
import { MdKeyboardArrowDown } from 'react-icons/md';
import { compose, withToggleState } from '@linkiwi/hoc';
import ClickListener from 'client/app/react/components/passive/click-listener';
import ButtonMenuItem from 'client/app/react/components/ui/overflow-menu/button-menu-item';
import { background, text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

class SelectList extends Component {
  static propTypes = {
    value: PropTypes.string.isRequired,
    options: PropTypes.arrayOf(PropTypes.shape({
      value: PropTypes.string.isRequired,
      label: PropTypes.string.isRequired,
    })).isRequired,
    style: PropTypes.object,
    onChange: PropTypes.func.isRequired,
    // HOC props
    isOpen: PropTypes.bool.isRequired,
    isHover: PropTypes.bool.isRequired,
    open: PropTypes.func.isRequired,
    hide: PropTypes.func.isRequired,
    handleMouseEnter: PropTypes.func.isRequired,
    handleMouseLeave: PropTypes.func.isRequired,
  };

  static defaultProps = {
    style: {},
  };

  ref = createRef();

  handleClick = this._handleClick.bind(this);

  _handleClick(evt) {
    const { isOpen, hide } = this.props;

    if (isOpen && this.ref.current && !this.ref.current.contains(evt.target)) {
      hide();
    }
  }

  render() {
    const {
      value,
      options,
      onChange,
      isOpen,
      open,
      hide,
      style,
      isHover,
      handleMouseEnter,
      handleMouseLeave,
    } = this.props;

    const containerStyle = {
      display: 'block',
      position: 'relative',
      width: '100%',
      ...style,
    };

    const buttonStyle = {
      alignItems: 'center',
      border: '1px solid transparent',
      borderRadius: '3px',
      display: 'flex',
      filter: 'none',
      justifyContent: 'space-between',
      width: '100%',
      ...isHover && {
        backgroundColor: background.dark.IOTA,
      },
      ...isOpen && {
        backgroundColor: 'transparent',
        border: `1px solid ${background.primary.BETA}`,
      },
    };

    const textStyle = {
      color: text.dark.BETA,
      textAlign: 'left',
      transition: transition.all.DEFAULT,
      ...(isHover || isOpen) && {
        color: text.dark.ALPHA,
      },
    };

    const iconStyle = {
      color: text.dark.GAMMA,
      fontSize: sizes.kilo,
      transition: transition.all.DEFAULT,
      ...isHover && {
        color: text.dark.BETA,
      },
      ...isOpen && {
        color: text.dark.BETA,
        transform: 'rotate(180deg)',
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
      transform: 'perspective(600px) rotateX(-15deg)',
      transformOrigin: 'top',
      transition: transition.all.DEFAULT,
      visibility: 'hidden',
      width: 'inherit',
      zIndex: 1,
      ...isOpen && {
        marginTop: '6px',
        opacity: 1,
        visibility: 'visible',
        transform: 'perspective(600px) rotate(0deg)',
      },
    };

    const selected = options.find((option) => option.value === value) || {};

    return (
      <Fragment>
        <ClickListener handler={this.handleClick} />

        <div style={containerStyle}>
          <Button
            ref={this.ref}
            color={background.dark.EPSILON}
            style={buttonStyle}
            onMouseEnter={handleMouseEnter}
            onMouseLeave={handleMouseLeave}
            onClick={isOpen ? hide : open}
          >
            <Text size="kilo" style={textStyle}>
              {selected.label || 'Select an item…'}
            </Text>

            <Spacing size="small" style={{ display: 'inherit' }} left>
              <MdKeyboardArrowDown style={iconStyle} />
            </Spacing>
          </Button>

          <div style={menuStyle}>
            {options.map((option) => (
              <ButtonMenuItem
                key={option.value}
                label={option.label}
                onClick={() => option.value !== value && onChange(option.value)}
              />
            ))}
          </div>
        </div>
      </Fragment>
    );
  }
}

export default compose(
  withToggleState({
    key: 'isOpen',
    enable: 'open',
    disable: 'hide',
  }),
  withToggleState({
    key: 'isHover',
    enable: 'handleMouseEnter',
    disable: 'handleMouseLeave',
  }),
)(SelectList);
