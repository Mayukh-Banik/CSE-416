const webpack = require('webpack');
const { override } = require('customize-cra');

module.exports = function override(config, env) {
    // Only apply React Refresh in development mode
    // if (env === 'development') {
    //     config = addBabelPlugin('react-refresh/babel')(config);

        // Add polyfills for Webpack 5
        config.resolve.fallback = {
            ...config.resolve.fallback,
            stream: require.resolve('stream-browserify'),
            buffer: require.resolve('buffer/'),
            crypto: require.resolve('crypto-browserify'),
            assert: require.resolve('assert/'),
            http: require.resolve('stream-http'),
            https: require.resolve('https-browserify'),
            os: require.resolve('os-browserify/browser'),
            path: require.resolve('path-browserify'),
            process: require.resolve('process/browser'),
        };

        config.plugins.push(
            new webpack.ProvidePlugin({
                process: 'process/browser',
                Buffer: ['buffer', 'Buffer'],
            })
        );

        return config;
    }

    // Disable ESLint warnings/errors during compilation
    // const eslintPlugin = config.plugins.find(
    //     (plugin) => plugin.constructor && plugin.constructor.name === 'ESLintWebpackPlugin'
    // );
    // if (eslintPlugin) {
    //     eslintPlugin.options.failOnError = false;
    //     eslintPlugin.options.failOnWarning = false;
    // }

//     return config;
// };
