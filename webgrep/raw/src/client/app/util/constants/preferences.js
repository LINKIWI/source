// Enum of known client preference keys.
export const PREFERENCE_KEYS = {
  // Number of lines of surrounding context for each matching line in code search results.
  CODE_SEARCH_CONTEXT_LINES: 'CODE_SEARCH_CONTEXT_LINES',
  // Default number of matches to populate for each search query.
  CODE_SEARCH_INITIAL_MATCH_LIMIT: 'CODE_SEARCH_INITIAL_MATCH_LIMIT',
  // Default number of collapsed file result matches to display for each search query.
  CODE_SEARCH_FILE_RESULTS_LIMIT: 'CODE_SEARCH_FILE_RESULTS_LIMIT',
  // Behavior when clicking on a line in a code search result.
  CODE_SEARCH_RESULT_NAVIGATION: 'CODE_SEARCH_RESULT_NAVIGATION',
  // Source preview syntax highlighter theme.
  SYNTAX_HIGHLIGHT_THEME: 'SYNTAX_HIGHLIGHT_THEME',
};

// Mapping of known client preference keys to an enum of available values.
export const PREFERENCE_VALUES = {
  [PREFERENCE_KEYS.CODE_SEARCH_RESULT_NAVIGATION]: {
    SOURCE_LINK_NEW_TAB: 'SOURCE_LINK_NEW_TAB',
    SOURCE_LINK_SAME_TAB: 'SOURCE_LINK_SAME_TAB',
    SOURCE_PREVIEW: 'SOURCE_PREVIEW',
  },
  [PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME]: {
    COY: 'COY',
    MATERIAL_LIGHT: 'MATERIAL_LIGHT',
    SOLARIZED_LIGHT: 'SOLARIZED_LIGHT',
    PRISM: 'PRISM',
    DUOTONE: 'DUOTONE',
    GITHUB: 'GITHUB',
  },
};
