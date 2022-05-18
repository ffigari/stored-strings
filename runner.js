const chokidar = require('chokidar')
const { spawn } = require('child_process')

const ts = () => (new Date).toISOString();

// ansi colors: https://stackoverflow.com/a/41407246/2923526
const RED = '\x1b[31m';
const RESET = '\x1b[0m';

const spawnServer = (existingProcess) => {
  if (existingProcess) {
    existingProcess.kill('SIGHUP')
  }
  try {
    existingProcess = spawn('node', ['api/index.js'])
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


if (process.argv[2] === "watch") {
  let webAPIProcess
  chokidar.watch('./api').on('all', (event, path) => {
    webAPIProcess = spawnServer(webAPIProcess)
  })
} else {
  spawnServer();
}
