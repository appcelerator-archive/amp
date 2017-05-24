import { Injectable } from '@angular/core';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { DashboardService } from './dashboard.service'
import { Subject } from 'rxjs/Subject'
import { Graph } from '../../models/graph.model';
import { GraphStats } from '../../models/graph-stats.model';
import * as d3 from 'd3';

@Injectable()
export class GraphCounter {
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
  colors : any

  constructor(
    private httpService : HttpService,
    private menuService : MenuService,
    private dashboardService : DashboardService) { }

  init(graph : Graph, chartContainer : any) {
    this.createGraph(graph, chartContainer);
    this.dashboardService.onNewData.subscribe(
      () => {
        this.updateGraph(graph);
      }
    )
    this.menuService.onWindowResize.subscribe(
      (win) => {
        this.svg.selectAll("*").remove();
        this.resizeGraph(graph, chartContainer)
      }
    );
  }

  destroy() {
    this.svg.selectAll("*").remove();
  }

  computeSize(graph : Graph) {
    this.margin.top = graph.height * 0.15
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
        .attr('width', graph.width)
        .attr('height', graph.height)
      .append("g")
        .attr("transform", "translate(" + this.margin.left + "," + this.margin.top + ")")

    this.created=true
    this.updateGraph(graph)
  }

  resizeGraph(graph : Graph, chartContainer : any) {
    if (!this.created) {
      return
    }
    this.element = chartContainer.nativeElement;
    this.computeSize(graph)
    //console.log("resize: "+graph.title+": "+this.width+","+this.height)
    d3.select('svg')
      .attr('width', graph.width)
      .attr('height', graph.height)
    d3.select("g").attr("transform", "translate(" + this.margin.left + "," + this.margin.top + ")")
    this.updateGraph(graph)
  }

  updateGraph(graph : Graph)
  {
    this.data = this.dashboardService.getData(graph)

    this.svg.selectAll("*").remove();
    let wwt = this.dashboardService.getTextWidth(graph.title, "10", "Arial")
    let fontSize = this.width/wwt*7*0.90;

    if (graph.title != '') {
      this.svg.append("text")
       .attr("class", "wtitle")
       .attr("transform", "translate("+this.width/2+","+this.height/2+")")
       .style("text-anchor", "middle")
       .attr("dy", "-.36em")
       .style("font-size", fontSize+'px')
       .text(graph.title);
     }

     if (this.data.length>0) {
       let val = Math.floor(this.data[0].values[graph.field])


       let dec = ".75em"
       if (graph.title == "") {
         dec = ".34em"
         wwt = this.dashboardService.getTextWidth(val, "10", "Arial")
         fontSize = this.width/wwt*7*0.90;
       }

       this.svg.append("text")
        .attr("class", "wtitle")
         .attr("transform", "translate("+this.width/2+","+this.height/2+")")
        .style("text-anchor", "middle")
        .attr("dy", dec)
        .style("font-size", fontSize+'px')
        .text(val);

        if (graph.alert) {
          let color="green"
          if (val>graph.alertMin) {
            color="orange"
            if (val>graph.alertMax) {
              color="red"
            }
          }
          this.svg.append("rect")
            .attr('width', this.width+this.margin.left+this.margin.right)
            .attr('height', this.height+this.margin.top+this.margin.bottom)
            .attr("transform", "translate(-"+this.margin.left+",-"+this.margin.top+")")
            .attr('stroke', 'lightgrey')
            .style('fill', color)
            .attr('fill-opacity', 0.4)
        }
    }
  }

}
