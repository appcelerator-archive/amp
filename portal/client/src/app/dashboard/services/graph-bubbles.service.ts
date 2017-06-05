import { Injectable } from '@angular/core';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { DashboardService } from './dashboard.service'
import { Subject } from 'rxjs/Subject'
import { Graph } from '../../models/graph.model';
import * as d3 from 'd3';

@Injectable()
export class GraphBubbles {
  onNewData = new Subject();
  private margin: any = { top: 40, bottom: 30, left: 60, right: 20};
  private svg : any
  private xScale : any;
  private yScale : any;
  private sScale : any;
  private xAxis: any;
  private yAxis: any;
  private legend : any
  private focus : any
  private element: any
  private created = false
  private chart: any;
  private width: number;
  private height: number;
  data = []

  constructor(
    private httpService : HttpService,
    private menuService : MenuService,
    private dashboardService : DashboardService) { }

  destroy() {
    this.svg.selectAll("*").remove();
  }

  computeSize(graph : Graph) {
    this.margin.top = graph.height * 0.10
    this.margin.bottom = graph.height * 0.15
    this.margin.left = graph.width * 0.15
    this.margin.right = graph.width * 0.15
    this.width = graph.width - this.margin.left - this.margin.right;
    this.height = graph.height - this.margin.top - this.margin.bottom;
  }

  createGraph(graph : Graph, chartContainer : any) {
    this.element = chartContainer.nativeElement;

    this.computeSize(graph)
    this.svg = d3.select(this.element)
      .append('svg')
      .attr('width', 2000)//graph.width)
      .attr('height', 2000)//graph.height)

    this.created=true
    this.updateGraph(graph)
  }

  resizeGraph(graph : Graph) {
    if (!this.created) {
      return
    }
    this.computeSize(graph)
    this.updateGraph(graph)
  }

  updateGraph(graph : Graph) {
    this.data = this.dashboardService.getCurrentData(graph)
    if (this.data.length == 0) {
      return
    }

    this.svg.selectAll("*").remove();

    this.xScale = d3.scaleLinear()
      .range([0, this.width])
    this.yScale = d3.scaleLinear()
      .range([this.height, 0]);

    let xmax = d3.max(this.data, d => d.values[graph.bubbleXField])
    let ymax = d3.max(this.data, d => d.values[graph.bubbleYField])

    let xunit = this.dashboardService.computeUnit(graph.bubbleXField, xmax, "").unit
    let yunit = this.dashboardService.computeUnit(graph.bubbleYField, ymax, "").unit
    this.data = this.dashboardService.adjustCurrentXYDataToUnit(xunit, yunit, graph.bubbleXField, graph.bubbleYField, this.data)
    xmax = xmax / this.dashboardService.unitdivider(xunit)
    ymax = ymax / this.dashboardService.unitdivider(yunit)

    this.xScale.domain([0, xmax]);
    this.yScale.domain([0, ymax]);

    //let wwt = this.dashboardService.getTextWidth(graph.title, "10", "Arial")
    let fontSize = this.height/10
    let dx = this.margin.left
    let dy = this.margin.top

    this.xAxis = this.svg.append('g')
      .attr('class', 'axis axis-x')
      .attr('transform', "translate(" + [dx, this.height+dy] +")")
      .call(d3.axisBottom(this.xScale))
      .style("font-size", fontSize/2+'px')

    this.yAxis = this.svg.append('g')
      .attr('class', 'axis axis-y')
      .attr('transform', "translate(" + [dx, dy] +")")
      .call(d3.axisLeft(this.yScale))
      .style("font-size", fontSize/2+'px')

    let ethis=this

    let d=0
    let size = 400
    if (graph.bubbleScale == 'large') {
      size = size*2
    }
    if (graph.bubbleScale == 'small') {
      size = size/2
    }
    if (graph.field != 'none') {
      this.sScale = d3.scaleLinear().range([0, size]);
      this.sScale.domain([0, d3.max(this.data, (d) => { return d.values[graph.field]; })]);
    }
    for (let dat of this.data) {
      d++
      let x = this.xScale(dat.valueUnitx)+dx
      let y = this.yScale(dat.valueUnity)+dy
      let s = size
      if (graph.field != 'none') {
        s = Math.sqrt(this.sScale(dat.values[graph.field]))
      }
      this.svg.append('circle')
        .attr('class', 'circle')
        .attr('r', s)
        .attr("transform", "translate(" + [x, y] + ")")
        .style('fill', ethis.dashboardService.getObjectColor(graph, dat.group))
        .style("stroke", 'black')
    }

    if (graph.title) {
      let xt = -5
      let anchor = 'left'
      if (graph.centerTitle) {
        xt = (this.width)/2;
        anchor = 'middle'
      }
      this.svg.append("text")
       .attr("class", "wtitle")
       .attr("transform", "translate(" + [xt+dx,dy-this.margin.top] + ")")
       .attr("dy", "1em")
       .style("text-anchor", anchor)
       .style("font-size", fontSize+'px')
       .text(graph.title);
     }

     let yTitle = this.dashboardService.yTitleMap[graph.bubbleYField]
     //if (graph.yTitle) {
       this.svg.append("text")
         .attr("class", "y-title")
         .attr("transform", "rotate(-90)")
         .attr("y", dx - this.margin.left)
         .attr("x", dy - (this.height+this.margin.top+this.margin.bottom) / 2)
         .attr("dy", "1em")
         .style("text-anchor", "middle")
         .style("font-size", fontSize*2/3+'px')
         .text(yTitle+" ("+yunit+")");
      // }

     let xTitle = this.dashboardService.yTitleMap[graph.bubbleXField]
     //if (graph.yTitle) {
       this.svg.append("text")
         .attr("class", "y-title")
         .attr("y", this.height+this.margin.bottom+dy)
         .attr("x", dx+(this.width) / 2)
         .attr("dy", "-0.50em")
         .style("text-anchor", "middle")
         .style("font-size", fontSize*2/3+'px')
         .text(xTitle+" ("+xunit+")");
       //}
   }
}
