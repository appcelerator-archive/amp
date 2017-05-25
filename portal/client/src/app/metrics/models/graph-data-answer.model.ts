import { GraphHistoricData } from '../../models/graph-historic-data.model'
import { GraphLine} from './graph-line.model'

export class GraphDataAnswer {
  public lines: GraphLine[]
  public data: GraphHistoricData[]

  constructor(lines: GraphLine[], data: GraphHistoricData[]) {
    this.lines = lines
    this.data = data
  }
}
