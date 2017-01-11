import GrpcClient from './grpcClient'
import Stack from './stack'
import Topic from './topic'
import Func from './function'
import Storage from './storage'

export default class AmpApi extends GrpcClient {
  constructor (base) {
    super()
    this.base = base || 'http://amplifier-api.local.appcelerator.io/v1/'
  }
  async logs (query) {
    const data = await this.getJson('log', query)
    return data.entries || []
  }
  async stats (query) {
    const data = await this.getJson('stats', query)
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
  async functions () {
    const data = await this.getJson('function')
    return data.functions ? data.functions.map(f => new Func(f, this)) : []
  }
  async createFunction (name, image) {
    const func = {name, image}
    const result = await this.postJson('function', {function: func})
    return new Func(result.function, this)
  }
  async storageObjects () {
    const data = await this.getJson('storage')
    return data.list ? data.list.map(s => new Storage(s, this)) : []
  }
  async createStorageObject (key, val) {
    const results = await this.putJson(`storage/${key}`, {val})
    return new Storage(results, this)
  }
}
