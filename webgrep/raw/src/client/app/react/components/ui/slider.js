import PropTypes from 'prop-types';
import React, { createRef, Component } from 'react';
import { Text } from 'react-elemental';
import { compose, withToggleState } from '@linkiwi/hoc';
import { background, text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

class Slider extends Component {
  static propTypes = {
    value: PropTypes.number.isRequired,
    annotation: PropTypes.string.isRequired,
    onChange: PropTypes.func.isRequired,
    onFinalize: PropTypes.func.isRequired,
    // HOC props
    isAnchorHover: PropTypes.bool.isRequired,
    handleAnchorMouseEnter: PropTypes.func.isRequired,
    handleAnchorMouseLeave: PropTypes.func.isRequired,
    isDragging: PropTypes.bool.isRequired,
    startDrag: PropTypes.func.isRequired,
    endDrag: PropTypes.func.isRequired,
  };

  container = createRef();

  handleMouseLeave = this._handleMouseLeave.bind(this);

  handleMouseMove = this._handleMouseMove.bind(this);

  handleTouchMove = this._handleTouchMove.bind(this);

  updateDrag = this._updateDrag.bind(this);

  finalizeDrag = this._finalizeDrag.bind(this);

  _handleMouseLeave(evt) {
    if (evt.relatedTarget !== window) {
      this.finalizeDrag();
    }
  }

  _handleMouseMove(evt) {
    if (this.props.isDragging) {
      this.updateDrag(evt);
    }
  }

  _handleTouchMove(evt) {
    const [touch] = evt.touches;
    if (!touch) {
      return;
    }

    this.updateDrag(touch);
  }

  _updateDrag(evt) {
    if (!this.container.current) {
      return;
    }

    const bounds = this.container.current.getBoundingClientRect();
    const offset = (evt.clientX - bounds.left) / (bounds.right - bounds.left);

    if (offset > 0 && offset < 1) {
      this.props.onChange(offset);
    }
  }

  _finalizeDrag() {
    const { endDrag, onFinalize } = this.props;

    endDrag();
    onFinalize();
  }

  render() {
    const {
      value,
      annotation,
      isDragging,
      isAnchorHover,
      handleAnchorMouseEnter,
      handleAnchorMouseLeave,
      startDrag,
    } = this.props;

    const containerStyle = {
      alignItems: 'center',
      cursor: 'grab',
      display: 'flex',
      height: '35px',
      touchAction: 'none',  // Don't scroll the page while touchmove is being fired
      width: '100%',
      ...isDragging && {
        cursor: 'grabbing',
      },
    };

    const sliderStyle = {
      backgroundColor: background.dark.DELTA,
      borderRadius: '2px',
      height: '2px',
      position: 'relative',
      width: '100%',
    };

    const annotationContainerStyle = {
      left: `${100 * value}%`,
      marginLeft: '-25px',
      marginTop: '-25px',
      pointerEvents: 'none',
      position: 'absolute',
      MozUserSelect: 'none',
      MsUserSelect: 'none',
      WebkitUserSelect: 'none',
      userSelect: 'none',
      transition: transition.all.DEFAULT,
      width: '50px',
      ...(isAnchorHover || isDragging) && {
        transition: 'unset',
      },
    };

    const annotationStyle = {
      color: text.dark.GAMMA,
      transition: transition.all.DEFAULT,
      ...(isAnchorHover || isDragging) && {
        color: text.dark.BETA,
      },
    };

    const anchorContainerStyle = {
      alignItems: 'center',
      backgroundColor: 'white',
      display: 'flex',
      height: '16px',
      justifyContent: 'center',
      left: `${100 * value}%`,
      marginLeft: '-8px',
      marginTop: '-7px',
      position: 'absolute',
      transition: transition.all.DEFAULT,
      width: '16px',
      ...(isAnchorHover || isDragging) && {
        transition: 'unset',
      },
    };

    const anchorStyle = {
      backgroundColor: background.primary.BETA,
      borderRadius: '1px',
      boxShadow: '0px 2px 3px 0px rgba(0, 0, 0, 0.15)',
      height: '12px',
      width: '12px',
      transition: transition.all.DEFAULT,
      ...isAnchorHover && {
        backgroundColor: background.primary.ALPHA,
      },
      ...isDragging && {
        backgroundColor: background.primary.ALPHA,
        boxShadow: '0px 2px 3px 0px rgba(0, 0, 0, 0.25)',
      },
    };

    return (
      <div
        ref={this.container}
        onMouseDown={compose(startDrag, this.updateDrag)}
        onMouseUp={this.finalizeDrag}
        onMouseLeave={this.handleMouseLeave}
        onMouseMove={this.handleMouseMove}
        onTouchStart={compose(startDrag, this.handleTouchMove)}
        onTouchEnd={this.finalizeDrag}
        onTouchCancel={this.finalizeDrag}
        onTouchMove={this.handleTouchMove}
        style={containerStyle}
      >
        <div style={sliderStyle}>
          <div style={annotationContainerStyle}>
            <Text size="lambda" style={annotationStyle} bold center>
              {annotation}
            </Text>
          </div>

          <div
            onMouseEnter={handleAnchorMouseEnter}
            onMouseLeave={handleAnchorMouseLeave}
            onMouseUp={handleAnchorMouseLeave}
            style={anchorContainerStyle}
          >
            <div style={anchorStyle} />
          </div>
        </div>
      </div>
    );
  }
}

export default compose(
  withToggleState({
    key: 'isAnchorHover',
    enable: 'handleAnchorMouseEnter',
    disable: 'handleAnchorMouseLeave',
  }),
  withToggleState({
    key: 'isDragging',
    enable: 'startDrag',
    disable: 'endDrag',
  }),
)(Slider);
