import { Injectable } from '@angular/core';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { DashboardService } from './dashboard.service'
import { Subject } from 'rxjs/Subject'
import { Graph } from '../../models/graph.model';
import { GraphCurrentData } from '../../models/graph-current-data.model';
import * as d3 from 'd3';

@Injectable()
export class GraphCounter {
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
    this.margin.top = graph.height * 0
    this.margin.bottom = 0
    this.margin.left = 0
    this.margin.right = 0
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

  updateGraph(graph : Graph)
  {
    this.data = this.dashboardService.getCurrentData(graph)

    this.svg.selectAll("*").remove();

    let dx = this.margin.left
    let dy = this.margin.top
    let title = graph.title
    let val = this.data.length
    let sval = ""+val
    let dec = "0.9em"
    let dect = "-.02em"
    if (this.data.length>0) {
      if (graph.field != 'number') {
        val = 0
        for (let dat of this.data) {
          //console.log(dat)
          val += dat.values[graph.field]
        }
        let unit=this.dashboardService.computeUnit(graph, Math.floor(val))
        val = unit.val
        sval = unit.sval
      }
    }
    let fontCoef = 0.7
    if (graph.counterHorizontal) {
      title+=" "+sval
      dect = ".36em"
      fontCoef = 0.8
    }
    let wwt = this.dashboardService.getTextWidth(title, "10", "Arial")
    if (title == "") {
      dec = ".34em"
      wwt = this.dashboardService.getTextWidth(sval, "10", "Arial")
    }
    let fontSize = this.width/wwt*10*fontCoef;

    if (title != '') {
      this.svg.append("text")
       .attr("class", "wtitle")
       .attr("transform", "translate("+ [this.width/2+dx, this.height/2+dy] + ")")
       .style("text-anchor", "middle")
       .attr("dy", dect)
       .style("font-size", fontSize+'px')
       .text(title);
    }

    if (!graph.counterHorizontal) {
      this.svg.append("text")
        .attr("class", "wtitle")
         .attr("transform", "translate(" + [this.width/2+dx, this.height/2+dy] + ")")
        .style("text-anchor", "middle")
        .attr("dy", dec)
        .style("font-size", fontSize+'px')
        .text(sval);
    }

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
        .attr("transform", "translate(" + [0,0] +")")
        .attr('stroke', 'lightgrey')
        .style('fill', color)
        .attr('fill-opacity', 0.4)
    }
  }

}
