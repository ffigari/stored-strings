const express = require('express')
const bodyParser = require('body-parser')
const app = express()
const port = 3000;

// https://expressjs.com/en/starter/static-files.html
app.use(express.static('public'))

// parse application/x-www-form-urlencoded
app.use(bodyParser.urlencoded({ extended: false }))
app.get('/favicon.ico', (req, res) => {
  res.status(204).send()
})

const wrapBody = (body) => `
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

const l = [];
app.get('/informe', (req, res) => {
  res.send(wrapBody(`
    <div class="my-3">
      <p>
        Con <a href="/list">esta webapp</a> busco explorar la generación de recursos
        por parte de usuarios autenticados. Para esto cada usuario tendrá que
        poder crearse una cuenta, generar strings y luego verlos en algún lado.
        Se explorará tangencialmente la opción de generar tanto recursos
        públicos como privados.
      </p>
      <p>
        La implementación se realizará de manera iterativa.
        <br>
        Inicialmente se buscará únicamente la generación y persistencia de
        strings sin autenticación. Esto ya implicará la existencia de un
        servidor web que sirva una interfaz y la conecte con una base de datos.
        Tal interfaz contendrá un form para completar el campo necesario y
        aparte una lista donde ver los strings generados. En esta primer fase se
        evitará el uso de js en la ui.
        <br>
        En la primer fase se evitará el uso de js en la ui pero luego sí se
        agregará validación frontend a través de algún framework reactivo. En
        este paso se pretende explorar diferencias implementativas entre ambas
        opciones.
        <br>
        db; leer un poco distintos métodos para autenticar
        <br>
        público vs privado
      </p>
    </div>
  `))
})
app.get('/i/', (req, res) => {
  return res.redirect('/')
})
app.get('/', (req, res) => {
  res.send(wrapBody(`
    <div class="my-3">
      <img
        src="/index.jpg"
        alt="yo solo quiero codear lptm"
      >

      Ahora estoy terminando mi tesis de licenciatura en Ciencias de la
      Computación UBA. Se trata de eye tracking web en el contexto neuro y <a
      href="https://github.com/ffigari/rastreador-ocular">anda por acá</a> lo
      que fui armando.

      También armé <a href="/list">esta pavada</a> que lo único que hace es
      guardar strings hasta que reseteo el server.
      <br>
      Yo en verdad quiero buscar entender la programación como un alfarero
      entiende la arcilla. Y dsp al margén de eso preocuparme únicamente por el
      amor en sus diversas expresiones.
    </div>
  `))
})
app.get('/list', (req, res) => {
  // TODO: Retrieve strings from db
  res.send(wrapBody(`
    <div class="my-3">
      <h3>generación autenticada de strings</h1>
      <p>
        Created strings are shown below. You can always
        <a href="/create">create</a> more. Tmb podés leer un
        <a href="/informe">informe del trabajo realizado</a>
      </p>
      ${l.length === 0 ? `
        No string has yet been created.
      ` : `
        Textos creados:
        <ul>${l.map(x => `
          <li>${x}</li>
        `).join("")}</ul>
      `}
      <div class="footer"><a href="/">home</a></div>
    </div>

  `));
})

app.get('/create', (req, res) => {
  res.send(wrapBody(`
    <form class="my-3" method="POST">
      Acá podes guardar strings. Si no podés volver a ver <a href="/list">las
      que están creadas</a>.
      <div class="mb-3">
        <label for="important-text" class="form-label">El texto</label>
        <input
          type="text"
          class="form-control"
          id="important-text"
          name="foo-string"
          aria-describedby="input-help"
        >
        <div id="input-help" class="form-text">
          Este texto va a quedar bien guardado
        </div>
      </div>

      <input type="submit" class="btn btn-primary" value="Mandale">
    </form>
  `));
})

app.post('/create', (req, res) => {
  // TODO: Validar que no venga vacío el campo siguiente
  //       Qué se hace a nivel HTML?
  if (!req.body['foo-string']) {
    return res.status(400).send();
  }
  // TODO: Store this in database
  l.push(req.body['foo-string'])
  res.send(wrapBody(`
    <div class="my-3">
      <p>your string was correctly created! <a href="/list">check them all</a>
      or <a href="/create">create another one</a>.
    </div>
  `))
})

app.listen(port, () => {
  console.log(`escuchando en puerto ${port}`);
})
