import PropTypes from 'prop-types';
import React from 'react';
import { sizes, spacing } from 'react-elemental';
import createElement from 'react-syntax-highlighter/dist/esm/create-element';
import SyntaxHighlighter from 'react-syntax-highlighter/dist/esm/prism-async-light';
import * as styles from 'react-syntax-highlighter/dist/esm/styles/prism';
import { withToggleState } from '@linkiwi/hoc';
import { text } from 'client/app/util/style/color';

/**
 * Syntax-highlighted block of source code, with some styling defaults applied.
 */
const SourceCodeBlock = ({
  language,
  theme,
  focusLine,
  isFocused,
  setFocused,
  children,
}) => (
  <SyntaxHighlighter
    language={language}
    renderer={({ rows, stylesheet, useInlineStyles }) => rows.map((node, idx) => {
      // Following logic is borrowed from the react-syntax-highlighter default renderer
      const key = `code-segment${idx}`;
      const line = createElement({
        node,
        stylesheet,
        useInlineStyles,
        key,
      });

      // Take advantage of lazy ref setting to scroll the desired line into the viewport after it
      // has been mounted
      if (focusLine !== null && focusLine === idx + 1) {
        return (
          <div
            key={key}
            ref={(ref) => {
              // Only scroll into view once on first mount
              if (ref && !isFocused) {
                ref.scrollIntoView({ behavior: 'smooth', block: 'center' });
                setFocused();
              }
            }}
            style={{ backgroundColor: text.highlight.CONTEXT }}
          >
            {line}
          </div>
        );
      }

      return line;
    })}
    style={styles[theme]}
    customStyle={{
      backgroundColor: 'transparent',
      border: 0,
      fontSize: sizes.lambda,
      margin: 0,
      marginBottom: 0,
    }}
    codeTagProps={{
      style: {
        fontFamily: 'secondary--regular',
        textAlign: 'left',
        whiteSpace: 'pre',
        wordSpacing: 'normal',
        wordBreak: 'normal',
        overflowWrap: 'normal',
        lineHeight: 1.5,
        hyphens: 'none',
      },
    }}
    lineNumberContainerStyle={{
      float: 'left',
      paddingLeft: '1px',
      paddingRight: spacing.default,
      userSelect: 'none',
      MozUserSelect: 'none',
    }}
    lineNumberStyle={{
      color: text.dark.GAMMA,
    }}
    lineProps={{
      style: {
        display: 'block',
      },
    }}
    showLineNumbers
    wrapLines
  >
    {children}
  </SyntaxHighlighter>
);

SourceCodeBlock.propTypes = {
  language: PropTypes.string.isRequired,
  theme: PropTypes.string.isRequired,
  focusLine: PropTypes.number,
  children: PropTypes.node.isRequired,
  // HOC props
  setFocused: PropTypes.func.isRequired,
  isFocused: PropTypes.bool.isRequired,
};

SourceCodeBlock.defaultProps = {
  focusLine: null,
};

export default withToggleState({
  key: 'isFocused',
  enable: 'setFocused',
})(SourceCodeBlock);
