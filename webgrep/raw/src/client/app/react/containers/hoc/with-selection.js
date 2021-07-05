import React, { Component } from 'react';

/**
 * HOC factory for a higher-order component that tracks DOM selection changes.
 *
 * @param {Function|Component} WrappedComponent Underlying component
 * @returns {Function} Factory that accepts a component as input and returns an HOC with selection
 *                     changes built into its lifecycle methods. Selection properties are injected
 *                     into the underlying component in prop `selection`.
 */
const withSelection = (WrappedComponent) =>
  class WithSelectionHOC extends Component {
    state = {
      selection: {
        anchor: null,
        text: '',
        bounds: {
          x: 0,
          y: 0,
          width: 0,
          height: 0,
          top: 0,
          right: 0,
          bottom: 0,
          left: 0,
        },
      },
    };

    componentDidMount() {
      document.addEventListener('selectionchange', this.handleSelectionChange);
    }

    componentWillUnmount() {
      document.removeEventListener('selectionchange', this.handleSelectionChange);
    }

    handleSelectionChange = this._handleSelectionChange.bind(this);

    _handleSelectionChange() {
      return this.setState((state) => {
        const selection = document.getSelection();
        const bounds = (selection && selection.rangeCount > 0 && selection.toString()) ?
          selection.getRangeAt(0).getBoundingClientRect() :
          state.selection.bounds;

        return {
          selection: {
            anchor: selection.anchorNode,
            text: selection.toString(),
            bounds,
          },
        };
      });
    }

    render() {
      const props = {
        ...this.props,
        ...this.state,
      };

      return (
        <WrappedComponent {...props} />
      );
    }
  };

export default withSelection;
