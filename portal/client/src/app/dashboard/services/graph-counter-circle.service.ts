import { Injectable } from '@angular/core';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { DashboardService } from './dashboard.service'
import { Subject } from 'rxjs/Subject'
import { Graph } from '../../models/graph.model';
import { GraphCurrentData } from '../../models/graph-current-data.model';
import * as d3 from 'd3';

@Injectable()
export class GraphCounterCircle {
  onNewData = new Subject();
  private margin: any = { top: 0, bottom: 0, left: 0, right: 0};
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
  data = []
  colors : any

  constructor(
    private httpService : HttpService,
    private menuService : MenuService,
    private dashboardService : DashboardService) { }

  destroy() {
    this.svg.selectAll("*").remove();
  }

  computeSize(graph : Graph) {
    this.margin.top = Math.floor(graph.height * 0.1)
    this.margin.bottom = 10
    this.margin.left = 10
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
    let tmpdata = this.dashboardService.getCurrentData(graph)
    if (tmpdata.length == 0) {
      return
    }

    this.svg.selectAll("*").remove();

    let title = graph.title
    let val = this.data.length
    let sval =  ""
    if (this.data.length>0) {
      if (graph.field != 'number') {
        val = 0
        for (let dat of this.data) {
          val += dat.values[graph.field]
        }
        let unit=this.dashboardService.computeUnit(graph.field, Math.floor(val), "")
        val = unit.val
        sval = unit.sval
      }
    }

    let values0 : { [name:string]: number; } = {}
    values0[graph.field] = 0
    let values : { [name:string]: number; } = {}
    values[graph.field] = val
    this.data = [
      new GraphCurrentData(graph.object, values0),
      new GraphCurrentData(graph.object, values)
    ]

    let xDomain = this.data.map(d => d.group);
    let yDomain = [0, d3.max(this.data, d => d.values[graph.field])];

    let fontSize = this.height/10
    let dx = this.margin.left
    let dy = this.margin.top

    this.svg.selectAll("*").remove();

    let arcs = d3.pie<GraphCurrentData>()
      .value((d) => {
          let val = d.values[graph.field];
          let format = this.dashboardService.computeUnit(graph.field, val,"")
          return format.val
        }
      )
      (this.data)


    let arc = d3.arc()
      .outerRadius(Math.min(this.height,this.width)/2)
      .innerRadius(Math.min(this.height,this.width)/4)
      .padAngle(0.03)
      .cornerRadius(8)

    let pieG = this.svg.selectAll("g")
      .data([this.data])
      .enter()
      .append("g")
      .attr("transform", "translate("+[this.width/2+dx, this.height/2+dy]+")")

    let block = pieG.selectAll(".arc")
      .data(arcs)

    var newBlock = block.enter().append("g").classed("arc", true)

    let athis=this
    newBlock.append("path")
      .attr("d", arc)
      .attr("id", function(d, i) { return "arc-" + i })
      .attr("stroke", "gray")
      .attr("fill", function(d,i){ return athis.dashboardService.getObjectColor(graph, d.data.group) })

    this.svg.append("text")
     .attr("class", "wtitle")
     .attr("transform", "translate("+[this.width/2+dx,this.height/2+dy]+")")
     .style("text-anchor", "middle")
     .attr("dy", ".36em")
     .style("font-size", fontSize/2+'px')
     .text(graph.field);
   /*
   this.svg.append("text")
    .attr("class", "wtitle")
    .attr("transform", "translate("+[this.width/2+dx,this.height/2+dy] +")")
    .style("text-anchor", "middle")
    .style("font-size", fontSize/2+'px')
    .attr("dy", ".95em")
    .text(this.dashboardService.unit[graph.field]);
  */

/*
    newBlock.append("text")
      .attr("transform", function(d) {
        d.outerRadius = 100;
        return "translate(" + arc.centroid(d) + ")";
      })
      .style("text-anchor", "middle")
      .attr("dy", ".35em")
      .style("font-size", fontSize/2+'px')
      .text(function(d) {
        let val = d.data.values[graph.field]
        let format = athis.dashboardService.computeUnit(graph, val)
        return format.sval
      });
*/


    if (graph.title) {
      this.svg.append("text")
       .attr("class", "wtitle")
       .attr("transform", "translate("+[this.width/2+dx,-this.margin.top*0.5+dy]+")")
       .style("text-anchor", "middle")
       .attr("dy", ".36em")
       .style("font-size", fontSize+'px')
       .text(graph.title);
     }
   }

 }
