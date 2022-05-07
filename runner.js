const chokidar = require('chokidar')
const { spawn } = require('child_process')

const ts = () => (new Date).toISOString();

const spawnServer = (existingProcess) => {
  if (existingProcess) {
    existingProcess.kill('SIGHUP')
  }
  try {
    existingProcess = spawn('node', ['api.js'])
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
    console.error(`[api:${ts()}] ${data}`);
  });
  return existingProcess
}


if (process.argv[2] === "watch") {
  let webAPIProcess
  chokidar.watch('./api.js').on('all', (event, path) => {
    webAPIProcess = spawnServer(webAPIProcess)
  })
} else {
  spawnServer();
}
