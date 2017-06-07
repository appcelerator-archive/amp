import { Injectable } from '@angular/core';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { DashboardService } from './dashboard.service'
import { Subject } from 'rxjs/Subject'
import { Graph } from '../../models/graph.model';
import { GraphHistoricData } from '../../models/graph-historic-data.model';
import * as d3 from 'd3';

@Injectable()
export class GraphLines {
  onNewData = new Subject();
  private margin: any = { top: 40, bottom: 30, left: 60, right: 20};
  private svg : any
  private xScale : any;
  private yScale : any;
  private xAxis: any;
  private yAxis: any;
  private legend : any
  private focus : any
  private element: any
  private created = false
  private chart: any;
  private width: number;
  private height: number;
  private data : GraphHistoricData[] = []
  private names : string[] = []



  constructor(
    private httpService : HttpService,
    private menuService : MenuService,
    private dashboardService : DashboardService) { }


  destroy() {
    this.svg.selectAll("*").remove();
  }

  computeSize(graph : Graph) {
    this.margin.top = graph.height * 0.1
    this.margin.bottom = graph.height * 0.2
    this.margin.left = graph.width * 0.15
    this.margin.right = 10
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
    let ans = this.dashboardService.getHistoricData(graph)
    this.data = ans.data
    this.names = ans.names
    if (this.data.length == 0) {
      return
    }

    this.svg.selectAll("*").remove();

    this.xScale = d3.scaleTime().range([0, this.width]);
    this.yScale = d3.scaleLinear().range([this.height, 0]);

    let fontSize = this.height/10
    let dx = this.margin.left
    let dy = this.margin.top

    // Scale the range of the data
    this.xScale.domain(d3.extent(this.data, (d) => { return d.date; }));
    let ymax=0
    for (let tmp of this.data) {
      for (let yy=0; yy<tmp.graphValues.length; yy++) {
        if (tmp.graphValues[yy]>ymax) {
          ymax = tmp.graphValues[yy]
        }
      }
    }
    let yunit = this.dashboardService.computeUnit(graph.field, ymax, "").unit
    this.data = this.dashboardService.adjustHistoricDataToUnit(yunit, graph.field, this.data)
    ymax = ymax / this.dashboardService.unitdivider(yunit)
    let yDomain = [0, ymax];

    this.yScale.domain([0, ymax])
    for (let ll=0; ll<this.names.length; ll++) {
      let valueline = d3.line<GraphHistoricData>()
        .defined( d => { return d.graphValuesUnit[ll] !== undefined; })
        .x((d: GraphHistoricData) => { return this.xScale(d.date); })
        .y((d: GraphHistoricData) => { return this.yScale(d.graphValuesUnit[ll]); })

      this.svg.append("path")
        .data([this.data])
        .style("stroke", (d) => this.dashboardService.getObjectColor(graph, this.names[ll]))
        .attr("transform", "translate(" + [dx, dy] + ")")
        .style("fill", 'none')
        .style("stroke-width", 2)
        .attr("d", valueline);
    }

    // add the X Axis
    if (this.width>80) {
      this.xAxis = this.svg.append("g")
        .attr("class", "axisx")
        .attr("transform", "translate(" + [dx, this.height+dy] + ")")
        .style("font-size", fontSize/2+'px')
        .call(d3.axisBottom(this.xScale).ticks(5));
    }

    // add the Y Axis
    if (this.height>80) {
      this.yAxis = this.svg.append("g")
        .attr("class", "axisy")
        .attr("transform", "translate(" + [dx, dy] + ")")
        .style("font-size", fontSize/2+'px')
        .call(d3.axisLeft(this.yScale));

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

     graph.yTitle = this.dashboardService.yTitleMap[graph.field]
     if (graph.yTitle) {
       this.svg.append("text")
         .attr("class", "y-title")
         .attr("y", dx - this.margin.left)
         .attr("x", dy - (this.height+this.margin.top+this.margin.bottom) / 2)
         .attr("transform", "rotate(-90)")
         .attr("dy", "1em")
         .style("text-anchor", "middle")
         .style("font-size", fontSize/2+'px')
         .text(graph.yTitle+" ("+yunit+")");
       }
    }

    /*
    this.svg.append("rect")
      .attr('width', this.width)
      .attr('height', this.height)
      .attr("transform", "translate(" + [dx, dy] + ")")
      .attr('stroke', 'lightgrey')
      .style('fill', 'none')
    */

  }
}
