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

module.exports.migrateDB = async (
  connectionString, dbName, migrations, { direction, count }
) => {
  if (!['up', 'down'].includes(direction)) {
    throw new Error('incorrect direction');
  }
  if (count !== 'all' && isNaN(parseInt(count))) {
    throw new Error('invalid count');
  }
  if (direction === 'down') {
    throw new Error('down direction not implemented');
  }
  if (count !== 'all') {
    throw new Error('number count not implemented');
  }

  // TODO: Clients should be ended / released
  let dbClient;
  try {
    dbClient = await getDBClient(connectionString, dbName);
  } catch (e) {
    if (e.code === '3D000') {  // If DB does not exist
                               // https://www.postgresql.org/docs/current/errcodes-appendix.html
      const defaultDBClient = await getDefaultDBClient(connectionString);

      await defaultDBClient.query(`create database ${dbName}`);
      await defaultDBClient.end();
      dbClient = await getDBClient(connectionString, dbName);

      await dbClient.query(`
        create table migrations (
          name text not null check (name <> ''),
          applied_at timestamp default current_timestamp,
          application_order serial,
          unique (name)
        )
      `);
      await dbClient.query(`insert into migrations values ('creation')`);
    } else {
      throw e;
    }
  }
  const lastRunMigration = (await dbClient.query(`
    select name
    from migrations
    order by application_order desc
    limit 1
  `)).rows[0].name;

  let lastMigrationIdx;
  if (lastRunMigration === 'creation') {
    lastMigrationIdx = 0;
  } else {
    throw new Error('not implemented');
  }

  const migrationsToRun = migrations.slice(lastRunMigration)
  for (const m of migrationsToRun) {
    try {
      await dbClient.query('begin');
      await m.up(dbClient);
      await dbClient.query(`
        insert into migrations (name) values ('${m.name}')
      `);
      await dbClient.query('commit');
    } catch (e) {
      await dbClient.query('rollback');
      throw e;
    }
  }
  await dbClient.end()
  // TODO: Migrate up or down
  //   - run each m inside a tx
  //   - before running them, verify run migrations are consistent with input
  //     migrations; not sure which verifications have to be done
  //
  // How do I ensure migrations are reversible? A pre-commit hook which goes up
  // and down with a development db?
}

module.exports.ensureDBState = async (connectionString, dbName) => {
  const client = await getDBClient(connectionString, dbName);
  // TODO: Check migrations are up to date
  console.log(await client.query('select now()'))
  throw `'ensureDBState' not implemented`
}
