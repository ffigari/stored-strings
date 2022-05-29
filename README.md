repo personal, interfaz web levantada en http://ffig.ar/i

```sh
nvm use
npm i
lsof -i:3000 -t && kill $( lsof -i:3000 -t )
node runner.js <action>  # donde <action> puede ser
                         #   - spawn: levantar una instancia del main server en
                         #            el puerto 3000
                         #   - watch: tmb levantar esa instancia pero con reload
                         #            cada vez que se realizan cambios en el
                         #            código
```
```sh
firefox localhost:3000/
```
