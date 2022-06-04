module.exports.dbName = 'stored_strings';

module.exports.migrations = [{
  name: 'strings_table',
  up: async (client) => {
    await client.query(`
      create table strings (
        id    serial primary key,
        value text
      )
    `);
  },
  down: async (client) => {
    throw new Error('TODO: down migration');
  },
}];
