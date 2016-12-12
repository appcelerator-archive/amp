import queryString from 'query-string'

function sleep (ms) {
  return new Promise((resolve, reject) => {
    setTimeout(resolve, ms)
  })
}

class Stack {
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
}

class Topic {
  constructor (data, api) {
    this.id = data.id
    this.name = data.name
    this.api = api
  }
  async remove () {
    const results = await this.api.deleteJson(`topic/${this.id}`)
    await sleep(100)
    return results
  }
}

export default class AmpApi {
  constructor (base) {
    this.base = base || 'http://amplifier-api.local.appcelerator.io/v1/'
  }
  async getJson (path, query) {
    if (query) {
      path += '?' + queryString.stringify(query)
    }
    const request = await fetch(this.base + path)
    const json = await request.json()
    return json
  }
  async postJson (path, body, query) {
    if (query) {
      path += '?' + queryString.stringify(query)
    }
    const request = await fetch(this.base + path, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined
    })
    const json = await request.json()
    return json
  }
  async deleteJson (path, query) {
    if (query) {
      path += '?' + queryString.stringify(query)
    }
    const request = await fetch(this.base + path, {
      method: 'DELETE'
    })
    const json = await request.json()
    return json
  }
  async logs (query) {
    const data = await this.getJson('log', query)
    return data.entries || []
  }
  async stacks (query) {
    const data = await this.getJson('stack', query)
    return data.list ? data.list.map(s => new Stack(s, this)) : []
  }
  async createStack (stack) {
    const result = await this.postJson('stack', {stack})
    const data = {
      id: result.stack_id,
      name: stack.name,
      state: 'Stopped',
    }
    return new Stack(data, this)
  }
  async topics () {
    const data = await this.getJson('topic')
    return data.topics ? data.topics.map(t => new Topic(t, this)) : []
  }
  async createTopic (name) {
    const topic = {name}
    const result = await this.postJson('topic', {topic})
    return new Topic(result.topic, this)
  }
}
