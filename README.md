web personal, levantada en http://ffig.ar/i

```sh
nvm use
npm i
lsof -i:3000 -t && kill $( lsof -i:3000 -t )
node runner.js  # o `node runner.js watch` para que haga reload al realizar
                # cambios en '/api'
```
```sh
firefox localhost:3000/
```
