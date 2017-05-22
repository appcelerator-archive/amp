import { Injectable } from '@angular/core';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { Subject } from 'rxjs/Subject'
import { Graph } from '../../models/graph.model';
import { StatsRequest } from '../../models/stats-request.model';
import { StatsRequestItem } from '../models/stats-request-item.model';
import { GraphStats } from '../../models/graph-stats.model'

@Injectable()
export class DashboardService {
    graphs : Graph[] = []
    editor = false
    onNewData = new Subject<string>();
    onGraphSelect = new Subject<Graph>()
    yTitleMap = {}
    x0 = 20
    y0 = 20
    w0 = 300
    h0 = 150
    refresh : number = 30
    timer : any
    requestMap = {}
    nbGraph = 1
    public showEditor = false;
    public graphColors = ['DodgerBlue', 'slateblue', 'blue', 'magenta', 'pink', 'green', 'ping', 'orange', 'red', 'yellow', 'blue']
    editorGraph : Graph = new Graph(1, this.x0, this.y0, this.w0, this.h0, 'editor', [''], '','')
    public notSelected : Graph = new Graph(0, 0, 0, 0, 0, "", [], "","")
    public selected : Graph = this.notSelected

  constructor(
    private httpService : HttpService,
    private menuService : MenuService) {
      this.notSelected.title = ""
      this.notSelected.object=""
      this.notSelected.field=""
      this.notSelected.border=true
      this.yTitleMap['cpu-usage'] = 'cpu usage (%)'
      this.yTitleMap['mem-limit'] = 'memory limit (bytes)'
      this.yTitleMap['mem-maxusage'] = 'memory max usage (bytes)'
      this.yTitleMap['mem-usage'] = 'memory usage (bytes)'
      this.yTitleMap['mem-usage-p'] = 'memory usage (%)'
      this.yTitleMap['net-total-bytes'] = 'network traffic (bytes)'
      this.yTitleMap['net-rx-bytes'] = 'network rx traffic (bytes)'
      this.yTitleMap['net-rx-packets'] = 'network rx traffic (packets)'
      this.yTitleMap['net-tx-bytes'] = 'network tx traffic (bytes)'
      this.yTitleMap['net-tx-packets'] = 'network tx traffic (packets)'
      this.yTitleMap['io-total'] = 'io r/w (bytes)'
      this.yTitleMap['io-write'] = 'io write (bytes)'
      this.yTitleMap['io-read'] = 'io read (bytes)'
      this.cancelRequests()
      this.timer = setInterval(this.executeRequest, this.refresh*1000)
    }

  cancelRequests() {
    if (this.timer) {
      console.log("clear interval")
      clearInterval(this.timer);
    }
  }

  addGraph(type : string) {
    this.x0 += 20
    this.y0 += 20
    this.nbGraph++;
    let graph = new Graph(this.nbGraph, this.x0, this.y0, this.w0, this.h0, type, [''], "essai de titre",'')
    graph.title = this.notSelected.title
    graph.object = this.notSelected.object
    graph.field = this.notSelected.field
    graph.border = this.notSelected.border
    this.graphs.push(graph)
  }

  removeSelectedGraph() {
    let list = []
    for (let graph of this.graphs) {
      if (graph.id != this.selected.id) {
        list.push(graph)
      }
    }
    this.graphs = list
  }

  setRefreshPeriod(refresh : number) {
    this.refresh = refresh;
  }

  setObject(name : string) {
    this.selected.object = name
    this.addRequest(this.selected)
    this.onNewData.next()
  }

  setField(name : string) {
    this.selected.field = name
    this.addRequest(this.selected)
    this.onNewData.next()
  }

  setTitle(title : string) {
    this.selected.title = title
    this.onNewData.next()
  }

  setBorder(border : boolean) {
    this.selected.border = border
    this.onNewData.next()
  }

  getTextWidth(text, fontSize, fontFace) : number {
    var a = document.createElement('canvas');
    var b = a.getContext('2d');
    b.font = fontSize + 'px ' + fontFace;
    return b.measureText(text).width;
  }

  executeRequests() {
    for (let id in this.requestMap) {
      this.executeRequest(this.requestMap[id])
    }
  }

  executeRequest(req : StatsRequestItem) {
    if (!req) {
      return
    }
    console.log("execute request: "+req.request.group)
    console.log(req.request)
    this.httpService.statsCurrent(req.request).subscribe(
      (data) => {
        req.result=data
        console.log(data)
        this.onNewData.next(req.id)
      },
      (err) => {
        console.log("request error")
        console.log(err)
      }
    )
  }

  addRequest(graph : Graph) : string {
    let item = this.requestMap[graph.object]
    if (item) {
      graph.requestId = item.id;
      return
    }
    let req = new StatsRequest()
    if (graph.object == "stack") {
        req.group="stack_name"
    } else if (graph.object == "service") {
      req.group="service_name"
    } else if (req.group == "container") {
      req.group="container_short_name"
    } else {
      return
    }

    req.period = "now-10m"
    req.stats_cpu = true
    req.stats_mem = true
    req.stats_net = true
    req.stats_io = true
    let id = graph.object
    let newItem = new StatsRequestItem(id, req)
    newItem.subscriberNumber=1
    this.requestMap[id]=newItem
    graph.requestId = id
    this.executeRequest(newItem)
    return id;
  }

  getData(id : string) : GraphStats[] {
    if (!id) {
      return []
    }
    let item = this.requestMap[id]
    if (!item) {
      return []
    }
    if (!item.result) {
      return []
    }
    return item.result;
  }




}
