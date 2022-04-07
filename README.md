```sh
nvm use
npm i
lsof -i:3000 -t && kill $( lsof -i:3000 -t )
node runner.js
```
```sh
firefox localhost:3000/
```

