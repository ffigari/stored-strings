const { Client } = require('pg')

module.exports.putDBsUpToDate = () => {
  // TODO: Retrieve app folders in root dir which also have a db.js file
  //       Said file should contain the information about the name of the db and
  //       the migrations
  throw `'putDBsUpToDate' not implemented`
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
  throw 'not ready'
}
