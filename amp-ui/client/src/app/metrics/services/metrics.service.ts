import { Injectable } from '@angular/core';
import { GraphHistoricData } from '../models/graph-historic-data.model';
import { StatsRequest } from '../models/stats-request.model';
import { HttpService } from '../../services/http.service';
import { Subject } from 'rxjs/Subject'

@Injectable()
export class MetricsService {
  histoData : GraphHistoricData[] = []
  statsRequest : StatsRequest
  onNewData = new Subject();
  timer : any

  constructor(private httpService : HttpService) { }

  setHistoricRequest(request : StatsRequest, period : number) {
    this.statsRequest = request
    if (period < 5) {
      period = 5;
    }
    this.updateHistoricData()
    this.timer = setInterval(() => this.updateHistoricData(), period * 1000)
  }

  cancelRequests() {
    if (this.timer) {
      clearInterval(this.timer);
    }
  }

  getHistoricData(fields : string[]) : GraphHistoricData[] {
    let data = []
    this.histoData.forEach( (ele) => {
      let ret = []
      for (let field of fields) {
        ret.push(ele.values[field])
      }
      let newEle = new GraphHistoricData(ele.date, undefined)
      newEle.graphValues = ret
      data.push(newEle)
    })
    //console.log(data)
    return data
  }

  updateHistoricData() {
    this.httpService.stats(this.statsRequest).subscribe(
      data => {
        this.histoData = data
        //console.log(data)
        this.onNewData.next()
      },
      error => {
        //console.log("loadStacksError")
        console.log(error)
      }
    );
  }

  setPeriod(period :string, group : string) : string {
    if (this.statsRequest) {
      this.statsRequest.period = period
      this.statsRequest.time_group = group
      this.updateHistoricData()
      return this.buildPeriodLabel(period, group)
    }
  }

  buildPeriodLabel(period : string, group : string) : string {
    return "time selection"
  }

}
