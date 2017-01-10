export default class Storage {
  constructor (data, api) {
    this.key = data.key
    this.val = data.val
    this.api = api
  }
  async update () {
    const results = await this.api.putJson(`storage/${this.key}`, {val: this.val})
    return results
  }
  async remove () {
    const results = await this.api.deleteJson(`storage/${this.key}`)
    return results
  }
}
