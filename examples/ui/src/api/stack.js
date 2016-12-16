import sleep from './sleep'

export default class Stack {
  constructor (data, api) {
    this.id = data.id
    this.name = data.name
    this.state = data.state
    this.api = api
  }
  async tasks () {
    const data = await this.api.getJson(`stack/${this.id}/tasks`)
    return data.message
  }
  async logs () {
    const data = await this.api.logs({
      stack: this.name
    })
    return data
  }
  async start () {
    const results = await this.api.postJson(`stack/${this.id}/start`)
    await sleep(100)
    this.state = 'Running'
    return results
  }
  async stop () {
    const results = this.api.postJson(`stack/${this.id}/stop`)
    await sleep(100)
    this.state = 'Stopped'
    return results
  }
  async remove () {
    const results = await this.api.deleteJson(`stack/${this.id}`)
    await sleep(100)
    return results
  }
  async details () {
    const results = await this.api.getJson(`stack/${this.id}`)
    return results.stack
  }
}
