const { send } = require('micro')
const sleep = require('then-sleep')

module.exports = async function (req, res) {
  console.log(`${req.method} ${req.url}`)
  await sleep(500)
  send(res, 200, `[${req.method} ${req.url}] hello.`)
}
