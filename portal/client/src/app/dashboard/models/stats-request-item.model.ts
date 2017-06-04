
import { StatsRequest } from '../../models/stats-request.model'
import { GraphCurrentData } from '../../models/graph-current-data.model'
import { GraphHistoricData } from '../../models/graph-historic-data.model'

export class StatsRequestItem {
  public id : string
  public request : StatsRequest
  public graphTitle: string
  public currentResult : GraphCurrentData[]
  public historicResult : GraphHistoricData[]
  public subscriberNumber : number


  constructor(id: string, req : StatsRequest, title: string) {
    this.id = id
    this.request = req
    this.graphTitle = title
  }

}
