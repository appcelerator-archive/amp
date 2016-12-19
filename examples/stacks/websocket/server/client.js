const WebSocket = require('ws')
const ws = new WebSocket('ws://localhost:4000/')
 
ws.on('open', () => {
  process.stdin.on('data', d => ws.send(d))
  process.stdin.resume()
})
 
ws.on('message', (data, flags) => {
  process.stdout.write(data)
})

ws.on('close', () => {
  process.exit()
})
