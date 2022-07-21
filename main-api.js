const express = require('express')
const bodyParser = require('body-parser')

const { wrapBody } = require('./common/index.js')
const storedStringsAPI = require('./stored-strings/api.js')
const recetarioAPI = require('./recetario/api.js')

const { ensureDBState } = require('./src/db-migrations.js')

const env = {
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

const main = async () => {
  const app = express()
  const port = 3000;

  // https://expressjs.com/en/starter/static-files.html
  app.use(express.static('public'))

  // parse application/x-www-form-urlencoded
  app.use(bodyParser.urlencoded({ extended: false }))
  app.get('/favicon.ico', (req, res) => {
    res.status(204).send()
  })

  // redirection middleware
  app.all('*', (req, res, next) => {
    if (req.url === '/') {
      return res.redirect('/i')
    }
    if (req.url.slice(0, 2) !== '/i') {
      return res.redirect('/i' + req.url)
    }
    next()
  })

  const mainRouter = express.Router()
  // TODO: Existing DBs should be read from the dirs
  // await ensureDBState(env.connectionString, 'stored_strings');
  // TODO: Existing APIs should be read from the dirs
  await storedStringsAPI.addItselfTo(mainRouter);
  await recetarioAPI.addItselfTo(mainRouter);
  mainRouter.get('/', (req, res) => {
    res.send(wrapBody(`
      <div class="my-3">
        <p>
          Qué tal? Mi nombre es Francisco Figari.
          Quisiera poder entender la programación como un alfarero entiende la
          arcilla y a la par de esa exloración retribuir a la sociedad que me
          formó.
        </p>

        <p>
          Busco dedicarme a entender el
          potencial de la programación en conectarnos con
          nuestra cotidianeidad.
        </p>

        <p>
          Ando un poco perdido respecto de cómo hacerlo.
        </p>

        <p>
          Ahora para fin de año tendría que estar recibido de computador
          científico de la UBA.
          Ando terminando mi tesis en la cual estudio la aplicabilidad de
          eye tracking web en análisis clínicos remotos por navegador.
          En <a href="https://github.com/ffigari/rastreador-ocular">este
          repo</a> anda lo que tengo armado hasta ahora.
          También armé <a href="/stored-strings/list">esta pavada</a> que lo
          único que hace es guardar strings hasta que reseteo el server.
        </p>

        <p>
          En cualquier momento mi mood está muy probablemente sintetizado por la
          siguiente imagen:
        </p>

        <img
          src="/index.jpg"
          alt="yo solo quiero codear lptm"
        >
      </div>
    `))
  })
  app.use('/i', mainRouter)
  app.listen(port, () => {
    console.log(`escuchando en puerto ${port}`);
  })
};

main().then().catch(e => {
  console.error(e)
  process.exit(-1)
});
