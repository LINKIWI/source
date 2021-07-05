import React from 'react';

/**
 * Read from the clipboard.
 *
 * @param {Function} cb Function invoked with (err, text) with clipboard contents.
 */
const read = (cb) => (
  window.navigator.clipboard ?
    window.navigator.clipboard.readText()
      .then((text) => cb(null, text))
      .catch((err) => cb(err)) :
    cb(new Error('clipboard unavailable'))
);

/**
 * Write to the clipboard.
 *
 * @param {String} text Text to write to the clipboard.
 * @param {Function} cb Function invoked with (err) on completion.
 */
const write = (text, cb) => (
  window.navigator.clipboard ?
    window.navigator.clipboard.writeText(text)
      .then(() => cb())
      .catch((err) => cb(err)) :
    cb(new Error('clipboard unavailable'))
);

/**
 * HOC factory for creating an HOC that injects methods to R/W to the system keyboard via native
 * DOM APIs.
 *
 * @param {Component|Function} WrappedComponent Underlying component.
 * @return {Function} Higher-order component with an additional `clipboard` prop injected.
 */
const withClipboard = (WrappedComponent) => (props) => (
  <WrappedComponent
    {...props}
    clipboard={{ read, write }}
  />
);

export default withClipboard;
