import { Component, OnInit, OnDestroy, ViewChild, ElementRef, Input, ViewEncapsulation } from '@angular/core';
import * as d3 from 'd3';
import { MetricsService } from '../../services/metrics.service';
import { GraphHistoricData } from '../../models/graph-historic-data.model';


@Component({
  selector: 'app-graph-line',
  template: `<div class="d3-chart" #chart></div>`,
  styleUrls: ['./lines.component.css'],
  encapsulation: ViewEncapsulation.None
})

export class LinesComponent implements OnInit, OnDestroy {
  @ViewChild('chart') private chartContainer: ElementRef;
  @Input() private fields: Array<string>;
  @Input() private yTitle: string;
  @Input() private title: string;
  bisectDate = d3.bisector((d : GraphHistoricData) => d.date).left;
  formatValue = d3.format(',.2f');
  //formatCurrency = d => `{this.formatValue(d)}`;
  private margin: any = { top: 20, bottom: 40, left: 50, right: 20};
  private svg : any
  valuelines : any[]
  private x : any;
  private y : any;
  private xAxis: any;
  private yAxis: any;
  private colors = ['steelblue', 'red']
  private data : GraphHistoricData[] = [];
  private legend : any
  private focus : any

  private chart: any;
  private width: number;
  private height: number;


  constructor(private metricsService : MetricsService) { }

  ngOnInit() {
    this.createGraph();
    this.metricsService.onNewData.subscribe(
      () => {
        this.updateGraph()
      }
    )
  }

  ngOnDestroy() {
    //this.metricsService.onNewData.unsubscribe()
  }

  createGraph() {
    this.data = this.metricsService.getHistoricData(this.fields)
    // set the dimensions and margins of the graph
    let element = this.chartContainer.nativeElement;
    this.width = element.offsetWidth - this.margin.left - this.margin.right;
    this.height = element.offsetHeight - this.margin.top - this.margin.bottom;
    let svg = d3.select(element)
      .append('svg')
        .attr('width', element.offsetWidth)
        .attr('height', element.offsetHeight)
      .append("g")
        .attr("transform", "translate(" + this.margin.left + "," + this.margin.top + ")");

    this.svg = svg
    this.chart = svg.append('g')
      .attr('class', 'lines')
      .attr('transform', `translate(${this.margin.left}, ${this.margin.top})`);

    this.x = d3.scaleTime().range([0, this.width]);
    this.y = d3.scaleLinear().range([this.height, 0]);

    // Scale the range of the data
    this.x.domain(d3.extent(this.data, (d) => { return d.date; }));
    this.y.domain([0, d3.max(this.data, (ds) => {
      return d3.max(ds.graphValues, (d) => { return d; }
    )})]);

    // define the lines
    this.valuelines = []
    for (let ll=0; ll<this.fields.length; ll++) {
      this.valuelines.push(
        d3.line<GraphHistoricData>()
          .x((d: GraphHistoricData) => { return this.x(d.date); })
          .y((d: GraphHistoricData) => { return this.y(d.graphValues[ll]); })
      )

      svg.append("path")
        .data([this.data])
        .attr("class", this.fields[ll]+" line")
        .style("stroke", this.colors[ll])
        .attr("d", this.valuelines[ll]);

    }

    // add the X Axis
    this.xAxis = svg.append("g")
      .attr("class", "axisx")
      .attr("transform", "translate(0," +  this.height + ")")
      .call(d3.axisBottom(this.x).ticks(5));

    // add the Y Axis
    this.yAxis = svg.append("g")
      .attr("class", "axisy")
      .call(d3.axisLeft(this.y));

    if (this.title != '') {
      svg.append("text")
       .attr("class", "title")
       .attr("transform", "translate(10,-10)")
       .style("text-anchor", "middle")
       .text(this.title);
    }

    if (this.yTitle != '') {
      svg.append("text")
        .attr("class", "y-title")
        .attr("transform", "rotate(-90)")
        .attr("y", 0 - this.margin.left)
        .attr("x",0 - (this.height / 2))
        .attr("dy", "1em")
        .style("text-anchor", "middle")
        .text(this.yTitle);
    }

    this.focus = this.svg.append('g')
      .attr('class', 'focus')
      .style('display', 'none');

    this.focus.append('circle')
      .attr('class', 'select-circle')
      .attr('r', 4)
      .style('fill', 'none')
      .style("stroke", 'black')

    this.focus.append('text')
      .attr('x', 9)
      .attr('dy', '.35em')
      .attr('stroke', 'red')

    svg.append("rect")
      .attr("class", "toolTip")
      .attr('width', 50)
      .attr('height', 20)
      .attr('stroke', 'black')
      .attr('fill', 'red')

    svg.append('rect')
      .attr('class', 'overlay')
      .attr('width', this.width)
      .attr('height', this.height)
      .on('mouseover', () => this.focus.style('display', null))
      .on('mouseout', () => this.focus.style('display', 'none'))
      .on('mousemove', () => { this.mousemove(this.x) });

  }

  updateGraph() {
    this.data = this.metricsService.getHistoricData(this.fields)
    //console.log(this.data)

    // Scale the range of the data again
    this.x.domain(d3.extent(this.data, (d) => { return d.date; }));
    this.y.domain([0, d3.max(this.data, (ds) => {
      return d3.max(ds.graphValues, (d) => { return d; }
    )})]);
    this.xAxis.transition().call(d3.axisBottom(this.x).ticks(5));
    this.yAxis.transition().call(d3.axisLeft(this.y));

    // Select the section we want to apply our changes to
    //let svg = d3.select("app-graph-line").transition();
    let svg = this.svg.transition();
    // Make the changes
    for (let ll=0; ll<this.fields.length; ll++) {
      svg.select("."+this.fields[ll])
        .duration(0)
        .attr("d", this.valuelines[ll](this.data));
    }
    /*
    svg.select("axisx") // change the x axis
      .duration(0)
      .call(this.xAxis);

    svg.select("axisy") // change the y axis
      .duration(0)
      .call(this.yAxis);
      */
  }

  mousemove(x) {
    let pt = d3.mouse(d3.event.currentTarget)
    let x0 = x.invert(pt[0]);
    let i = this.bisectDate(this.data, x0, 1);
    let d0 : GraphHistoricData = this.data[i - 1];
    let d1 : GraphHistoricData = this.data[i];
    let d = x0 - d0.date.getTime() > d1.date.getTime() - x0 ? d1 : d0;
    this.focus.attr('transform', `translate(${this.x(d.date)}, ${this.y(d.graphValues[0])})`);
    this.focus.select('text')
      .attr("class", "text-label")
      .text(this.formatValue(d.graphValues[0]));
    this.focus.selectAll('rect').classed('toolTip', true)
      .style("left", 10)
      .style("top", 10)
  }

}


/*
// add the X gridlines
this.xAxis = svg.append("g")
  .attr("class", "axisx grid")
  .attr("transform", "translate(0," + this.height + ")")
  .call(d3.axisBottom(this.x).tickSize(-this.height)) //.tickFormat("")
  .selectAll("text")
    .style("text-anchor", "end")
    .attr("dx", "-.8em")
    .attr("dy", ".15em")
    .attr("transform", "rotate(-55)");

// add the Y gridlines
this.yAxis = svg.append("g")
  .attr("class", "axisy grid")
  .call(d3.axisLeft(this.y).tickSize(-this.width)
    //.tickFormat()
)
*/
