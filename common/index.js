const { readdir } = require('fs').promises

module.exports.wrapBody = (body) => `
  <!DOCTYPE html>
  <html>
    <head>
      <meta charset="utf-8">
      <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">
      <meta name="viewport" content="width=device-width, initial-scale=1">
      <style>
        /* https://stackoverflow.com/a/2810372/2923526 */
        img {
          height: 50%;
          width: auto;
        }
        .footer{ 
          position: fixed;     
          text-align: center;    
          bottom: 0px; 
          width: 100%;
        }
      </style>
    </head>
    <body>
      <div class="container">${body}</div>
      <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-ka7Sk0Gln4gmtz2MlQnikT1wXgYsOg+OMhuP+IlRH9sENBO0LRn5q+8nbTov4+1p" crossorigin="anonymous"></script>
    </body>
  </html>`;

module.exports.readDBs = async () => {
  const dirs = (await readdir('.', { withFileTypes: true })).filter((
    f
  ) => f.isDirectory() && ![
    '.git', 'common', 'node_modules', 'devops', 'src', 'public'
  ].includes(f.name)).map(f => f.name);
  let dirsWithDB = [];
  for (const d of dirs) {
    if ((await readdir(`./${d}`)).includes('db.js'))
      dirsWithDB.push(d)
  }
  return dirsWithDB.map(d => require(`../${d}/db.js`));
}

module.exports.env = {
  get connectionString() {
    const cs = process.env.CONNECTION_STRING;
    if (!cs) {
      throw 'missing connection string for db'
    }
    if(!cs.match(/^postgresql:\/\/[^:]+:[^@]+@[^:]+:[^\/]+$/)) {
      throw 'db\'s connection string does not match the expected format `postgresql://<host>:<password>@<server>:<port>`';
    }
    return cs;
  }
}
