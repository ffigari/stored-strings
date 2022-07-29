const chokidar = require('chokidar')
const { spawn } = require('child_process')

const { putDBsUpToDate } = require('./src/db-migrations.js');

const { printCV } = require('./cv/index.js');

const ts = () => (new Date).toISOString();

// ansi colors: https://stackoverflow.com/a/41407246/2923526
const RED = '\x1b[31m';
const RESET = '\x1b[0m';

const spawnServer = (existingProcess) => {
  if (existingProcess) {
    existingProcess.kill('SIGHUP')
  }
  try {
    existingProcess = spawn('node', ['main-api.js'])
  } catch (e) {
    console.error(e)
    throw e
  }
  console.log(`[runner:${ts()}] new server spawned`)
  existingProcess.stdout.on('data', (data) => {
    data = data.toString().replace(/\n$/, "");
    console.log(`[api:${ts()}] ${data}`);
  });

  existingProcess.stderr.on('data', (data) => {
    data = data.toString().replace(/\n$/, "");
    console.error(`${RED}[api:${ts()}] ${data}${RESET}`);
  });
  return existingProcess
}

const main = () => {
  const action = process.argv[2];
  if (!action) {
    throw 'action was not indicated'
  }
  // TODO: Este archivo no debería ser exclusivo a la main api
  if (action === "print-cv") {
    printCV();
  } else if (action === "watch") {
    let webAPIProcess
    // TODO: Review why this spawn lots of servers on startup
    chokidar.watch('.').on('all', (event, path) => {
      webAPIProcess = spawnServer(webAPIProcess)
    })
  } else if (action === "spawn") {
    spawnServer();
  } else if (action === "migrate") {
    // TODO: here the existing dbs should also be found automatically
    putDBsUpToDate();
  } else {
    throw `action ${action} was not recognized`;
  }
}

try {
  main();
} catch (e) {
  console.error(`${RED}[runner:${ts()}] ${e}${RESET}`);
  process.exit(-1);
}
