import sleep from './sleep'

export default class Topic {
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
