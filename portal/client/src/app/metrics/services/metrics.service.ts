import { Injectable } from '@angular/core';
import { GraphHistoricData } from '../../models/graph-historic-data.model';
import { StatsRequest } from '../../models/stats-request.model';
import { GraphDataAnswer } from '../models/graph-data-answer.model';
import { GraphLine} from '../models/graph-line.model';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { Subject } from 'rxjs/Subject'

@Injectable()
export class MetricsService {
  histoData : GraphHistoricData[] = []
  public lineVisibleMap = {}
  lines : GraphLine[] = []
  graphColorsStack = ['blue', 'slateblue', 'blue', 'pink', 'green', 'pink', 'orange', 'red', 'yellow', 'blue']
  graphColorsService = ['slateblue', 'blue', 'DodgerBlue', 'pink', 'green', 'orange', 'red', 'yellow', 'blue']
  graphColorsContainer = ['green', 'orange', 'blue', 'magenta', 'pink', 'green', 'orange', 'red', 'yellow', 'blue']
  graphColors = ['dodgerBlue', 'pink', 'blue', 'pink', 'green', 'orange', 'red', 'yellow', 'blue']
  statsRequest : StatsRequest
  onNewData = new Subject();
  timePeriod = "now-10m"
  timeGroup = "30s"
  periodRefresh = 30
  object = ""
  type = ""
  ref = ""
  clickPossible : true

  constructor(
    private httpService : HttpService,
    private menuService : MenuService) { }

  setHistoricRequest(request : StatsRequest, period : number) {
    this.statsRequest = request
    if (period < 5) {
      period = 5;
    }
    this.updateHistoricData()
    this.cancelRequests()
    this.menuService.setCurrentTimer(setInterval(() => this.updateHistoricData(), period * 1000))
  }

  cancelRequests() {
    this.menuService.clearCurrentTimer()
  }

  getColor(index : number) {
    if (this.object == 'stack') {
      return this.graphColorsStack[index]
    } else if (this.object == 'service') {
      return this.graphColorsService[index]
    } else if (this.object == 'container') {
      return this.graphColorsContainer[index]
    } else {
      return this.graphColors[index]
    }
  }

  getHistoricData(fields : string[], object : string, graphType : string) : GraphDataAnswer {
    let data = []
    let lines : GraphLine[] = []
    //console.log(fields)
    //console.log("type="+graphType+" data.length="+this.histoData.length)
    //console.log(this.histoData)
    if (graphType=='single') {
      for (let ii=0; ii<fields.length; ii++) {
        let name = fields[ii]
        lines.push(new GraphLine(name, this.graphColors[ii]))
        if (this.lineVisibleMap[name] === undefined) {
          this.lineVisibleMap[name]=true
        }
      }
      this.histoData.forEach( (ele) => {
        let ret = []
        for (let field of fields) {
          ret.push(ele.values[field])
        }
        let newEle = new GraphHistoricData(ele.date)
        newEle.graphValues = ret
        data.push(newEle)
      })
    }
    else {
      if (this.histoData.length>0) {
        let date = this.histoData[0].date
        let ret = []
        let localLineRefMap = {}
        this.histoData.forEach( (ele : GraphHistoricData) => {
          if (date.getTime() !== ele.date.getTime()) {
            let newEle = new GraphHistoricData(ele.date)
            newEle.graphValues = ret
            data.push(newEle)
            date = ele.date
            ret = []
          }
          if (ele.name.toLowerCase() != 'nostack' && ele.name.toLowerCase() != 'noservice') {
            ret.push(ele.values[fields[0]])
            if (localLineRefMap[ele.name] === undefined) {
              localLineRefMap[ele.name] = ele.name
              let line = new GraphLine(ele.name, this.graphColors[lines.length])
              lines.push(line)
              if (this.lineVisibleMap[line.name] === undefined) {
                this.lineVisibleMap[line.name]=true
              }
            }
          }
        })
        let newEle = new GraphHistoricData(date)
        newEle.graphValues = ret
        data.push(newEle)
      }
    }
    //console.log(lines)
    //console.log(data)
    this.lines=lines
    return new GraphDataAnswer(lines, data)
  }

  updateHistoricData() {
    this.statsRequest.format=true
    this.httpService.statsHistoric(this.statsRequest).subscribe(
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

  set(object : string, type : string, ref : string) {
    this.object = object
    this.type = type
    this.ref = ref
  }

  setTimePeriod(period :string, group : string) {
    this.timePeriod = period
    this.timeGroup = group
    if (this.statsRequest) {
      this.statsRequest.period = period
      this.statsRequest.time_group = group
      this.updateHistoricData()
    }
  }

  setContainerAvg(val : boolean) {
    this.statsRequest.avg = val
    this.updateHistoricData()
  }

  setRefreshPeriod(period :string) {
    this.periodRefresh = parseInt(period)
    this.menuService.clearCurrentTimer()
    this.menuService.setCurrentTimer(setInterval(() => this.updateHistoricData(), parseInt(period) * 1000))
  }

  toggleLine(name: string) {
    this.lineVisibleMap[name]=!this.lineVisibleMap[name]
    this.onNewData.next()
  }

  route(ref : string) {
    if (this.object == 'global') {
      if (this.type == 'multi') {
        this.menuService.navigate(['/amp', 'metrics', 'stack', 'multi', 'all'])
      } else {
        this.menuService.navigate(['/amp', 'metrics', 'stack', 'multi', 'all'])
      }
    }
    if (this.object == 'stack') {
      if (this.type == 'multi') {
        this.menuService.navigate(['/amp', 'metrics', 'service', 'multi', ref])
      } else {
        this.menuService.navigate(['/amp', 'metrics', 'service', 'single', this.ref])
      }
    }
    if (this.object == 'service') {
      if (this.type == 'multi') {
        this.menuService.navigate(['/amp', 'metrics', 'task', 'multi', ref])
    } else {
        this.menuService.navigate(['/amp', 'metrics', 'task', 'single', this.ref])
      }
    }
  }
}
