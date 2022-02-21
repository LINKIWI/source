module.exports = {
  collectCoverageFrom: ['<rootDir>/src/**/*.js'],
  setupFiles: ['./test/setup.js'],
  testEnvironment: './test/environment',
  testURL: 'http://localhost',
  verbose: true,
  moduleNameMapper: {
    '\\.(png)': 'identity-obj-proxy',
  },
  transformIgnorePatterns: [
    'node_modules/(?!(react-syntax-highlighter)/)',
  ],
};
