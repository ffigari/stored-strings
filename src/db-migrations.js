const { Client } = require('pg')

// TODO: This has to receive the connection string and the name of the db,
//       instead of assuming here anything about the structure of the repository
module.exports.putDBUpToDate = (dbName) => {
  console.log(`TODO: Create '${dbName}' db`)
  console.log(`TODO: Create migration\'s table at '${dbName}'`)
  console.log(`TODO: Run '${dbName}'\'s migrations`)
  throw `'putDBUpToDate' not implemented`
}

module.exports.ensureDBState = async (connectionString, dbName) => {
  let client;
  try {
    client = new Client({
      connectionString: `${connectionString}`,
    })
    await client.connect();
  } catch (e) {
    console.error(e)
    throw 'could not connect to PostgreSQL server';
  }
  try {
    client = new Client({
      connectionString: `${connectionString}/${dbName}`,
    });
    await client.connect();
  } catch (e) {
    console.error(e);
    throw `db '${dbName}' was not found`;
  }
  console.log(await client.query('select now()'))
  throw `'ensureDBState' not implemented`
}
