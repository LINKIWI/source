import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { Spacing } from 'react-elemental';
import { Helmet } from 'react-helmet';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { compose, withForm } from '@linkiwi/hoc';
import withTransactionalTunnel from 'client/app/react/containers/hoc/with-transactional-tunnel';
import SearchResultsContainer from 'client/app/react/containers/search/results';
import SearchQueryContainer from 'client/app/react/containers/search/query';
import Layout from 'client/app/react/components/common/layout';
import { CONNECTIVITY_STATE } from 'client/app/react/components/search/meta/connection-status';
import { recordSearchHistoryItem } from 'client/app/redux/actions/search';
import { PREFERENCE_KEYS } from 'client/app/util/constants/preferences';
import { MAX_MATCHES } from 'client/app/util/constants/search';
import { string } from 'client/app/util/format';
import { URLStateSerializer, URLStateDeserializer } from 'client/app/util/data';
import { decodeURLState, encodeURLState } from 'client/app/util/navigation';
import { objLookup } from 'shared/util/data';

// Window widths at which to consider the layout to be horizontally compact.
const breakpoints = [
  { threshold: 700, compact: true },
];

/**
 * Container for managing query state and mediating search requests over the network.
 */
class SearchContainer extends Component {
  static propTypes = {
    // HOC props
    title: PropTypes.string.isRequired,
    width: PropTypes.number,
    contextLines: PropTypes.number.isRequired,
    initialMatchLimit: PropTypes.number.isRequired,
    indexIdentity: PropTypes.shape({
      name: PropTypes.string.isRequired,
      timestamp: PropTypes.number.isRequired,
    }).isRequired,
    actions: PropTypes.shape({
      recordSearchHistoryItem: PropTypes.func.isRequired,
    }).isRequired,
    results: PropTypes.shape({
      err: PropTypes.object,
      message: PropTypes.object,
      isConnected: PropTypes.bool.isRequired,
      sendMessage: PropTypes.func.isRequired,
    }).isRequired,
    form: PropTypes.shape({
      query: PropTypes.string.isRequired,
      regex: PropTypes.bool.isRequired,
      maxMatches: PropTypes.number.isRequired,
      caseSensitive: PropTypes.bool.isRequired,
      file: PropTypes.string.isRequired,
      repos: PropTypes.arrayOf(PropTypes.string.isRequired).isRequired,
    }).isRequired,
    handleFormChange: PropTypes.func.isRequired,
  };

  static defaultProps = {
    width: null,
  };

  componentDidUpdate(prevProps) {
    // Trigger another search request if the connection state changed from disconnected to connected
    if (this.props.results.isConnected && !prevProps.results.isConnected) {
      this.invoke();
    }
  }

  handleQueryChange = this._handleQueryChange.bind(this);

  handleRegexChange = this._handleRegexChange.bind(this);

  handleCaseSensitivityChange = this._handleCaseSensitivityChange.bind(this);

  handleMaxMatchesChange = this._handleMaxMatchesChange.bind(this);

  handleFilePathChange = this._handleFilePathChange.bind(this);

  handleRepositoriesChange = this._handleRepositoriesChange.bind(this);

  invoke = this._invoke.bind(this);

  _handleQueryChange(evt, cb) {
    const { form, initialMatchLimit, handleFormChange } = this.props;

    // Changes to the search query should reset the number of max matches when it has only seen
    // incremental changes (i.e. loading more search results). A special case is when the max
    // matches limit is set to "all matches," in which case it should not be altered.

    if (form.maxMatches === MAX_MATCHES.UNLIMITED) {
      handleFormChange('query')(evt);
      this.invoke({ query: evt.target.value }, cb);
    } else {
      this.props.handleFormChange('query')(evt);
      this.props.handleFormChange('maxMatches')(initialMatchLimit);
      this.invoke({ query: evt.target.value, maxMatches: initialMatchLimit }, cb);
    }
  }

  _handleRegexChange(state, cb) {
    this.props.handleFormChange('regex')(state);
    this.invoke({ regex: state }, cb);
  }

  _handleCaseSensitivityChange(state, cb) {
    this.props.handleFormChange('caseSensitive')(state);
    this.invoke({ caseSensitive: state }, cb);
  }

  _handleMaxMatchesChange(value, cb) {
    const { initialMatchLimit, handleFormChange } = this.props;

    // null max matches implies resetting to the default/initial match limit, as supplied in
    // preferences.

    if (value !== null) {
      handleFormChange('maxMatches')(value);
      this.invoke({ maxMatches: value }, cb);
    } else {
      handleFormChange('maxMatches')(initialMatchLimit);
      this.invoke({ maxMatches: initialMatchLimit }, cb);
    }
  }

  _handleFilePathChange(evt, cb) {
    this.props.handleFormChange('file')(evt);
    this.invoke({ file: evt.target.value }, cb);
  }

  _handleRepositoriesChange(repos, cb) {
    this.props.handleFormChange('repos')(repos);
    this.invoke({ repos }, cb);
  }

  _invoke(overrides = {}, cb) {
    const { results, form, indexIdentity, contextLines, initialMatchLimit } = this.props;
    const searchParams = { ...form, ...overrides };

    // Only re-encode URL state when a form parameter changes
    if (Object.keys(overrides).length) {
      encodeURLState(
        searchParams.query && searchParams,
        // Conditionally serialize request parameters into URL state, to keep the URL shorter
        {
          regex: (val) => val,
          caseSensitive: (val) => val,
          repos: (val) => val.length > 0,
          maxMatches: (val) => val !== initialMatchLimit,
          file: (val) => val,
        },
        // Special case serialization for non-primitive field values
        {
          repos: URLStateSerializer.array,
          query: URLStateSerializer.string,
          file: URLStateSerializer.string,
        },
      );
    }

    const request = {
      ...searchParams,
      // The index identity is used to "fingerprint" the client's search request context.
      // Any discrepancy between the client-asserted and server-asserted index identities implies
      // that the search index has changed since the client last updated its knowledge of the
      // server-side index metadata.
      indexIdentity,
      // Customize the number of context lines based on the global preference value.
      context: contextLines,
    };

    results.sendMessage(request);

    if (cb) {
      this.forceUpdate(cb);
    }
  }

  render() {
    const { results, form, title, width, actions } = this.props;
    const searchTitle = form.query ? `${form.query} - ${title}` : title;
    const isCompact = breakpoints.reduce((acc, breakpoint) =>
      (width < breakpoint.threshold ? breakpoint.compact : acc), false);
    const connectivity = (() => {
      if (results.err) {
        return CONNECTIVITY_STATE.DISCONNECTED;
      }

      if (results.isConnected) {
        return CONNECTIVITY_STATE.CONNECTED;
      }

      return CONNECTIVITY_STATE.CONNECTING;
    })();

    return (
      <Fragment>
        <Helmet>
          <title>
            {searchTitle}
          </title>
        </Helmet>

        <Layout width={width}>
          <Spacing size="huge" bottom>
            <Spacing size="large" top>
              <SearchQueryContainer
                value={form.query}
                regex={form.regex}
                caseSensitive={form.caseSensitive}
                maxMatches={form.maxMatches}
                filePath={form.file}
                filteredRepos={form.repos}
                results={results.message}
                onQueryChange={this.handleQueryChange}
                onRegexChange={this.handleRegexChange}
                onCaseSensitivityChange={this.handleCaseSensitivityChange}
                onMaxMatchesChange={this.handleMaxMatchesChange}
                onFilePathChange={this.handleFilePathChange}
                onRepositoriesChange={this.handleRepositoriesChange}
                onQueryRecord={actions.recordSearchHistoryItem}
                connectivity={connectivity}
                isCompact={isCompact}
              />
            </Spacing>
          </Spacing>

          <Spacing size="huge" bottom>
            <SearchResultsContainer
              maxMatches={form.maxMatches}
              results={results.message}
              onQueryChange={this.handleQueryChange}
              onQueryRecord={actions.recordSearchHistoryItem}
              onMaxMatchesChange={this.handleMaxMatchesChange}
              onFilePathChange={this.handleFilePathChange}
              onRepositoriesChange={this.handleRepositoriesChange}
            />
          </Spacing>
        </Layout>
      </Fragment>
    );
  }
}

const mapStateToProps = ({ config, context, meta, preferences }) => ({
  title: string(objLookup(config, ['client', 'site', 'title']) || 'webgrep', {
    page: 'Search',
    name: meta.name,
  }),
  width: context.window.width,
  indexIdentity: { name: meta.name, timestamp: meta.timestamp },
  contextLines: preferences[PREFERENCE_KEYS.CODE_SEARCH_CONTEXT_LINES],
  initialMatchLimit: preferences[PREFERENCE_KEYS.CODE_SEARCH_INITIAL_MATCH_LIMIT],
});

const mapDispatchToProps = (dispatch) => ({
  actions: bindActionCreators({ recordSearchHistoryItem }, dispatch),
});

export default compose(
  connect(mapStateToProps, mapDispatchToProps),
  withForm({
    initial: ({ initialMatchLimit }) => decodeURLState({
      regex: false,
      caseSensitive: false,
      maxMatches: initialMatchLimit,
      file: '',
      repos: [],
      query: '',
    }, {
      regex: URLStateDeserializer.boolean,
      caseSensitive: URLStateDeserializer.boolean,
      maxMatches: URLStateDeserializer.number,
      file: URLStateDeserializer.string,
      repos: URLStateDeserializer.array,
      query: URLStateDeserializer.string,
    }),
  }),
  withTransactionalTunnel({
    key: 'results',
    endpoint: '/api/search',
    messageBufferSize: 1,
  }),
)(SearchContainer);
