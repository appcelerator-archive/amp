import { Injectable } from '@angular/core';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { DashboardService } from './dashboard.service'
import { Subject } from 'rxjs/Subject'
import { Graph } from '../../models/graph.model';
import { GraphCurrentData } from '../../models/graph-current-data.model';
import * as d3 from 'd3';

@Injectable()
export class GraphPie {
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
  data = []

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
      //.append("g")


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
    let xDomain = this.data.map(d => d.group);
    let yDomain = [0, d3.max(this.data, d => d.values[graph.field])];

    let fontSize = this.height/10
    let dx = this.margin.left
    let dy = this.margin.top

    this.svg.selectAll("*").remove();

    let arcs = d3.pie<GraphCurrentData>()
      .value((d) => { return d.values[graph.field]; })
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
      .attr("fill", function(d,i){ return athis.dashboardService.getObjectColor(graph.object, d.data.group) })

    if (this.data.length == 1) {
      this.svg.append("text")
       .attr("class", "wtitle")
       .attr("transform", "translate("+[this.width/2+dx,this.height/2+dy]+")")
       .style("text-anchor", "middle")
       .style("font-size", fontSize+'px')
       .text(this.data[0].group);
     this.svg.append("text")
      .attr("class", "wtitle")
      .attr("transform", "translate("+[this.width/2+dx,this.height/2+dy] +")")
      .style("text-anchor", "middle")
      .style("font-size", fontSize+'px')
      .attr("dy", ".95em")
      .text(Math.floor(this.data[0].values[graph.field]));
    } else {
      /*
      newBlock.append("text")
        .attr("transform", function(d) {
          d.outerRadius = 100;
          return "translate(" + arc.centroid(d) + ")";
        })
        .style("text-anchor", "middle")
        .attr("dy", ".35em")
        .style("font-size", fontSize/2+'px')
        .text(function(d) { return d.data.group; });
      */
      newBlock.append("text")
        .attr("transform", function(d) {
          d.outerRadius = 100;
          return "translate(" + arc.centroid(d) + ")";
        })
        .style("text-anchor", "middle")
        .attr("dy", ".35em")
        .style("font-size", fontSize/2+'px')
        .text(function(d) { return Math.floor(d.data.values[graph.field]); });
    }

    if (graph.title) {
      this.svg.append("text")
       .attr("class", "wtitle")
       .attr("transform", "translate("+[this.width/2+dx,-this.margin.top*0.5+dy]+")")
       .style("text-anchor", "middle")
       .attr("dy", ".36em")
       .style("font-size", fontSize+'px')
       .text(graph.title);
     }

     /*
     this.svg.append("rect")
       .attr('width', this.width)
       .attr('height', this.height)
       .attr("transform", "translate("+[dx,dy]+")")
       .attr('stroke', 'lightgrey')
       .style('fill', 'none')
    */

  }
}
