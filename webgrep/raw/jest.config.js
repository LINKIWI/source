module.exports = {
  collectCoverageFrom: ['<rootDir>/src/**/*.js'],
  setupFiles: ['./test/setup.js'],
  testEnvironment: './test/environment',
  testURL: 'http://localhost',
  verbose: true,
  moduleNameMapper: {
    '\\.(png|woff|woff2)': 'identity-obj-proxy',
  },
  transformIgnorePatterns: [
    'node_modules/(?!(react-syntax-highlighter)/)',
  ],
};
