require('babel-register');
require('babel-polyfill');

module.exports = {
  migrations_directory: './migrations',
  solc: {
    optimizer: {
      enabled: true,
      runs: 200
    }
  }
};
