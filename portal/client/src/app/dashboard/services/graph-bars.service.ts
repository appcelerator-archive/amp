import { Injectable } from '@angular/core';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { DashboardService } from './dashboard.service'
import { Subject } from 'rxjs/Subject'
import { Graph } from '../../models/graph.model';
import * as d3 from 'd3';

@Injectable()
export class GraphBars {
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

  updateGraph(graph : Graph) {
    this.data = this.dashboardService.getData(graph.requestId)
    let xDomain = this.data.map(d => d.group);
    let yDomain = [0, d3.max(this.data, d => d.values[graph.field])];

    this.svg.selectAll("*").remove();

    this.xScale = d3.scaleBand()
      .range([0, this.width])
      .padding(0.1);
    this.yScale = d3.scaleLinear()
      .range([this.height, 0]);

    this.xScale.domain(this.data.map((d) => { return d.group; }));
    this.yScale.domain([0, d3.max(this.data, (d) => { return d.values[graph.field]; })]);

    //let wwt = this.dashboardService.getTextWidth(graph.title, "10", "Arial")
    let fontSize = this.height/10

    this.colors = d3.scaleOrdinal()
        .range(["#6F257F", "#CA0D59"]);


    this.xAxis = this.svg.append('g')
      .attr('class', 'axis axis-x')
      .attr('transform', `translate(${0}, ${this.height})`)
      .call(d3.axisBottom(this.xScale))
      .style("font-size", fontSize*2/3+'px');

    this.yAxis = this.svg.append('g')
      .attr('class', 'axis axis-y')
      .call(d3.axisLeft(this.yScale))
      .style("font-size", fontSize*2/3+'px')

    let ethis=this

    this.svg.selectAll(".bar")
      .data(this.data)
      .enter().append("rect")
      .attr("class", "bar")
      .attr("x", (d) => { return ethis.xScale(d[0]) })
      .attr("width", ethis.xScale.bandwidth())
      .attr("y", (d) => { return ethis.yScale(d[1]); })
      .attr("height", (d) => { return ethis.height - ethis.yScale(d.values[graph.field]) })
      .attr("fill", (d) => { return this.colors(d.group); })

    this.svg.append("rect")
      .attr('width', this.width)
      .attr('height', this.height)
      .attr('stroke', 'lightgrey')
      .style('fill', 'none')

    if (graph.title == '') {
      graph.title = graph.object
    }
    if (graph.title != '') {
      this.svg.append("text")
       .attr("class", "wtitle")
       .attr("transform", "translate(-5,-"+this.margin.top*0.1+")")
       .style("text-anchor", "left")
       .style("font-size", fontSize+'px')
       .text(graph.title);
     }

     graph.yTitle = this.dashboardService.yTitleMap[graph.field]
     if (graph.yTitle != '') {
       this.svg.append("text")
         .attr("class", "y-title")
         .attr("transform", "rotate(-90)")
         .attr("y", 0 - this.margin.left)
         .attr("x", 0 - (this.height / 2))
         .attr("dy", "1em")
         .style("text-anchor", "middle")
         .style("font-size", fontSize*2/3+'px')
         .text(graph.yTitle);
       }
   }
}
