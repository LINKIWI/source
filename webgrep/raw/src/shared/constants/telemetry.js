// Enum of known client-side initiated actions.
export const TELEMETRY_ACTIONS = {
  // Navigation event to a new client-side defined route.
  RENDER_VIEW_ROUTE: 'RENDER_VIEW_ROUTE',
  // Click on a file path in a file search result.
  CLICK_FILE_RESULT_PATH: 'CLICK_FILE_RESULT_PATH',
  // Rank position of the clicked file result.
  CLICK_FILE_RESULT_POSITION: 'CLICK_FILE_RESULT_POSITION',
  // Click on a line in a code search result.
  CLICK_CODE_RESULT_LINE: 'CLICK_CODE_RESULT_LINE',
  // Click on the file path within a code search result.
  CLICK_CODE_RESULT_PATH: 'CLICK_CODE_RESULT_PATH',
  // Rank position of the clicked code result, for both lines and paths.
  CLICK_CODE_RESULT_POSITION: 'CLICK_CODE_RESULT_POSITION',
  // Open the source preview for a code search result.
  SOURCE_PREVIEW: 'SOURCE_PREVIEW',
  // Copy of text to the system clipboard.
  CLIPBOARD_WRITE: 'CLIPBOARD_WRITE',
  // Add an entry to client-side search history.
  RECORD_SEARCH_HISTORY_ITEM: 'RECORD_SEARCH_HISTORY_ITEM',
  // Search for a historical (recent) search query.
  EXECUTE_SEARCH_RECENT: 'EXECUTE_SEARCH_RECENT',
  // Search for a string from a source preview modal.
  EXECUTE_SEARCH_SOURCE_PREVIEW: 'EXECUTE_SEARCH_SOURCE_PREVIEW',
  // Set search filters to narrow results to a single file within a repository.
  EXECUTE_SEARCH_SINGLE_FILE: 'EXECUTE_SEARCH_SINGLE_FILE',
  // Incrementally increase maximum number of matches, in order to load more search results.
  LOAD_MORE_MATCH_RESULTS: 'LOAD_MORE_MATCH_RESULTS',
  // Request for current server metadata via client-side polling.
  POLL_SERVER_METADATA: 'POLL_SERVER_METADATA',
  // Update the current client-side metadata to be consistent with the server-side value.
  COMMIT_SERVER_METADATA: 'COMMIT_SERVER_METADATA',
  // Modification of a client-side preference.
  SET_PREFERENCE: 'SET_PREFERENCE',
  // Reset preferences to defaults.
  RESET_DEFAULT_PREFERENCES: 'RESET_DEFAULT_PREFERENCES',
  // Continuous duration that the window is active in the foreground.
  ACTIVE_SESSION_LENGTH: 'ACTIVE_SESSION_LENGTH',
  // Total duration that the window is available (including both in the foreground and background).
  TOTAL_SESSION_LENGTH: 'TOTAL_SESSION_LENGTH',
};

// Map of telemetry actions to known, whitelisted tag keys for that action type.
export const TELEMETRY_TAGS = {
  RENDER_VIEW_ROUTE: ['path'],
  CLICK_FILE_RESULT_PATH: ['repo'],
  CLICK_FILE_RESULT_POSITION: ['repo'],
  CLICK_CODE_RESULT_LINE: ['repo'],
  CLICK_CODE_RESULT_PATH: ['repo'],
  CLICK_CODE_RESULT_POSITION: ['repo'],
  EXECUTE_SEARCH_SINGLE_FILE: ['repo'],
  SET_PREFERENCE: ['key'],
};
