const webpack = require('webpack');

module.exports = function override(config) {
  // Find the ESLint plugin rule
  const eslintPlugin = config.plugins.find(
    (plugin) => plugin.constructor && plugin.constructor.name === 'ESLintWebpackPlugin'
  );

  if (eslintPlugin) {
    // Disable failing on warnings and errors
    eslintPlugin.options.failOnError = false;
    eslintPlugin.options.failOnWarning = false;
  }

  // 기존의 폴리필 설정 추가 (crypto-browserify, stream-browserify 등)
  config.resolve.fallback = {
    ...config.resolve.fallback,
    stream: require.resolve("stream-browserify"),
    buffer: require.resolve("buffer/"),
    crypto: require.resolve("crypto-browserify"),
    assert: require.resolve("assert/"),
    http: require.resolve("stream-http"),
    https: require.resolve("https-browserify"),
    os: require.resolve("os-browserify/browser"),
    path: require.resolve("path-browserify"),
    process: require.resolve("process/browser.js"),  // process 폴리필에 확장자 명시
  };

  // 필요한 경우, process와 Buffer 폴리필을 Webpack에 추가
  config.plugins.push(
    new webpack.ProvidePlugin({
      process: 'process/browser.js',  // 확장자 명시
      Buffer: ['buffer', 'Buffer'],
    })
  );

  return config;
};
