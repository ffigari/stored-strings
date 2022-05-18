const express = require('express')
const { Client } = require('pg')

const { wrapBody } = require('./common/index.js')

module.exports.storedStrings = {
  addItselfTo: async (router, connectionString) => {
    const urlPath = '/stored-strings';
    const storedStringsRouter = express.Router()

    // TODO: Replace usage of this list below with adequate calls to the db
    const l = [];
    storedStringsRouter.get('/informe', (req, res) => {
      res.send(wrapBody(`
        <div class="my-3">
          <p>
            Con <a href="${urlPath}/list">esta webapp</a> busco explorar la
            generación de recursos por parte de usuarios autenticados. Para
            esto cada usuario tendrá que poder crearse una cuenta, generar
            strings y luego verlos en algún lado.  Se explorará tangencialmente
            la opción de generar tanto recursos públicos como privados.
          </p>
          <p>
            La implementación se realizará de manera iterativa.
            <br>
            Inicialmente se buscará únicamente la generación y persistencia de
            strings sin autenticación. Esto ya implicará la existencia de un
            servidor web que sirva una interfaz y la conecte con una base de
            datos.  Tal interfaz contendrá un form para completar el campo
            necesario y aparte una lista donde ver los strings generados. En
            esta primer fase se evitará el uso de js en la ui.
            <br>
            En la primer fase se evitará el uso de js en la ui pero luego sí se
            agregará validación frontend a través de algún framework reactivo.
            En este paso se pretende explorar diferencias implementativas entre
            ambas opciones.  <br> db; leer un poco distintos métodos para
            autenticar <br> público vs privado
          </p>
        </div>
      `))
    })
    storedStringsRouter.get('/list', (req, res) => {
      // TODO: Retrieve strings from db
      res.send(wrapBody(`
        <div class="my-3">
          <h3>generación autenticada de strings</h1>
          <p>
            Created strings are shown below. You can always
            <a href="${urlPath}/create">create</a> more. Tmb podés leer un
            <a href="${urlPath}/informe">informe del trabajo realizado</a>
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
    storedStringsRouter.get('/create', (req, res) => {
      res.send(wrapBody(`
        <form class="my-3" method="POST">
          Acá podes guardar strings. Si no podés volver a ver <a
          href="${urlPath}/list">las que están creadas</a>.
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

    storedStringsRouter.post('/create', (req, res) => {
      // TODO: Validar que no venga vacío el campo siguiente
      //       Qué se hace a nivel HTML?
      //       Va a convenir hacer validación pero tmb de frontend con algún
      //       framework reactivo
      if (!req.body['foo-string']) {
        return res.status(400).send();
      }
      // TODO: Store this in database
      l.push(req.body['foo-string'])
      res.send(wrapBody(`
        <div class="my-3">
          <p>your string was correctly created! <a href="${urlPath}/list">check
          them all</a> or <a href="${urlPath}/create">create another one</a>.
        </div>
      `))
    })
    router.use(urlPath, storedStringsRouter)
  }
}
