import { Injectable } from '@angular/core';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { Subject } from 'rxjs/Subject'
import { Graph } from '../../models/graph.model';
import * as d3 from 'd3';

@Injectable()
export class GraphLines {
  onNewData = new Subject();
  private margin: any = { top: 40, bottom: 30, left: 60, right: 20};
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
    private menuService : MenuService) { }

  init(graph : Graph, chartContainer : any) {
    this.createGraph(graph, chartContainer);
    this.resizeGraph(graph, chartContainer);
  }

  destroy() {
    this.svg.selectAll("*").remove();
  }

  createGraph(graph : Graph, chartContainer : any) {
    // set the dimensions and margins of the graph
    this.element = chartContainer.nativeElement;
    //console.log("create parent: "+this.element.offsetWidth+","+this.element.offsetHeight)
    //this.width = this.element.offsetWidth - this.margin.left - this.margin.right;
    //this.height = this.element.offsetHeight - this.margin.top - this.margin.bottom;
    this.width = graph.width - this.margin.left - this.margin.right;
    this.height = graph.height - this.margin.top - this.margin.bottom;
    //console.log("create: "+this.graph.title+": "+this.width+","+this.height)
    this.svg = d3.select(this.element)
      .append('svg')
        //.attr('width', this.element.offsetWidth)
        //.attr('height', this.element.offsetHeight)
        .attr('width',2000)// this.graph.width)
        .attr('height', 2000)//this.graph.height)
      .append("g")
        .attr("transform", "translate(" + this.margin.left + "," + this.margin.top + ")")
    //this.updateGraph()
    this.created=true
  }

  clearGraph() {
    this.svg.selectAll("*").remove();
  }

  resizeGraph(graph : Graph, chartContainer : any) {
    if (!this.created) {
      return
    }
    this.element = chartContainer.nativeElement;
    //console.log("resize parent: "+this.element.offsetWidth+","+this.element.offsetHeight)
    //this.width = this.element.offsetWidth - this.margin.left - this.margin.right;
    //this.height = this.element.offsetHeight - this.margin.top - this.margin.bottom;
    this.width = graph.width - this.margin.left - this.margin.right;
    this.height = graph.height - this.margin.top - this.margin.bottom;
    console.log("resize: "+graph.title+": "+this.width+","+this.height)
    d3.select('svg')
      //.attr('width', this.element.offsetWidth)
      //.attr('height', this.element.offsetHeight)
      .attr('width', graph.width)
      .attr('height', graph.height)
    //d3.select("g").attr("transform", "translate(" + this.margin.left + "," + this.margin.top + ")")
    this.updateGraph(graph)
  }

  updateGraph(graph : Graph) {
    this.chart = this.svg.append('g')
      .attr('transform', `translate(${this.margin.left}, ${this.margin.top})`);

    if (graph.title != '') {
      this.svg.append("text")
       .attr("class", "wtitle")
       .attr("transform", "translate(-"+(this.margin.left-5)+",-"+(this.margin.top-10)+")")
       .style("text-anchor", "left")
       .text(graph.title);
    }
  }

}
