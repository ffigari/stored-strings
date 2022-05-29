const { Client } = require('pg')

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
    // If the db does not exist..
    // (https://www.postgresql.org/docs/current/errcodes-appendix.html)
    if (e.code = '3D000') {
      // ..create it
      const defaultDBClient = await getDefaultDBClient(connectionString);
      await defaultDBClient.query(`create database ${dbName}`);
      dbClient = await getDBClient(connectionString, dbName);
    } else {
      throw e;
    }
  }
  console.log((await dbClient.query('select now()')).rows)
  console.log(`TODO: Create migration\'s table at '${dbName}'`)
  console.log(`TODO: Run '${dbName}'\'s migrations`)
  throw `'migrateDB' not implemented`
}

module.exports.ensureDBState = async (connectionString, dbName) => {
  const client = await getDBClient(connectionString, dbName);
  console.log(await client.query('select now()'))
  throw `'ensureDBState' not implemented`
}
