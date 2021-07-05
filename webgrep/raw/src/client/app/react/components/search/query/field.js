import PropTypes from 'prop-types';
import React, { createRef, Component, Fragment } from 'react';
import { colors, sizes, Button } from 'react-elemental';
import { MdClear, MdSearch } from 'react-icons/md';
import KeyboardListener from 'client/app/react/components/passive/keyboard-listener';
import TextField from 'client/app/react/components/ui/text-field';
import { NODE_ID, INPUT_FIELD_IDS } from 'client/app/util/constants/dom';
import { background, text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

/**
 * Search query field, with optional annotations.
 */
class QueryField extends Component {
  static propTypes = {
    value: PropTypes.string.isRequired,
    onChange: PropTypes.func.isRequired,
  };

  field = createRef();

  handleKeyboardFocus = this._handleKeyboardFocus.bind(this);

  handleClear = this._handleClear.bind(this);

  _handleKeyboardFocus(evt) {
    if (!INPUT_FIELD_IDS.includes(evt.target.id)) {
      evt.preventDefault();
      this.field.current.focus();
    }
  }

  _handleClear() {
    this.props.onChange({ target: { value: '' } });
  }

  render() {
    const { value, ...props } = this.props;

    const before = (
      <MdSearch style={{ color: colors.primary, fontSize: '20px' }} />
    );

    const after = (
      <Button
        color={background.dark.DELTA}
        onClick={this.handleClear}
        style={{
          alignItems: 'center',
          borderRadius: '50%',
          display: 'flex',
          justifyContent: 'center',
          height: '18px',
          opacity: 0,
          padding: 0,
          right: 0,
          transition: transition.all.DEFAULT,
          width: '18px',
          ...value && {
            opacity: 0.6,
          },
        }}
      >
        <MdClear style={{ color: text.dark.ALPHA, fontSize: sizes.lambda }} />
      </Button>
    );

    const afterEnhancerStyle = {
      pointerEvents: 'unset',
      ...value && { zIndex: 0 },
    };

    return (
      <Fragment>
        <KeyboardListener
          character="/"
          handler={this.handleKeyboardFocus}
        />

        <TextField
          ref={this.field}
          before={before}
          after={after}
          afterEnhancerStyle={afterEnhancerStyle}
          id={NODE_ID.SEARCH_QUERY_FIELD}
          placeholder="Search for code…"
          autoComplete="off"
          value={value}
          autoFocus
          {...props}
        />
      </Fragment>
    );
  }
}

export default QueryField;
