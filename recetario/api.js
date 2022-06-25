const express = require('express')

const { wrapBody } = require('./../common/index.js')

module.exports.addItselfTo = async (parentRouter, connectionString) => {
  const urlPath = '/recetario';
  const appRouter = express.Router();

  appRouter.get('/budin-de-banana', (req, res) => {
    res.send(wrapBody(`
      <div class="my-3">

        <h1>budín de banana</h1>

        <h2>ingredientes</h2>

        <ul>
          <li>75 grs de manteca</li>
          <li>1 banana, entre más madura mejor</li>
          <li>1 huevo</li>
          <li>una taza de azucar rubia</li>
          <li>una taza de harina leudante</li>
          <li>un chorrito de aceite</li>
          <li>dsp ya más opcional, chocolate trozeado, una cucharada de cacao
          en polvo, un chorrito de crema</li>
        </ul>

        <h2>preparación</h2>
        <ol>
          <li>Con un tenedor pisotear la banana junto a la manteca. La manteca
          conviene sacarla con tiempo y si no se le puede dar un toque de baño
          maría o de microondas </li>
          <li>Ir agregando el resto, mezclando bien en cada paso</li>
          <li>Los extras (eg, el chocolate) se pueden ir agregando de a
          <em>capas</em>. En lugar de mandar todo junto a la mezcla, se manda
          ponele un tercio. Ahí mezclás bien y tirás un tercio de la mezcla al
          molde. Dsp volvés a mandar otra parte a la mezcla, incorporás, vertés
          en el molde. Así debería ir quedando más chocolate arriba en lugar de
          irse todo para abajo, como cuando hacías una gelatina con fruta 
          cortada y quedaba toda contra el fondo.</li>
          <li>Mandar al horno ~20 minutos, los primero 5 a fuego fuerte y
          habiendo precalentado previamente el horno. Antes de sacarlo clavarle
          un palito para chequear que no le falte cocción a la masa.</li>
        </ol>
      </div>
    `));
  })

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

        <p>
          igual salió medio bizcochuelo así que qué sé yo
        </p>

      </div>
    `));
  })

  parentRouter.use(urlPath, appRouter);
}
