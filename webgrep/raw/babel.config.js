module.exports = {
  env: {
    development: {
      presets: [
        '@babel/preset-env',
        '@babel/preset-react',
      ],
      plugins: [
        ['@babel/plugin-proposal-decorators', { legacy: true }],
        '@babel/plugin-proposal-class-properties',
        ['module-resolver', { root: ['./src'] }],
        ['transform-inline-environment-variables', { include: ['VERSION'] }],
      ],
    },

    test: {
      presets: [
        '@babel/preset-env',
        '@babel/preset-react',
      ],
      plugins: [
        ['@babel/plugin-proposal-decorators', { legacy: true }],
        '@babel/plugin-proposal-class-properties',
        ['module-resolver', { root: ['./src'] }],
        ['transform-inline-environment-variables', { include: ['VERSION'] }],
      ],
    },

    production: {
      presets: [
        '@babel/preset-env',
        '@babel/preset-react',
      ],
      plugins: [
        ['@babel/plugin-proposal-decorators', { legacy: true }],
        '@babel/plugin-proposal-class-properties',
        'react-elemental',
        ['module-resolver', { root: ['./src'] }],
        ['transform-inline-environment-variables', { include: ['VERSION'] }],
        ['transform-react-remove-prop-types', { removeImports: true }],
      ],
    },
  },
};
