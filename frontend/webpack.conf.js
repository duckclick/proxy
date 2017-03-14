const path = require('path')
const webpack = require('webpack')

module.exports = {
  context: __dirname,
  resolve: {
    root: [
      path.join(__dirname, '/node_modules'),
      __dirname,
    ]
  },
  entry: 'src/index.js',
  output: {
    path: path.join(__dirname, '..', 'public'),
    filename: 'frame.js',
    sourceMapFilename: 'frame.js.map'
  },
  target: 'web',
  node: {
    process: false
  },
  plugins: [
    new webpack.DefinePlugin({
      BEAK_HOST: JSON.stringify(process.env.BEAK_HOST || 'http://localhost:7276')
    })
  ],
  devtool: 'inline-source-map',
  module: {
    loaders: [
      { test: /\.js$/, exclude: /node_modules/, loader: 'babel-loader' }
    ]
  }
}
