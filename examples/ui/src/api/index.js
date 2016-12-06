import queryString from 'query-string'

class Stack {
  constructor (data, api) {
    this.id = data.id
    this.name = data.name
    this.state = data.state
    this.api = api
  }
  async tasks () {
    const data = await this.api.getJson(`stack/${this.id}/tasks`)
    console.log(data.message)
    return data.message
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
  async logs (query) {
    const data = await this.getJson('log', query)
    return data.entries
  }
  async stacks () {
    const data = await this.getJson('stack')
    console.log(data)
    return data.list.map(s => new Stack(s, this)) || []
  }
}
