module.exports.dbName = 'stored_strings';

module.exports.migrations = [{
  name: 'strings_table',
  up: async (client) => {
    console.log('TODO: Implement UP migration')
  },
  down: async (client) => {
    // TODO: Implement down migration
  },
}];
