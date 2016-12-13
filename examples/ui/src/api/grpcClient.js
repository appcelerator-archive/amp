import queryString from 'query-string'

const codes = [
  'Ok',
  'Canceled',
  'Unknown',
  'Invalid argument',
  'Deadline exceeded',
  'Not found',
  'Already exists',
  'Permission denied',
  'Unauthenticated',
  'Resource exhaused',
  'Failed precondition',
  'Aborted',
  'Out of range',
  'Unimplemented',
  'Interal',
  'Unavailable',
  'Data loss',
]

export default class GrpcClient {
  handleError (json) {
    if (json.Error) {
      const error = new Error(`${codes[json.Code]} error: ${json.Error}`)
      error.code = json.Code
      throw error
    }
  }
  async getJson (path, query) {
    if (query) {
      path += '?' + queryString.stringify(query)
    }
    const request = await fetch(this.base + path)
      .catch(error => {
        console.error(error)
        throw new Error('Network error getting from ' + this.base + path)
      })
    const json = await request.json()
    this.handleError(json)
    return json
  }
  async putJson (path, body, query) {
    if (query) {
      path += '?' + queryString.stringify(query)
    }
    const request = await fetch(this.base + path, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined
    }).catch(error => {
      console.error(error)
      throw new Error('Network error putting to ' + this.base + path)
    })
    const json = await request.json()
    this.handleError(json)
    return json
  }
  async postJson (path, body, query) {
    if (query) {
      path += '?' + queryString.stringify(query)
    }
    const request = await fetch(this.base + path, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined
    }).catch(error => {
      console.error(error)
      throw new Error('Network error posting to ' + this.base + path)
    })
    const json = await request.json()
    this.handleError(json)
    return json
  }
  async deleteJson (path, query) {
    if (query) {
      path += '?' + queryString.stringify(query)
    }
    const request = await fetch(this.base + path, {
      method: 'DELETE'
    }).catch(error => {
      console.error(error)
      throw new Error('Network error deleting ' + this.base + path)
    })
    const json = await request.json()
    this.handleError(json)
    return json
  }
}
