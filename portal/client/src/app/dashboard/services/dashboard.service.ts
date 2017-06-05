import { Injectable } from '@angular/core';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { Subject } from 'rxjs/Subject'
import { Graph } from '../../models/graph.model';
import { DashboardInnerStats } from '../models/dashboard-inner-stats.model';
import { StatsRequest } from '../../models/stats-request.model';
import { StatsRequestItem } from '../models/stats-request-item.model';
import { GraphCurrentData } from '../../models/graph-current-data.model'
import { GraphHistoricData } from '../../models/graph-historic-data.model'
import { GraphHistoricAnswer } from '../../models/graph-historic-answer.model'
import { ColorsService } from './colors.service'
import * as d3 from 'd3';


@Injectable()
export class DashboardService {
    graphs : Graph[] = []
    editor = false
    onNewData = new Subject<string>();
    onGraphSelect = new Subject<Graph>()
    yTitleMap = {}
    unit = {}
    x0 = 280
    y0 = 5
    w0 = 300
    h0 = 150
    refresh : number = 30
    period : string = 'now-2m'
    requestMap = {}
    nbGraph = 1
    public innerStats : DashboardInnerStats = new DashboardInnerStats()
    public showEditor = false;
    public showAlert = false;
    public nodeColorIndex = 0
    public editorGraph : Graph = new Graph('graph1', this.x0, this.y0, this.w0, this.h0, 'editor','')
    public notSelected : Graph = new Graph('', 0, 0, 0, 0, "", "")
    public selected : Graph = this.notSelected
    public isVisible = false
    public ModeSetting = "setting"
    public ModeLinkLegendToGraph = "linkLegendToGraph"
    public clickMode = this.ModeSetting


  constructor(
    private httpService : HttpService,
    private menuService : MenuService,
    private colorsService: ColorsService) {
      this.notSelected.title = ""
      this.notSelected.object="stack"
      this.notSelected.field="cpu-usage"
      this.notSelected.topNumber=3
      this.notSelected.border=true
      this.yTitleMap['cpu-usage'] = 'cpu usage'
      this.yTitleMap['mem-limit'] = 'memory limit'
      this.yTitleMap['mem-maxusage'] = 'memory max usage'
      this.yTitleMap['mem-usage'] = 'memory usage'
      this.yTitleMap['mem-usage-p'] = 'memory usage'
      this.yTitleMap['net-total-bytes'] = 'network traffic'
      this.yTitleMap['net-rx-bytes'] = 'network rx traffic'
      this.yTitleMap['net-rx-packets'] = 'network rx traffic'
      this.yTitleMap['net-tx-bytes'] = 'network tx traffic'
      this.yTitleMap['net-tx-packets'] = 'network tx traffic'
      this.yTitleMap['io-total'] = 'io r/w'
      this.yTitleMap['io-write'] = 'io write'
      this.yTitleMap['io-read'] = 'io read'
      //
      this.unit['cpu-usage'] = '%'
      this.unit['mem-limit'] = 'bytes'
      this.unit['mem-maxusage'] = 'bytes'
      this.unit['mem-usage'] = 'bytes'
      this.unit['mem-usage-p'] = '%'
      this.unit['net-total-bytes'] = 'bytes'
      this.unit['net-rx-bytes'] = 'bytes'
      this.unit['net-rx-packets'] = 'packets'
      this.unit['net-tx-bytes'] = 'bytes'
      this.unit['net-tx-packets'] = 'packets'
      this.unit['io-total'] = 'bytes'
      this.unit['io-write'] = 'bytes'
      this.unit['io-read'] = 'bytes'
      this.cancelRequests()
      this.menuService.setCurrentTimer(setInterval(() => {
          this.redisplay()
        }, this.refresh * 1000)
      )
      this.menuService.onRefreshClicked.subscribe(
        () => {
          if (this.isVisible) {
            this.redisplay()
          }
        }
      )
    }

  cancelRequests() {
    this.menuService.clearCurrentTimer()
  }

  addGraph(type : string, offTop : number, offLeft : number) {
    this.nbGraph++;
    let graph = new Graph('graph'+this.nbGraph, this.x0-offLeft, this.y0-offTop, this.w0, this.h0, type, '')
    this.x0 += 2
    this.y0 += 2
    if (graph.type == "pie") {
      graph.width = graph.height
    }
    graph.title = this.notSelected.title
    graph.object = this.notSelected.object
    graph.field = this.notSelected.field
    if (graph.type == 'counterSquare' || graph.type == 'counterCircle') {
      graph.field = 'number'
      if (graph.type == 'counterSquare') {
        graph.counterHorizontal = true
        graph.height /= 2;
      }
    }
    if (graph.type == 'lines' || graph.type == 'areas') {
      graph.histoPeriod = 'now-10m'
    }
    if (graph.type == 'bubbles') {
      graph.bubbleXField = 'mem-usage'
      graph.bubbleYField = 'cpu-usage'
      graph.field = 'net-total-bytes'
      graph.bubbleScale = 'medium'
      graph.topNumber = 0
    }
    if (graph.type == 'areas') {
      graph.stackedAreas = true
    }
    this.graphs.push(graph)
    this.addRequest(graph)
    //this.onNewData.next()
    this.selected = graph
  }

  copySelected() {
    //
  }

  addLegend(object : string) {
    this.x0 += 2
    this.y0 += 2
    this.nbGraph++;
    let graph = new Graph('graph'+this.nbGraph, this.x0, this.y0, this.w0*2/3, this.h0, "legend", "Legend "+object+"s")
    graph.object=object
    this.graphs.push(graph)
  }

  addInnerStats() {
    this.x0 += 2
    this.y0 += 2
    this.nbGraph++;
    let graph = new Graph('graph'+this.nbGraph, this.x0, this.y0, this.w0*2/3, this.h0, "innerStats", "Inner stats")
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
    delete this.requestMap[this.selected.requestId];
    this.selected = this.notSelected
  }

  toggleEditor(offsetTop : number, offsetLeft : number) {
    if (this.showEditor) {
      this.showEditor = false
      this.selected = this.notSelected
    } else {
      this.showEditor = true
      this.editorGraph.x = -offsetLeft
      this.editorGraph.y = -offsetTop
    }
  }

  getTopLabel() : string {
    if (this.selected.topNumber == 0) {
      return 'all'
    }
    return 'top'+this.selected.topNumber
  }

  getObjectColor(graph : Graph, name : string) : string {
    return this.colorsService.getColor(graph, name)
  }

  computeUnit(field : string, val : number, refUnit: string) : {val: number, sval: string, unit: string} {
    if (this.unit[field] == '%') {
      return { val: val, sval: val.toFixed(1)+' %', unit: '%'}
    }
    if (this.unit[field]!='bytes') {
      return { val: val, sval: val.toFixed(0)+" "+this.unit[field], unit: this.unit[field]}
    }
    if (refUnit) {
      //force to compute in the refUnit unit
      return {val: val, sval: (val/this.unitdivider(refUnit)).toFixed(1)+" "+refUnit, unit: refUnit }
    }
  	if (val < 1024) {
  		return {val: val, sval: val.toFixed(0)+' Bytes', unit: 'Bytes'}
  	} else if (val < 1048576) {
  		return {val: (val/1024), sval: (val/1024).toFixed(1)+' KB', unit: 'KB'}
  	} else if (val < 1073741824) {
  		return {val: (val/1048576), sval: (val/1048576).toFixed(1)+ ' MB', unit: 'MB'}
  	}
  	return {val: (val/1073741824), sval: (val/1073741824).toFixed(1)+' GB', unit: 'GB'}
  }

  adjustCurrentDataToUnit(unit : string, field : string, data : GraphCurrentData[]) : GraphCurrentData[] {
    let div = this.unitdivider(unit)
    for (let gdata of data) {
      gdata.valueUnit = gdata.values[field]/div
    }
    return data
  }

  adjustCurrentXYDataToUnit(unitx : string, unity : string, fieldx : string, fieldy : string, data : GraphCurrentData[]) : GraphCurrentData[] {
    let divx = this.unitdivider(unitx)
    let divy = this.unitdivider(unity)
    for (let gdata of data) {
      gdata.valueUnitx = gdata.values[fieldx]/divx
      gdata.valueUnity = gdata.values[fieldy]/divy
    }
    return data
  }

  adjustHistoricDataToUnit(unit : string, field : string, data : GraphHistoricData[]) : GraphHistoricData[] {
    let div = this.unitdivider(unit)
    for (let hdata of data) {
      for (let ii=0; ii< hdata.graphValues.length; ii++) {
        hdata.graphValuesUnit[ii] = hdata.graphValues[ii] / div
      }
    }
    return data
  }

  unitdivider(unit : string) : number {
    if (unit == 'KB') {
      return 1024
    } else if (unit == 'MB') {
      return 1048576
    } else if (unit == 'GB') {
      return 1073741824
    }
    return 1;
  }

  computeUnitFormat(graph : Graph, val : number, unit : string): string {
    if (this.unit[graph.field] != "bytes") {
      return val +" "+unit
    }
    let div = this.unitdivider(unit)
    return val/div+" "+unit
  }

  setRefreshPeriod(refresh : number) {
    this.refresh = refresh;
    this.cancelRequests()
    this.menuService.setCurrentTimer(setInterval(() => {
      this.redisplay()
    },this.refresh * 1000))
  }

  redisplay() {
    this.colorsService.refresh()
    this.executeRequests()
  }

  setPeriod(period : string) {
    this.period = period;
    for (let id in this.requestMap) {
      let req = this.requestMap[id]
      if (req) {
        req.period = period
      }
    }
    this.menuService.onRefreshClicked.next()
    for (let graph of this.graphs) {
      this.addRequest(graph)
    }
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
      this.redisplay()
  }

  setTitle(title : string) {
    this.selected.title = title
    this.onNewData.next()
  }

  setTitleCenter(val : boolean) {
    this.selected.centerTitle = val
    this.onNewData.next()
  }

  setRoundedBox(val : boolean) {
    this.selected.roundedBox = val
    this.onNewData.next()
  }

  setAlert(val : boolean) {
    this.selected.alert = val;
    this.onNewData.next()
  }
  setMinAlert(val : string) {
    this.selected.alertMin = val;
    this.onNewData.next()
  }

  setMaxAlert(val : string) {
    this.selected.alertMax = val;
    this.onNewData.next()
  }

  setBorder(border : boolean) {
    this.selected.border = border
    this.onNewData.next()
  }

  setCriterion(name: string) {
    this.selected.criterion = name
    this.addRequest(this.selected)
    this.onNewData.next()
  }

  setCriterionValue(val : string) {
    this.selected.criterionValue = val
    this.addRequest(this.selected)
    this.onNewData.next()
  }

  setContainerAvg(val : boolean) {
    this.selected.containerAvg = val
    this.addRequest(this.selected)
    this.onNewData.next()
  }

  setHistoPeriod(val : string) {
    this.selected.histoPeriod = val
    this.addRequest(this.selected)
    this.onNewData.next()
  }

  setCounterHorizontal(val : boolean) {
    this.selected.counterHorizontal = val
    this.onNewData.next()
  }

  setBubbleXField(name : string) {
    this.selected.bubbleXField = name
    this.addRequest(this.selected)
    this.onNewData.next()
  }

  setBubbleYField(name : string) {
    this.selected.bubbleYField = name
    this.addRequest(this.selected)
    this.onNewData.next()
  }

  setBubbleScale(name : string) {
    this.selected.bubbleScale = name
    this.addRequest(this.selected)
    this.onNewData.next()
  }

  setStackedAreas(val : boolean) {
    this.selected.stackedAreas = val
    if (val) {
      this.selected.percentAreas = false
    }
    this.onNewData.next()
  }

  setPercentAreas(val : boolean) {
    this.selected.percentAreas = val
    if (val) {
      this.selected.stackedAreas = false
    }
    this.onNewData.next()
  }

  setTransparentLegend(val : boolean) {
    this.selected.transparentLegend = val
    this.onNewData.next()
  }

  setRemoveLocalLegend(val : boolean) {
    this.selected.removeLocalLegend = val
    this.onNewData.next()
  }

  linkLegendToGraph() {
    this.clickMode = this.ModeLinkLegendToGraph
  }

  unlinkLegend() {
    this.selected.legendGraphId = undefined
    this.clickMode = this.ModeSetting
    this.onNewData.next()
  }


  getTextWidth(text, fontSize, fontFace) : number {
    var a = document.createElement('canvas');
    var b = a.getContext('2d');
    b.font = fontSize + 'px ' + fontFace;
    return b.measureText(text).width;
  }

  executeRequests() {
    console.log("nbRequest: "+Object.keys(this.requestMap).length)
    this.innerStats.initNewRefresh()
    for (let id in this.requestMap) {
      this.executeRequest(this.requestMap[id])
    }
  }

  executeRequest(req : StatsRequestItem) {
    if (!req) {
      return
    }
    //console.log(req.id)
    //console.log(req.request)
    let t0 = new Date().getTime()
    if (!req.request.time_group) {
      this.httpService.statsCurrent(req.request).subscribe(
        (data) => {
          //console.log("data size: "+data.length)
          this.innerStats.setRequestTime(t0, req.graphTitle)
          req.currentResult = data
          req.historicResult = []
          this.onNewData.next(req.id)
        },
        (err) => {
          this.innerStats.setRequestError()
          console.log("request error")
          console.log(err)
        }
      )
    } else {
      this.httpService.statsHistoric(req.request).subscribe(
        (data) => {
          this.innerStats.setRequestTime(t0, req.graphTitle)
          req.historicResult = data
          req.currentResult = []
          this.onNewData.next(req.id)
        },
        (err) => {
          this.innerStats.setRequestError()
          console.log("request error")
          console.log(err)
        }
      )
    }
  }

  sumRequest(data : GraphCurrentData[]) : GraphCurrentData[] {
    //console.log(data)
    return data
  }

  addRequest(graph : Graph) : string {
    if (graph.type == "legend" || graph.type == "innerStats") {
      return
    }
    if (graph.title == '' || graph.title == 'stacks' || graph.title == 'services' || graph.title == 'containers' || graph.title == 'nodes') {
      graph.title = graph.object
      if (graph.object != 'all') {
        graph.title += 's'
      }
      if (graph.type=='counterSquare' && graph.counterHorizontal) {
        graph.title +=": "
      }
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
    } else if (graph.object == 'all') {
      req.group=""
    } else {
      return
    }
    req.avg = graph.containerAvg
    if (!req.avg) {
      req.avg = false
    }

    //if (graph.type == "counter" && graph.field != "number") {
    //  req.group="container_short_name"
    //}

    req.period = this.period
    if (graph.type == 'lines' || graph.type == 'areas') {
      req.time_group = this.period.substring(4);
      req.period = graph.histoPeriod
    }
    req.stats_cpu = true
    req.stats_mem = true
    req.stats_net = true
    req.stats_io = true
    if (graph.criterion == 'stack_name') {
      req.filter_stack_name = graph.criterionValue
    } else if (graph.criterion == 'service_name') {
      req.filter_service_name = graph.criterionValue
    } else if (graph.criterion == 'container_id') {
      req.filter_container_id = graph.criterionValue
    } else if (graph.criterion == 'node_id') {
      req.filter_node_id = graph.criterionValue
    }
    let id = graph.id
    let newItem = new StatsRequestItem(id, req, graph.title)
    newItem.subscriberNumber=1
    this.requestMap[id]=newItem
    graph.requestId = id
    this.executeRequest(newItem)
    return id;
  }

  getCurrentData(graph : Graph) : GraphCurrentData[] {
    let item = this.requestMap[graph.id]
    if (!item) {
      return []
    }
    if (!item.currentResult) {
      return []
    }
    this.sortCurrentByField(item.currentResult, graph.field)
    if (graph.topNumber == 0 || graph.type == 'counterSquare' || graph.type == 'counterCircle') {
      return item.currentResult
    }
    return item.currentResult.slice(0, graph.topNumber)
  }

  sortCurrentByField(data : GraphCurrentData[], field : string) {
    data.sort((a, b) => {
      if (a.values[field] < b.values[field]) {
        return 1;
      }
      return -1
    })
  }

  getHistoricData(graph : Graph) : GraphHistoricAnswer {
    let item = this.requestMap[graph.id]
    if (!item) {
      return new GraphHistoricAnswer([], [])
    }
    if (!item.historicResult) {
      return new GraphHistoricAnswer([], [])
    }
    //let list = this.sortHistoricByField(item.historicResult, graph.field)
    let dateMap = {}
    let list : GraphHistoricData[] = []
    let names : string[] = []
    let nameMap = {}
    for (let dat of item.historicResult) {
      let pdata = dateMap[dat.sdate]
      if (!pdata) {
        pdata = new GraphHistoricData(dat.date)
        dateMap[dat.sdate] = pdata
        list.push(pdata)
      }
      let max = nameMap[dat.name]
      if (!max) {
        names.push(dat.name)
        nameMap[dat.name]=dat.values[graph.field]
      } else if (dat.values[graph.field] > max) {
        nameMap[dat.name]=dat.values[graph.field]
      }
      pdata.graphValues.push(dat.values[graph.field])
      pdata.graphValuesUnit.push(0) //graphValues and valuesUnit should have the same size
    }
    if (graph.topNumber > 0) {
      names = names.slice(0, graph.topNumber)
      for (let dat of list) {
        dat.graphValues = dat.graphValues.slice(0, graph.topNumber)
      }
    }
    return new GraphHistoricAnswer(names, list)
  }

  clear() {
    this.requestMap = {}
    this.graphs = []
    this.selected = this.notSelected
    this.nbGraph = 1
  }

  getData() : string {
    return JSON.stringify(this.graphs)
  }

  setData(data : string) {
    this.clear()
    let graphs = JSON.parse(data)
    this.nbGraph = 1
    for (let graph of graphs) {
      this.ascendingCompatibilityAdjustment(graph)
      this.nbGraph++
      graph.id="graph"+this.nbGraph
      this.graphs.push(graph)
      this.addRequest(graph)
    }
    this.menuService.onRefreshClicked.next()
  }

  ascendingCompatibilityAdjustment(graph : Graph) {
    //nothing to do for now
  }

  defaultDefaultDashboard() : string {
    return `[{"id":"graph2","x":6,"y":7,"width":316,"height":104,"type":"counterSquare","fields":[],"title":"Stacks CPU: ","border":true,"modeParameter":false,"topNumber":3,"alert":true,"alertMin":"30","alertMax":"100","criterion":"","criterionValue":"","stackedAreas":true,"legendNames":[],"legendColors":[],"containerAvg":false,"roundedBox":true,"object":"all","field":"cpu-usage","counterHorizontal":true,"requestId":"graph2"},{"id":"graph3","x":5,"y":120,"width":317,"height":103,"type":"counterSquare","fields":[],"title":"Stacks Mem: ","border":true,"modeParameter":false,"topNumber":3,"alert":true,"alertMin":"4GB","alertMax":"7GB","criterion":"","criterionValue":"","stackedAreas":true,"legendNames":[],"legendColors":[],"containerAvg":false,"roundedBox":true,"object":"all","field":"mem-usage","counterHorizontal":true,"requestId":"graph3"},{"id":"graph4","x":4,"y":234,"width":318,"height":110,"type":"counterSquare","fields":[],"title":"Stacks Net: ","border":true,"modeParameter":false,"topNumber":3,"alert":true,"alertMin":"","alertMax":"","criterion":"","criterionValue":"","stackedAreas":true,"legendNames":[],"legendColors":[],"containerAvg":false,"roundedBox":true,"object":"all","field":"net-total-bytes","counterHorizontal":true,"requestId":"graph4"},{"id":"graph5","x":2,"y":355,"width":319,"height":107,"type":"counterSquare","fields":[],"title":"Stacks IO: ","border":true,"modeParameter":false,"topNumber":3,"alert":true,"alertMin":"","alertMax":"","criterion":"","criterionValue":"","stackedAreas":true,"legendNames":[],"legendColors":[],"containerAvg":false,"roundedBox":true,"object":"all","field":"io-total","counterHorizontal":true,"requestId":"graph5"},{"id":"graph6","x":485,"y":22,"width":210,"height":54,"type":"counterSquare","fields":[],"title":"Stack number: ","border":true,"modeParameter":false,"topNumber":3,"alert":true,"alertMin":"","alertMax":"1","criterion":"","criterionValue":"","stackedAreas":true,"legendNames":[],"legendColors":[],"containerAvg":false,"roundedBox":true,"object":"stack","field":"number","counterHorizontal":true,"requestId":"graph6"},{"id":"graph7","x":484,"y":83,"width":211,"height":53,"type":"counterSquare","fields":[],"title":"Service number: ","border":true,"modeParameter":false,"topNumber":3,"alert":true,"alertMin":"","alertMax":"10","criterion":"","criterionValue":"","stackedAreas":true,"legendNames":[],"legendColors":[],"containerAvg":false,"roundedBox":true,"object":"service","field":"number","counterHorizontal":true,"requestId":"graph7"},{"id":"graph8","x":482,"y":145,"width":212,"height":60,"type":"counterSquare","fields":[],"title":"Container number: ","border":true,"modeParameter":false,"topNumber":3,"alert":true,"alertMin":"","alertMax":"14","criterion":"","criterionValue":"","stackedAreas":true,"legendNames":[],"legendColors":[],"containerAvg":false,"roundedBox":true,"object":"container","field":"number","counterHorizontal":true,"requestId":"graph8"},{"id":"graph9","x":735,"y":232,"width":559,"height":263,"type":"bubbles","fields":[],"title":"services","border":true,"modeParameter":false,"topNumber":5,"alert":false,"alertMin":"","alertMax":"","criterion":"","criterionValue":"","stackedAreas":true,"legendNames":[],"legendColors":[],"containerAvg":false,"roundedBox":true,"object":"service","field":"net-total-bytes","bubbleXField":"mem-usage","bubbleYField":"cpu-usage","bubbleScale":"medium","requestId":"graph9"},{"id":"graph10","x":742,"y":30,"width":173,"height":176,"type":"pie","fields":[],"title":"First 3 services CPU ","border":true,"modeParameter":false,"topNumber":3,"alert":false,"alertMin":"","alertMax":"","criterion":"","criterionValue":"","stackedAreas":true,"legendNames":[],"legendColors":[],"containerAvg":false,"roundedBox":true,"object":"service","field":"cpu-usage","requestId":"graph10"},{"id":"graph11","x":1109,"y":28,"width":183,"height":177,"type":"pie","fields":[],"title":"First 3 services Net","border":true,"modeParameter":false,"topNumber":3,"alert":false,"alertMin":"","alertMax":"","criterion":"","criterionValue":"","stackedAreas":true,"legendNames":[],"legendColors":[],"containerAvg":false,"roundedBox":true,"object":"service","field":"net-total-bytes","requestId":"graph11"},{"id":"graph12","x":924,"y":30,"width":175,"height":176,"type":"pie","fields":[],"title":"First 3 services Mem","border":true,"modeParameter":false,"topNumber":3,"alert":false,"alertMin":"","alertMax":"","criterion":"","criterionValue":"","stackedAreas":true,"legendNames":[],"legendColors":[],"containerAvg":false,"roundedBox":true,"object":"service","field":"mem-usage","requestId":"graph12"},{"id":"graph13","x":479,"y":266,"width":217,"height":177,"type":"legend","fields":[],"title":"Legend services","border":false,"modeParameter":false,"topNumber":3,"alert":false,"alertMin":"","alertMax":"","criterion":"","criterionValue":"","stackedAreas":true,"legendNames":[],"legendColors":[],"containerAvg":false,"roundedBox":true,"object":"service","transparentLegend":false},{"id":"graph14","x":336,"y":8,"width":102,"height":101,"type":"lines","fields":[],"title":"all","border":true,"modeParameter":false,"topNumber":3,"alert":false,"alertMin":"","alertMax":"","criterion":"","criterionValue":"","stackedAreas":true,"legendNames":[],"legendColors":[],"containerAvg":false,"roundedBox":true,"object":"all","field":"cpu-usage","histoPeriod":"now-30m","requestId":"graph14","yTitle":"cpu usage"},{"id":"graph15","x":335,"y":120,"width":100,"height":103,"type":"lines","fields":[],"title":"all","border":true,"modeParameter":false,"topNumber":3,"alert":false,"alertMin":"","alertMax":"","criterion":"","criterionValue":"","stackedAreas":true,"legendNames":[],"legendColors":[],"containerAvg":false,"roundedBox":true,"object":"all","field":"mem-usage","histoPeriod":"now-30m","requestId":"graph15","yTitle":"cpu usage"},{"id":"graph16","x":335,"y":235,"width":99,"height":108,"type":"lines","fields":[],"title":"all","border":true,"modeParameter":false,"topNumber":3,"alert":false,"alertMin":"","alertMax":"","criterion":"","criterionValue":"","stackedAreas":true,"legendNames":[],"legendColors":[],"containerAvg":false,"roundedBox":true,"object":"all","field":"net-total-bytes","histoPeriod":"now-30m","requestId":"graph16","yTitle":"cpu usage"},{"id":"graph17","x":332,"y":357,"width":101,"height":104,"type":"lines","fields":[],"title":"all","border":true,"modeParameter":false,"topNumber":3,"alert":false,"alertMin":"","alertMax":"","criterion":"","criterionValue":"","stackedAreas":true,"legendNames":[],"legendColors":[],"containerAvg":false,"roundedBox":true,"object":"all","field":"cpu-usage","histoPeriod":"now-30m","requestId":"graph17","yTitle":"cpu usage"}]`

  }

}
