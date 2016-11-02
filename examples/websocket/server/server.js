const port = 4000
const { Server } = require('ws')
const wss = new Server({ port })
const pty = require('pty.js')

wss.on('connection', ws => {
  const term = pty.spawn('bash', [], {
    name: 'xterm-color',
    cols: 80,
    rows: 30,
    cwd: '/app',
    env: process.env
  })
  ws.on('message', m => term.write(m))
  term.on('data', d => {
    ws.send(d)
    process.stdout.write(d)
  })
  term.on('exit', (c) => ws.close())
})

console.log('listening on', port)
