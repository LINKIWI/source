import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { Spacing } from 'react-elemental';
import { MdRoundedCorner } from 'react-icons/md';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import AdminControl from 'client/app/react/components/admin/control';
import AdminPanel from 'client/app/react/components/admin/panel';
import SelectList from 'client/app/react/components/ui/select-list';
import Separator from 'client/app/react/components/ui/separator';
import { setPreference } from 'client/app/redux/actions/preferences';
import { PREFERENCE_KEYS, PREFERENCE_VALUES } from 'client/app/util/constants/preferences';
import { objLookup } from 'shared/util/data';

class AdminPreferencesInterfaceContainer extends Component {
  static propTypes = {
    isCompact: PropTypes.bool.isRequired,
    // HOC props
    preferences: PropTypes.shape({
      [PREFERENCE_KEYS.CODE_SEARCH_RESULT_NAVIGATION]: PropTypes.string.isRequired,
      [PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME]: PropTypes.string.isRequired,
    }).isRequired,
    sourcePreviewOpts: PropTypes.object,
    actions: PropTypes.shape({
      setPreference: PropTypes.func.isRequired,
    }).isRequired,
  };

  static defaultProps = {
    sourcePreviewOpts: null,
  };

  handleChoiceChange = this._handleChoiceChange.bind(this);

  _handleChoiceChange(key) {
    return (value) => this.props.actions.setPreference(key, value);
  }

  render() {
    const { preferences, sourcePreviewOpts, isCompact } = this.props;

    return (
      <AdminPanel
        title="Interface"
        subtitle="User interface behavior preferences, persisted per-browser"
        icon={<MdRoundedCorner />}
      >
        <AdminControl
          title="Code search result navigation"
          description="Customize the behavior when clicking on a code line search result"
          isCompact={isCompact}
        >
          <SelectList
            value={preferences[PREFERENCE_KEYS.CODE_SEARCH_RESULT_NAVIGATION]}
            onChange={this.handleChoiceChange(PREFERENCE_KEYS.CODE_SEARCH_RESULT_NAVIGATION)}
            options={[
              {
                value: PREFERENCE_VALUES[PREFERENCE_KEYS.CODE_SEARCH_RESULT_NAVIGATION]
                  .SOURCE_LINK_NEW_TAB,
                label: 'Navigate to source link in a new tab',
              },
              {
                value: PREFERENCE_VALUES[PREFERENCE_KEYS.CODE_SEARCH_RESULT_NAVIGATION]
                  .SOURCE_LINK_SAME_TAB,
                label: 'Navigate to source link in the same tab',
              },
              sourcePreviewOpts && {
                value: PREFERENCE_VALUES[PREFERENCE_KEYS.CODE_SEARCH_RESULT_NAVIGATION]
                  .SOURCE_PREVIEW,
                label: 'Open interactive source preview',
              },
            ].filter(Boolean)}
          />
        </AdminControl>

        {sourcePreviewOpts && (
          <Fragment>
            <Spacing top bottom>
              <Separator />
            </Spacing>

            <AdminControl
              title="Syntax highlighter theme"
              description="Color scheme to use in interactive source previews"
              linkTitle="Browse Prism themes"
              linkHref="https://prismjs.com"
              isCompact={isCompact}
            >
              <SelectList
                value={preferences[PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME]}
                onChange={this.handleChoiceChange(PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME)}
                options={[
                  {
                    value: PREFERENCE_VALUES[PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME].COY,
                    label: 'Coy',
                  },
                  {
                    value: PREFERENCE_VALUES[PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME].MATERIAL_LIGHT,
                    label: 'Material Light',
                  },
                  {
                    value: PREFERENCE_VALUES[PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME]
                      .SOLARIZED_LIGHT,
                    label: 'Solarized Light',
                  },
                  {
                    value: PREFERENCE_VALUES[PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME].PRISM,
                    label: 'Prism',
                  },
                  {
                    value: PREFERENCE_VALUES[PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME].DUOTONE,
                    label: 'Duotone',
                  },
                  {
                    value: PREFERENCE_VALUES[PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME].GITHUB,
                    label: 'Github',
                  },
                ]}
                style={{ minWidth: '95px' }}
              />
            </AdminControl>
          </Fragment>
        )}
      </AdminPanel>
    );
  }
}

const mapStateToProps = ({ config, preferences }) => ({
  preferences,
  sourcePreviewOpts: objLookup(config, ['server', 'source']),
});

const mapDispatchToProps = (dispatch) => ({
  actions: bindActionCreators({ setPreference }, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(AdminPreferencesInterfaceContainer);
