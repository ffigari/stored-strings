const chokidar = require('chokidar')
const { spawn } = require('child_process')

const ts = () => (new Date).toISOString();

const spawnServer = (existingProcess) => {
  if (existingProcess) {
    existingProcess.kill('SIGHUP')
  }
  existingProcess = spawn('node', ['api.js'])
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

let webAPIProcess
chokidar.watch('./api.js').on('all', (event, path) => {
  webAPIProcess = spawnServer(webAPIProcess)
})
