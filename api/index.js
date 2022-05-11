const express = require('express')
const bodyParser = require('body-parser')

const { wrapBody } = require('./common/index.js')
const { storedStrings } = require('./stored-strings')

const app = express()
const port = 3000;

// https://expressjs.com/en/starter/static-files.html
app.use(express.static('public'))

// parse application/x-www-form-urlencoded
app.use(bodyParser.urlencoded({ extended: false }))
app.get('/favicon.ico', (req, res) => {
  res.status(204).send()
})

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
storedStrings.addItselfTo(content)
content.get('/', (req, res) => {
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

      También armé <a href="/stored-strings/list">esta pavada</a> que lo único
      que hace es guardar strings hasta que reseteo el server.
      <br>
      Yo en verdad quiero buscar entender la programación como un alfarero
      entiende la arcilla. Y dsp al margén de eso preocuparme únicamente por el
      amor en sus diversas expresiones.
    </div>
  `))
})
app.use('/i', content)
app.listen(port, () => {
  console.log(`escuchando en puerto ${port}`);
})
