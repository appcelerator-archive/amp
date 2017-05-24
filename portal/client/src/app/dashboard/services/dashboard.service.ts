import { Injectable } from '@angular/core';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { Subject } from 'rxjs/Subject'
import { Graph } from '../../models/graph.model';
import { StatsRequest } from '../../models/stats-request.model';
import { StatsRequestItem } from '../models/stats-request-item.model';
import { GraphStats } from '../../models/graph-stats.model'
import * as d3 from 'd3';

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
    period : string = '2m'
    timer : any
    requestMap = {}
    nbGraph = 1
    public showEditor = false;
    public showAlert = false;
    public graphColors = []
    editorGraph : Graph = new Graph(1, this.x0, this.y0, this.w0, this.h0, 'editor', [''], '','')
    public notSelected : Graph = new Graph(0, 0, 0, 0, 0, "", [], "","")
    public selected : Graph = this.notSelected

  constructor(
    private httpService : HttpService,
    private menuService : MenuService) {
      for (let i=0;i<20;i++) {
        this.graphColors.push(d3.interpolateCool(Math.random()))
      }
      this.notSelected.title = ""
      this.notSelected.object="stack"
      this.notSelected.field="cpu-usage"
      this.notSelected.topNumber=3
      this.notSelected.border=true
      this.yTitleMap['cpu-usage'] = 'cpu usage (%)'
      this.yTitleMap['mem-limit'] = 'memory limit (bytes)'
      this.yTitleMap['mem-maxusage'] = 'memory max usage (bytes)'
      this.yTitleMap['mem-usage'] = 'memory usage (MB)'
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
      this.timer = setInterval(() => this.executeRequests(), this.refresh * 1000)
      this.menuService.onRefreshClicked.subscribe(
        () => this.executeRequests()
      )
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
    if (graph.type == "pie") {
      graph.width = graph.height
    }
    graph.title = this.notSelected.title
    graph.object = this.notSelected.object
    graph.field = this.notSelected.field
    graph.border = this.notSelected.border
    this.graphs.push(graph)
    this.addRequest(graph)
    this.onNewData.next()
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

  getTopLabel() : string {
    if (this.selected.topNumber == 0) {
      return 'all'
    }
    return 'top'+this.selected.topNumber
  }

  setRefreshPeriod(refresh : number) {
    this.refresh = refresh;
    this.cancelRequests()
    this.timer = setInterval(() => this.executeRequests(), this.refresh * 1000)
  }

  setPeriod(period : string) {
    this.period = period;
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

  setTop(top : number) {
      this.selected.topNumber = top
      this.addRequest(this.selected)
      this.onNewData.next()
  }

  setTitle(title : string) {
    this.selected.title = title
    this.onNewData.next()
  }

  setAlert(val : boolean) {
    this.selected.alert = val;
    this.onNewData.next()
  }
  setMinAlert(val : string) {
    this.selected.alertMin = +val;
    this.onNewData.next()
  }

  setMaxAlert(val : string) {
    this.selected.alertMax = +val;
    this.onNewData.next()
  }

  setBorder(border : boolean) {
    this.selected.border = border
    this.onNewData.next()
  }

  setCriterionValue(val : string) {
    this.selected.criterionValue = val
    this.addRequest(this.selected)
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
        req.result = data
        this.onNewData.next(req.id)
      },
      (err) => {
        console.log("request error")
        console.log(err)
      }
    )
  }

  addRequest(graph : Graph) : string {
    if (graph.title == '' || graph.title == 'stacks' || graph.title == 'services' || graph.title == 'containers' || graph.title == 'nodes') {
      graph.title = graph.object+'s'
    }
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
    } else if (graph.object == "container") {
      req.group="container_short_name"
    } else if (graph.object == "node") {
      req.group="node_id"
    } else {
      return
    }

    req.period = "now-2m"
    req.stats_cpu = true
    req.stats_mem = true
    req.stats_net = true
    req.stats_io = true
    console.log(graph)
    if (graph.object == 'stack') {
      req.filter_stack_name = graph.criterionValue
    } else if (graph.object == 'service') {
      req.filter_service_name = graph.criterionValue
    } else if (graph.object == 'container') {
      req.filter_container_id = graph.criterionValue
    } else if (graph.object == 'node') {
      req.filter_node_id = graph.criterionValue
    }
    let id = graph.object+"-"+graph.criterionValue
    let newItem = new StatsRequestItem(id, req)
    newItem.subscriberNumber=1
    this.requestMap[id]=newItem
    graph.requestId = id
    this.executeRequest(newItem)
    return id;
  }

  getData(graph : Graph) : GraphStats[] {
    if (!graph.requestId) {
      return []
    }
    let item = this.requestMap[graph.requestId]
    if (!item) {
      return []
    }
    if (!item.result) {
      return []
    }
    let list = this.sortByField(item.result, graph.field)
    if (graph.topNumber == 0) {
      return list
    }
    return list.slice(0, graph.topNumber)
  }

  sortByField(data : GraphStats[], field : string) : GraphStats[] {
    return data.sort((a, b) => {
      if (a.values[field] < b.values[field]) {
        return 1;
      }
      return -1
    })
  }

  isVisible(type : string) {
    if (type == 'object' || type == 'field') {
      if (this.selected.type != 'text') {
        return true
      }
      return false
    }
    if (type == 'top') {
      if (this.selected.type != 'text' && this.selected.type != 'counter') {
        return true
      }
      return false;
    }
    if (type == 'alert' || type == 'criterion' || type == 'criterionValue') {
      if (this.selected.type == 'counter') {
        return true
      }
      return false
    }
    return false
  }

}
