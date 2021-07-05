/* eslint-disable import/no-extraneous-dependencies */

const CopyPlugin = require('copy-webpack-plugin');
const path = require('path');
const TerserPlugin = require('terser-webpack-plugin');
const webpack = require('webpack');

// Array of known environment variables whose value may be injected into the frontend as a key in
// process.env, allowing Node-like environment variable access in client-side logic.
const BUILD_ENV_VARS = [
  'NODE_ENV',
  'VERSION',
];

const isProduction = process.env.NODE_ENV === 'production';

module.exports = {
  mode: isProduction ? 'production' : 'development',
  entry: {
    main: './src/client',
  },
  output: {
    path: path.resolve(__dirname, './dist/client'),
    publicPath: '/assets/',
    filename: 'js/[name].js',
    chunkFilename: 'js/[name].js',
  },
  module: {
    rules: [
      {
        test: /src\/(client|shared)\/.+\.js$/,
        exclude: /node_modules/,
        loader: 'babel-loader',
      },
      {
        test: /src\/client\/.+\.(svg|png)$/,
        use: [
          {
            loader: 'url-loader',
            options: {},
          },
        ],
      },
      {
        test: /src\/client\/.+\.(woff|woff2)$/,
        loader: 'base64-inline-loader',
      },
    ],
  },
  optimization: {
    minimizer: [
      new TerserPlugin({
        parallel: true,
        extractComments: false,
        terserOptions: {
          ecma: 6,
          output: {
            comments: false,
          },
        },
      }),
    ],
  },
  performance: {
    hints: false,
  },
  plugins: [
    new webpack.ProgressPlugin(),
    new webpack.DefinePlugin({
      'process.env': BUILD_ENV_VARS
        .filter((key) => key in process.env)
        .reduce((acc, key) => ({
          ...acc,
          [key]: JSON.stringify(process.env[key]),
        }), {}),
    }),
    new CopyPlugin({
      patterns: [
        {
          from: 'src/client/resources/img/',
          to: 'img/',
        },
      ],
    }),
    isProduction && new webpack.LoaderOptionsPlugin({
      minimize: true,
    }),
  ].filter(Boolean),
  devServer: {
    historyApiFallback: true,
    publicPath: '/',
    port: 3000,
  },
};
