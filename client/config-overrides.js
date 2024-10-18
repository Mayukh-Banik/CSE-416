module.exports = function override(config) {
    const eslintRule = config.module.rules.find(
      (rule) => rule.use && rule.use.some((use) => use.loader && use.loader.includes('eslint-loader'))
    );
  
    if (eslintRule) {
      eslintRule.use = eslintRule.use.map((use) => {
        if (use.loader && use.loader.includes('eslint-loader')) {
          // Disable the ESLint loader
          return {
            ...use,
            options: {
              ...use.options,
              failOnError: false,
              failOnWarning: false,
            },
          };
        }
        return use;
      });
    }
  
    return config;
  };