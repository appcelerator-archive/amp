import { GraphHistoricData } from './graph-historic-data.model'

export class GraphHistoricAnswer {
  public names : string[]
  public data: GraphHistoricData[]

  constructor(names : string[], data: GraphHistoricData[]) {
    this.names = names
    this.data = data
  }
}
