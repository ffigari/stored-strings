const express = require('express')
const bodyParser = require('body-parser')

const { wrapBody } = require('./common/index.js')
const { storedStrings } = require('./stored-strings')

// TODO: This import should not need to go back, 'src' should be added to PATH
const { ensureDBState } = require('./../src/dbMigrations.js')

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

  const content = express.Router()
  await ensureDBState(env.connectionString, 'stored_strings');
  throw 'foo';
  await storedStrings.addItselfTo(content);
  content.get('/', (req, res) => {
    res.send(wrapBody(`
      <div class="my-3">
        <p>
          Qué tal? Mi nombre es Francisco Figari. Busco entender la
          programación como un alfarero entiende la arcilla. A la par de tal
          exploración busco retribuir a la sociedad que me permitió formarme.
          Quisiera ver una Buenos Aires letrada en programación, tal que sus
          residentes puedan verse adecuadamente empoderados. No es evidente qué
          significan estos objetivos ni en qué dirección navegar para
          lograrlos.
        </p>

        <p>
          Ahora estoy terminando mi tesis de licenciatura en Ciencias de la
          Computación, UBA. Se trata de eye tracking web aplicado a la
          neurociencia y <a
          href="https://github.com/ffigari/rastreador-ocular">anda por acá</a>
          lo que fui armando.  También armé <a href="/stored-strings/list">esta
          pavada</a> que lo único que hace es guardar strings hasta que reseteo
          el server.
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
  app.use('/i', content)
  app.listen(port, () => {
    console.log(`escuchando en puerto ${port}`);
  })
};

main().then().catch(e => {
  console.error(e)
  process.exit(-1)
});
