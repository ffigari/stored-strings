const express = require('express')

const { wrapBody } = require('./../common/index.js')

module.exports.addItselfTo = async (parentRouter, connectionString) => {
  const urlPath = '/desarollo';
  const appRouter = express.Router();

  appRouter.get('/', (req, res) => {
    res.send(wrapBody(`
      <div class="my-3">

        <h1>desarollo</h1>

        <h2>migraciones de bases de datos</h2>

        Una migración es una especificación de cómo debe modificarse el estado
        de una base de datos. Tal especificación es independiente de instancias
        concretas, las cuales pueden ser múltiples y en particular estar bajo
        uso activo de usuarios. La especificación debe incluir cómo aplicar los
        cambios deseados pero también cómo deshacerlos. Esto permite realizar
        un "rollback" para volver al estado anterior, en caso de que por ejemplo
        se inserte un bug crítico con el deploy asociado a tal migración.
        <br>
        En este sentido, a cada migración se la puede pensar como una operación
        <em>inversible</em> que toma una db y devuelve otra. Esta operación es a
        su vez una secuencia de suboperaciones atómicas. Ejemplos de estas
        suboperaciones son <span style="font-family: monospace">asd</span>

      </div>
    `));
  })

  parentRouter.use(urlPath, appRouter);
}
