import { Injectable } from '@angular/core';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { DashboardService } from './dashboard.service'
import { ColorsService } from './colors.service'
import { GraphColor } from '../../dashboard/models/graph-color.model'
import { Subject } from 'rxjs/Subject'
import { Graph } from '../../models/graph.model';
import * as d3 from 'd3';

@Injectable()
export class GraphLegend {
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
    private dashboardService : DashboardService,
    private colorsService : ColorsService) { }


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
    if (graph.transparentLegend) {
      this.addClass("graph-transparent")
    } else {
      this.removeClass("graph-transparent")
    }

    let fontSize = this.height/7
    let dx = this.margin.left
    let dy = this.margin.top

    this.svg.on('click', () => this.titleClick());
    //this.svg.attr('fill-opacity', 0.4)
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
    //let colorObject = this.dashboardService.graphObjectColorMap[graph.object]
    let colorList : GraphColor[] = this.colorsService.getColors(graph.object, graph.legendGraphId)
    if (colorList.length ==0) {
      return
    }

    let yy = dy*1.2
    let dh = this.height/colorList.length
    for (let col of colorList) {
        this.svg.append("rect")
          .attr('width', dh*.90)
          .attr('height', dh*.90)
          .attr("transform", "translate(" + [dx, yy] + ")")
          .attr('stroke', 'lightgrey')
          .style('fill', col.color)

        this.svg.append("text")
         .attr("class", "wtitle")
         .attr("transform", "translate(" + [dx+dh, yy+dh*.65] + ")")
         .style("text-anchor", "start")
         .style("font-size", fontSize*.75+'px')
         .text(col.name);
        yy+=dh
    }
  }

  titleClick() {
  }

  addClass(name : string) {
    if (this.element.className.indexOf(name)>=0) {
      return
    }
    this.element.className += " " + name
  }

  removeClass(name : string) {
    if (this.element.className.indexOf(name)<0) {
      return
    }
    //warning: works only for this component
    this.element.className = this.element.className.replace(name, "")
  }

}
