import PropTypes from 'prop-types';
import React, { Fragment, PureComponent } from 'react';
import withHashState from 'client/app/react/containers/hoc/with-hash-state';
import CodeSnippetContainer from 'client/app/react/containers/search/results/code-snippet';
import KeyboardListener from 'client/app/react/components/passive/keyboard-listener';
import { modulo } from 'client/app/util/number';
import { Spacing } from 'react-elemental';
import { INPUT_FIELD_IDS } from 'client/app/util/constants/dom';
import { background } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

/**
 * List of results for matching code snippets.
 */
class CodeResultsContainer extends PureComponent {
  static propTypes = {
    repositories: PropTypes.object.isRequired,
    snippets: PropTypes.array.isRequired,
    onSearchQueryChange: PropTypes.func.isRequired,
    onSearchRepositoryFilePathChange: PropTypes.func.isRequired,
    // HOC props
    hash: PropTypes.number,
    setHash: PropTypes.func.isRequired,
  };

  static defaultProps = {
    hash: null,
  };

  componentDidMount() {
    const { hash } = this.props;

    if (hash) {
      this._scrollPermalink(hash);
    }
  }

  componentDidUpdate(prevProps) {
    if (this.props.hash !== prevProps.hash) {
      this._scrollPermalink(this.props.hash);
    }
  }

  navigateResults = this._navigateResults.bind(this);

  _scrollPermalink(idx) {  // eslint-disable-line class-methods-use-this
    const [node] = document.getElementsByName(`#R${idx}`);
    if (node) {
      node.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  }

  _navigateResults(delta) {
    const { snippets, hash, setHash } = this.props;

    return (evt) => {
      // Don't allow results navigation when focused on any input fields
      if (!INPUT_FIELD_IDS.includes(evt.target.id)) {
        if (hash === null) {
          setHash(0);
        } else {
          // Bidirectionally cycle through all available snippets
          setHash(modulo((hash) + delta, snippets.length));
        }
      }
    };
  }

  render() {
    const {
      snippets,
      repositories,
      hash,
      onSearchQueryChange,
      onSearchRepositoryFilePathChange,
    } = this.props;

    return (
      <Fragment>
        <KeyboardListener
          character="j"
          handler={this.navigateResults(1)}
        />
        <KeyboardListener
          character="k"
          handler={this.navigateResults(-1)}
        />

        {snippets.map((snippet, idx) => (
          <Spacing
            key={`${snippet.repo}-${snippet.path}`}
            name={`#R${idx}`}
            style={hash === idx ? {
              borderRadius: '3px',
              boxShadow: `0 0 0 1.5px ${background.primary.DELTA}`,
              transition: transition.all.DEFAULT,
            } : {}}
            bottom={idx < snippets.length - 1}
          >
            <CodeSnippetContainer
              position={idx + 1}
              permalink={`#R${idx}`}
              repositories={repositories}
              snippet={snippet}
              onSearchQueryChange={onSearchQueryChange}
              onSearchRepositoryFilePathChange={onSearchRepositoryFilePathChange}
              isHighlighted={hash === idx}
            />
          </Spacing>
        ))}
      </Fragment>
    );
  }
}

export default withHashState({
  serializer: (value) => `#R${value}`,
  deserializer: (hash) => parseInt(hash.substr(2), 10),
})(CodeResultsContainer);
