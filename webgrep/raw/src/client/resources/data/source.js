import { PREFERENCE_KEYS, PREFERENCE_VALUES } from 'client/app/util/constants/preferences';

/* eslint-disable max-len */

// Static mapping of known file extensions to corresponding language for syntax highlighting.
// Unknown file extensions (i.e. those that do not appear in this object) default to plain text.
// All available language constants are available here: https://github.com/react-syntax-highlighter/react-syntax-highlighter/blob/master/AVAILABLE_LANGUAGES_PRISM.MD
export const CANONICAL_FILE_EXTENSION_LANGUAGES = {
  '.babelrc': 'javascript',
  '.c': 'c',
  '.cc': 'cpp',
  '.cpp': 'cpp',
  '.css': 'css',
  Dockerfile: 'docker',
  '.ex': 'elixir',
  '.exs': 'elixir',
  '.gitignore': 'ignore',
  '.go': 'go',
  '.groovy': 'groovy',
  '.h': 'c',
  '.htm': 'markup',
  '.html': 'markup',
  '.ini': 'ini',
  '.java': 'java',
  Jenkinsfile: 'groovy',
  '.js': 'jsx',
  '.json': 'javascript',
  '.jsx': 'jsx',
  '.kt': 'kotlin',
  '.kts': 'kotlin',
  '.less': 'less',
  Makefile: 'makefile',
  '.md': 'markdown',
  '.php': 'php',
  '.pl': 'perl',
  '.pp': 'ruby',
  '.proto': 'protobuf',
  '.py': 'python',
  '.rb': 'ruby',
  '.rs': 'rust',
  '.sass': 'sass',
  '.scala': 'scala',
  '.scss': 'scss',
  '.sh': 'bash',
  '.sql': 'sql',
  '.svg': 'markup',
  '.swift': 'swift',
  '.toml': 'toml',
  '.ts': 'typescript',
  '.yaml': 'yaml',
  '.yml': 'yaml',
};

// Static mapping of syntax highlight theme preference values to the corresponding Prism style name.
// All available style names are available here: https://github.com/react-syntax-highlighter/react-syntax-highlighter/blob/master/AVAILABLE_STYLES_PRISM.MD
export const SYNTAX_HIGHLIGHT_THEME_NAMES = {
  [PREFERENCE_VALUES[PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME].COY]: 'coy',
  [PREFERENCE_VALUES[PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME].MATERIAL_LIGHT]: 'materialLight',
  [PREFERENCE_VALUES[PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME].SOLARIZED_LIGHT]: 'solarizedlight',
  [PREFERENCE_VALUES[PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME].PRISM]: 'prism',
  [PREFERENCE_VALUES[PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME].DUOTONE]: 'duotoneLight',
  [PREFERENCE_VALUES[PREFERENCE_KEYS.SYNTAX_HIGHLIGHT_THEME].GITHUB]: 'ghcolors',
};
