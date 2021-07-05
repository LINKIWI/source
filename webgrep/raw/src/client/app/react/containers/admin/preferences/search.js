import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { Spacing } from 'react-elemental';
import { MdSearch } from 'react-icons/md';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { compose, withForm } from '@linkiwi/hoc';
import AdminControl from 'client/app/react/components/admin/control';
import AdminPanel from 'client/app/react/components/admin/panel';
import Separator from 'client/app/react/components/ui/separator';
import Slider from 'client/app/react/components/ui/slider';
import { setPreference } from 'client/app/redux/actions/preferences';
import { PREFERENCE_KEYS } from 'client/app/util/constants/preferences';

// Upper limit for the number of context lines that can be set with the interactive slider.
const MAX_CONTEXT_LINES = 20;

// Upper limit for the match limit that can be set with the interactive slider.
const MAX_INITIAL_MATCH_LIMIT = 200;

// Upper limit for the collapsed file results that can be set with the interactive slider.
const MAX_FILE_RESULTS_LIMIT = 50;

class AdminPreferencesSearchContainer extends Component {
  static propTypes = {
    // HOC props
    form: PropTypes.shape({
      contextSlider: PropTypes.number.isRequired,
      matchLimitSlider: PropTypes.number.isRequired,
      fileResultsLimitSlider: PropTypes.number.isRequired,
    }).isRequired,
    handleFormChange: PropTypes.func.isRequired,
    preferences: PropTypes.shape({
      [PREFERENCE_KEYS.CODE_SEARCH_CONTEXT_LINES]: PropTypes.number.isRequired,
    }).isRequired,
    actions: PropTypes.shape({
      setPreference: PropTypes.func.isRequired,
    }).isRequired,
  };

  componentDidUpdate(prevProps) {
    const { preferences, handleFormChange } = this.props;
    const { preferences: prevPreferences } = prevProps;

    const contextLines = preferences[PREFERENCE_KEYS.CODE_SEARCH_CONTEXT_LINES];
    const prevContextLines = prevPreferences[PREFERENCE_KEYS.CODE_SEARCH_CONTEXT_LINES];
    const matchLimit = preferences[PREFERENCE_KEYS.CODE_SEARCH_INITIAL_MATCH_LIMIT];
    const prevMatchLimit = prevPreferences[PREFERENCE_KEYS.CODE_SEARCH_INITIAL_MATCH_LIMIT];
    const fileResultsLimit = preferences[PREFERENCE_KEYS.CODE_SEARCH_FILE_RESULTS_LIMIT];
    const prevFileResultsLimit = prevPreferences[PREFERENCE_KEYS.CODE_SEARCH_FILE_RESULTS_LIMIT];

    // Redux-induced change in number of slider values; propagate this change as an adjustment to
    // slider position
    if (prevContextLines !== contextLines) {
      handleFormChange('contextSlider')(contextLines / MAX_CONTEXT_LINES);
    }
    if (matchLimit !== prevMatchLimit) {
      handleFormChange('matchLimitSlider')(matchLimit / MAX_INITIAL_MATCH_LIMIT);
    }
    if (fileResultsLimit !== prevFileResultsLimit) {
      handleFormChange('fileResultsLimitSlider')(fileResultsLimit / MAX_FILE_RESULTS_LIMIT);
    }
  }

  handleChoiceChange = this._handleChoiceChange.bind(this);

  _handleChoiceChange(key) {
    return (value) => this.props.actions.setPreference(key, value);
  }

  render() {
    const {
      form: { contextSlider, matchLimitSlider, fileResultsLimitSlider },
      handleFormChange,
    } = this.props;
    const {
      CODE_SEARCH_CONTEXT_LINES,
      CODE_SEARCH_INITIAL_MATCH_LIMIT,
      CODE_SEARCH_FILE_RESULTS_LIMIT,
    } = PREFERENCE_KEYS;

    const contextLines = Math.round(contextSlider * MAX_CONTEXT_LINES);
    const matchLimit = Math.round(matchLimitSlider * MAX_INITIAL_MATCH_LIMIT);
    const fileResultsLimit = Math.round(fileResultsLimitSlider * MAX_FILE_RESULTS_LIMIT);

    return (
      <AdminPanel
        title="Search"
        subtitle="Code search preferences, persisted per-browser"
        icon={<MdSearch />}
      >
        <AdminControl
          title="Code context lines"
          description="Number of lines of surrounding context for each matching search result"
        >
          <Slider
            value={contextSlider}
            annotation={contextLines.toString()}
            onChange={handleFormChange('contextSlider')}
            onFinalize={() =>
              this.handleChoiceChange(CODE_SEARCH_CONTEXT_LINES)(contextLines)}
          />
        </AdminControl>

        <Spacing top bottom>
          <Separator />
        </Spacing>

        <AdminControl
          title="Default match limit"
          description="Initial maximum number of matched lines per search query"
        >
          <Slider
            value={matchLimitSlider}
            annotation={matchLimit.toString()}
            onChange={handleFormChange('matchLimitSlider')}
            onFinalize={() =>
              this.handleChoiceChange(CODE_SEARCH_INITIAL_MATCH_LIMIT)(matchLimit)}
          />
        </AdminControl>

        <Spacing top bottom>
          <Separator />
        </Spacing>

        <AdminControl
          title="Collapsed file results limit"
          description="Number of file result matches to display per search query before expansion"
        >
          <Slider
            value={fileResultsLimitSlider}
            annotation={fileResultsLimit.toString()}
            onChange={handleFormChange('fileResultsLimitSlider')}
            onFinalize={() =>
              this.handleChoiceChange(CODE_SEARCH_FILE_RESULTS_LIMIT)(fileResultsLimit)}
          />
        </AdminControl>
      </AdminPanel>
    );
  }
}

const mapStateToProps = ({ preferences }) => ({ preferences });

const mapDispatchToProps = (dispatch) => ({
  actions: bindActionCreators({ setPreference }, dispatch),
});

export default compose(
  connect(mapStateToProps, mapDispatchToProps),
  withForm({
    initial: ({ preferences }) => ({
      contextSlider: preferences[PREFERENCE_KEYS.CODE_SEARCH_CONTEXT_LINES] / MAX_CONTEXT_LINES,
      matchLimitSlider:
        preferences[PREFERENCE_KEYS.CODE_SEARCH_INITIAL_MATCH_LIMIT] / MAX_INITIAL_MATCH_LIMIT,
      fileResultsLimitSlider:
        preferences[PREFERENCE_KEYS.CODE_SEARCH_FILE_RESULTS_LIMIT] / MAX_FILE_RESULTS_LIMIT,
    }),
  }),
)(AdminPreferencesSearchContainer);
