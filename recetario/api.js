const express = require('express')

const { wrapBody } = require('./../common/index.js')

module.exports.addItselfTo = async (parentRouter, connectionString) => {
  const urlPath = '/recetario';
  const appRouter = express.Router();

  appRouter.get('/brownie', (req, res) => {
    res.send(wrapBody(`
      <div class="my-3">

        <h1>brownie</h1>

        <h2>ingredientes</h2>
        <ul>
          <li>100 grs de chocolate</li>
          <li>50 grs de manteca + manteca para enmantecar el molde</li>
          <li>un chorrito de aceite</li>
          <li>media taza de azucar moreno</li>
          <li>un huevo</li>
          <li>una cucharadita de esencia artificial de vainilla</li>
          <li>una pizca de sal</li>
          <li>media taza de harina 0000 no leudante</li>
          <li>una cucharada de cacao en polvo</li>
          <li>un puñado de nueces</li>
        </ul>

        <h2>preparación</h2>
        <ol>
          <li>Trozear el chocolate. Derretirlo a baño maría junto a la manteca,
          agregando tmb un chorrito de aceite y de agua.</li>
          <li>A la mezcla anterior agregar la taza de azucar, el huevo, la
          esencia de vainilia y la pizca de sal. Integrar todo.</li>
          <li>Tamizar y agregar a la mezcla la harina y el cacao. Romper las 
          nueces en pedacitos, agregarlas y mezclar.</li>
          <li>Verter la mezcla sobre un molde enmantecado y poner el molde en el
          horno recién prendido* a temperatura media. A los ~18 minutos clavarle
          un palito de madera. Si sale masa húmeda todavía le falta.</li>
        </ol>

        <p>
          * Lo de prender el horno al poner la mezcla (en lugar de
          precalentarlo) sería para que se cocine más pareja la mezcla.
        </p>

      </div>
    `));
  })

  parentRouter.use(urlPath, appRouter);
}
