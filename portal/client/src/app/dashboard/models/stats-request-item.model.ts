
import { StatsRequest } from '../../models/stats-request.model'
import { GraphStats } from '../../models/graph-stats.model'


export class StatsRequestItem {
  public id : string
  public request : StatsRequest
  public result : GraphStats[]
  public subscriberNumber : number


  constructor(id: string, req : StatsRequest) {
    this.id = id
    this.request = req
  }

}
