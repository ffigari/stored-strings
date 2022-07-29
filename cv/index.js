const fs = require('fs')

module.exports.printCV = () => fs.readFile(
  './cv/main.tex', 'utf8',
  (err, data) => console.log(data.trim()))

