const webpack = require('webpack');

module.exports = {
  resolve: {
    fallback: {
      "stream": require.resolve("stream-browserify"),  // Polyfill for stream
      "buffer": require.resolve("buffer/"),  // Polyfill for buffer
      "crypto": require.resolve("crypto-browserify"),  // If needed
      "assert": require.resolve("assert/"),  // If needed
      "http": require.resolve("stream-http"),  // If needed
      "https": require.resolve("https-browserify"),  // If needed
      "os": require.resolve("os-browserify/browser"),  // If needed
      "path": require.resolve("path-browserify"),  // If needed
    },
  },
  plugins: [
    new webpack.ProvidePlugin({
      process: 'process/browser',  // Polyfill for process
      Buffer: ['buffer', 'Buffer'],  // Polyfill for Buffer
    }),
  ],
};

