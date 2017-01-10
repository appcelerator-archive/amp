export default class Func {
  constructor (data, api) {
    this.id = data.id
    this.name = data.name
    this.image = data.image
    this.api = api
  }
  async run (string) {
    const base = 'http://localhost:4242/'
    const path = this.id
    const request = await fetch(base + path, {
      method: 'POST',
      body: string
    }).catch(error => {
      console.error(error)
      throw new Error('Network error posting to ' + base + path)
    })
    return await request.text()
  }
  async remove () {
    const results = await this.api.deleteJson(`function/${this.id}`)
    return results
  }
}
