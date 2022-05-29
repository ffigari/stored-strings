const { Client } = require('pg')

const MIGRATIONS_TABLE_NAME = 'migrations';

const getDefaultDBClient = async (connectionString) => {
  const client = new Client({
    connectionString: `${connectionString}`,
  })
  await client.connect();
  return client;
}

const getDBClient = async (connectionString, dbName) => {
  const client = new Client({
    connectionString: `${connectionString}/${dbName}`,
  });
  await client.connect();
  return client;
}

module.exports.migrateDB = async (connectionString, dbName) => {
  // TODO: Migrations should consider the possibility of the db being already
  //       migrated
  let dbClient;
  try {
    dbClient = await getDBClient(connectionString, dbName);
  } catch (e) {
    if (e.code === '3D000') {  // If DB does not exist
                               // https://www.postgresql.org/docs/current/errcodes-appendix.html
      const defaultDBClient = await getDefaultDBClient(connectionString);
      await defaultDBClient.query(`create database ${dbName}`);
      dbClient = await getDBClient(connectionString, dbName);
    } else {
      throw e;
    }
  }
  console.log(`TODO: Obtain last run migration`)
  let lastRunMigration;
  try {
    console.log(await dbClient.query(`select * from ${MIGRATIONS_TABLE_NAME}`));
    throw 'foo'
  } catch (e) {
    // TODO: Maybe this should be done along with the db creation
    if (e.code === '42P01') {  // If table does not exist
      await dbClient.query(`
        create table ${MIGRATIONS_TABLE_NAME} (
          name text not null check (name <> ''),
          applied_at timestamp default current_timestamp,
          application_order serial,
          unique (name)
        )
      `)
      const dbCreationMigrationName = 'creation';
      await dbClient.query(`
        insert into ${MIGRATIONS_TABLE_NAME}
        values ('${dbCreationMigrationName}')
      `)
      lastRunMigration = dbCreationMigrationName;
    } else {
      throw e
    }
  }
  // Migrations should be run inside a transaction
  console.log(`TODO: Verify migrations status, last migration run=${lastRunMigration}`)
  console.log(`TODO: Run '${dbName}'\'s migrations`)
  throw `'migrateDB' not implemented`
}

module.exports.ensureDBState = async (connectionString, dbName) => {
  const client = await getDBClient(connectionString, dbName);
  // Maybe the best here is to simply check that the last run migration
  // coincides with the last one created
  // But how are migration identified: an id? how do I choose it? manually? What
  // about git merges? Maybe just a name and doing a check that names do not
  // repeat? When, before running migrations? With a pre-commit hook?
  console.log(await client.query('select now()'))
  throw `'ensureDBState' not implemented`
}
