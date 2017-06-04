import { Injectable } from '@angular/core';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { DashboardService } from './dashboard.service'
import { Subject } from 'rxjs/Subject'
import { Graph } from '../../models/graph.model';
import * as d3 from 'd3';

@Injectable()
export class GraphText {
  onNewData = new Subject();
  private margin: any = { top: 0, bottom: 0, left: 30, right: 30};
  private svg : any
  private x : any;
  private y : any;
  private xAxis: any;
  private yAxis: any;
  private legend : any
  private focus : any
  private element: any
  private created = false
  private chart: any;
  private width: number;
  private height: number;

  constructor(
    private httpService : HttpService,
    private menuService : MenuService,
    private dashboardService : DashboardService) { }


  destroy() {
    this.svg.selectAll("*").remove();
  }

  computeSize(graph : Graph) {
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
    this.svg.selectAll("*").remove();
    if (graph.title) {
      let wwt = this.dashboardService.getTextWidth(graph.title, "10", "Arial")
      this.svg.append("text")
        .text(graph.title)
        .style("text-anchor", "middle")
        .attr("transform", "translate("+[this.width/2+this.margin.left, this.height/2+this.margin.top]+")")
        .style("font-size",this.width/wwt*10*0.90 +"px")
        //.style("font-size", (this.width / this.getComputedTextLength() *24) + "px"; })
        .attr("dy", ".35em");
    }
  }

}
