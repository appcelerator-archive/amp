import { Component, OnInit, OnDestroy, ViewChild, ElementRef, Input, ViewEncapsulation } from '@angular/core';
import * as d3 from 'd3';
import { MetricsService } from '../../services/metrics.service';
import { MenuService } from '../../../services/menu.service';
import { GraphHistoricData } from '../../../models/graph-historic-data.model';
import { GraphLine } from '../../models/graph-line.model';
import { Graph } from '../../../models/graph.model';

//style="height:parentdiv.offsetHeight;width:parentdiv.offsetWidth"
@Component({
  selector: 'app-graph-line',
  template: `<div class="d3-chart" #chart></div>`,
  styleUrls: ['./lines.component.css'],
  encapsulation: ViewEncapsulation.None
})

export class LinesComponent implements OnInit, OnDestroy {
  @ViewChild('chart') private chartContainer: ElementRef;
  @Input() private graph : Graph;
  private formatValue = d3.format(',.2f');
  private dateFormat = d3.timeFormat("%y-%m-%d %H:%M:%S")
  public selectedLine = -1
  private margin: any = { top: 40, bottom: 30, left: 60, right: 20};
  private lines : GraphLine[] = []
  private svg : any
  valuelines : any[]
  private x : any;
  private y : any;
  private xAxis: any;
  private yAxis: any;
  private data : GraphHistoricData[] = [];
  private legend : any
  private focus : any
  private element: any
  private created = false

  private chart: any;
  private width: number;
  private height: number;


  constructor(
    private metricsService : MetricsService,
    private menuService : MenuService) { }

  ngOnInit() {
    this.createGraph();
    this.metricsService.onNewData.subscribe(
      () => {
        //this.svg.remove();
        this.svg.selectAll("*").remove();
        this.updateGraph();
      }
    )
    this.menuService.onWindowResize.subscribe(
      (win) => {
        this.svg.selectAll("*").remove();
        this.resizeGraph()
      }
    );
  }

  ngOnDestroy() {
    this.svg.selectAll("*").remove();
    //this.metricsService.onNewData.unsubscribe()
  }

  createGraph() {
    // set the dimensions and margins of the graph
    this.element = this.chartContainer.nativeElement;
    this.width = this.graph.width - this.margin.left - this.margin.right;
    this.height = this.graph.height - this.margin.top - this.margin.bottom;
    //console.log("create: "+this.graph.title+": "+this.width+","+this.height)
    this.svg = d3.select(this.element)
      .append('svg')
        .attr('width',2000)// this.graph.width)
        .attr('height', 2000)//this.graph.height)
      .append("g")
        .attr("transform", "translate(" + this.margin.left + "," + this.margin.top + ")")
      .on('click', () => this.selectedClick());
    this.updateGraph()
    this.created=true
  }

  resizeGraph() {
    if (!this.created) {
      return
    }
    this.element = this.chartContainer.nativeElement;
    this.width = this.graph.width - this.margin.left - this.margin.right;
    this.height = this.graph.height - this.margin.top - this.margin.bottom;
    d3.select('svg')
      .attr('width', this.graph.width)
      .attr('height', this.graph.height)
    this.updateGraph()
  }

  updateGraph() {
    let ans = this.metricsService.getHistoricData(this.graph.fields, this.metricsService.object, this.metricsService.type)
    this.data = ans.data
    this.lines = ans.lines
    //console.log(this.lineRefs)

    this.chart = this.svg.append('g')
      .attr('class', 'lines')
      .attr('transform', `translate(${this.margin.left}, ${this.margin.top})`);

    this.x = d3.scaleTime().range([0, this.width]);
    this.y = d3.scaleLinear().range([this.height, 0]);

    // Scale the range of the data
    this.x.domain(d3.extent(this.data, (d) => { return d.date; }));
    let ymax=0
    for (let tmp of this.data) {
      for (let yy=0; yy<tmp.graphValues.length; yy++) {
        if (this.isVisible(this.lines[yy].name)) {
          if (tmp.graphValues[yy]>ymax) {
            ymax = tmp.graphValues[yy]
          }
        }
      }
    }
    this.y.domain([0, ymax])

    // define the lines
    this.valuelines = []
    //console.log(this.metricsService.lineVisibleMap)
    for (let ll=0; ll<this.lines.length; ll++) {
      if (!this.isVisible(this.lines[ll].name)) {
        this.valuelines.push(null)
      } else {
        this.valuelines.push(
          d3.line<GraphHistoricData>()
            .defined( d => { return d.graphValues[ll] !== undefined; })
            .x((d: GraphHistoricData) => { return this.x(d.date); })
            .y((d: GraphHistoricData) => { return this.y(d.graphValues[ll]); })
        )
        this.svg.append("path")
          .data([this.data])
          .attr("class", this.lines[ll].name+" line ")
          .style("stroke", this.metricsService.getColor(ll))
          .attr("d", this.valuelines[ll]);
      }
    }

    // add the X Axis
    if (this.width>80) {
      this.xAxis = this.svg.append("g")
        .attr("class", "axisx")
        .attr("transform", "translate(0," +  this.height + ")")
        .call(d3.axisBottom(this.x).ticks(5));
    }

    // add the Y Axis
    if (this.height>50) {
      this.yAxis = this.svg.append("g")
        .attr("class", "axisy")
        .call(d3.axisLeft(this.y));

      if (this.graph.yTitle != '') {
        this.svg.append("text")
          .attr("class", "y-title")
          .attr("transform", "rotate(-90)")
          .attr("y", 0 - this.margin.left)
          .attr("x", 0 - (this.height / 2))
          .attr("dy", "1em")
          .style("text-anchor", "middle")
          .text(this.graph.yTitle);
        }
    }

    if (this.graph.title != '') {
      this.svg.append("text")
       .attr("class", "wtitle")
       .attr("transform", "translate(-"+(this.margin.left-5)+",-"+(this.margin.top-10)+")")
       .style("text-anchor", "left")
       .text(this.graph.title);
    }

    this.svg.append("rect")
      .attr('width', this.width+this.margin.left+this.margin.right)
      .attr('height', this.height+this.margin.bottom)
      .attr("transform", "translate(-"+this.margin.left+",-"+(this.margin.top-15)+")")
      .attr('stroke', 'lightgrey')
      .style('fill', 'none')

    this.focus = this.svg.append('g')
      .attr('class', 'focus')
      .style('display', 'none');

    this.focus.append('circle')
      .attr('class', 'select-circle')
      .attr('r', 5)
      .style('fill', 'none')
      .style("stroke", 'black')

    this.svg.append('rect')
      .attr('class', 'overlay')
      .attr('width', this.width)
      .attr('height', this.height)
      .on('mouseover', () => this.focus.style('display', null))
      .on('mouseout', () => this.removeSelected)
      .on('mousemove', () => { this.mousemove(this.x, this.y) });

  }

  mousemove(x, y) {
    let pt = d3.mouse(d3.event.currentTarget)
    let x0 = x.invert(pt[0]);
    let y0 = y.invert(pt[1]);
    let dist
    let yy = 0
    let ptx
    let pty = 0
    let d
    for (let jj=0;jj<this.data.length; jj++) {
      let dn = this.data[jj]
      let xn = x(dn.date)
      for (let ii=0; ii < dn.graphValues.length; ii++) {
        if (this.isVisible(this.lines[ii].name)) {
          let yn = y(dn.graphValues[ii])
          let dist2 = (xn-pt[0])*(xn-pt[0]) + (yn-pt[1])*(yn-pt[1])
          if (dist === undefined || dist2 < dist) {
            dist = dist2
            pty = dn.graphValues[ii]
            ptx = dn.date
            d = dn
            yy = ii
          }
        }
      }
    }
    let epty = y(pty)
    let eptx = x(ptx)
    let dist2 = (eptx-pt[0])*(eptx-pt[0]) + (epty-pt[1])*(epty-pt[1])
    dist = Math.sqrt(dist2)
    if (dist<20) {
      this.removeSelected()
      this.selectedLine = yy
      this.focus.style('display', null)
      let valueLine = d3.line<GraphHistoricData>()
          .x((d: GraphHistoricData) => { return this.x(d.date); })
          .y((d: GraphHistoricData) => { return this.y(d.graphValues[yy]); })
      this.svg.append("path")
        .data([this.data])
        .attr("class", "selectedLine")
        .style("stroke", this.metricsService.getColor(yy))
        .attr("d", valueLine);
      this.focus.attr('transform', `translate(${eptx}, ${epty})`);
      this.svg.append("text")
         .attr("class", "info")
         .attr("transform", "translate(10,-10)")
         .style("text-anchor", "left")
         //.attr("font-size", "10")
         .style("fill", this.metricsService.getColor(yy))
         .text(this.lines[yy].displayedName+": "+this.dateFormat(ptx)+" -> "+this.formatValue(pty));
      this.focus.selectAll('rect').classed('toolTip', true)
        .style("left", 10)
        .style("top", 10)
    } else {
      this.removeSelected()
    }
  }

  removeSelected() {
    d3.select("path.selectedLine").remove();
    d3.select("text.info").remove();
    this.selectedLine = -1
    this.focus.style('display', 'none')
  }

  selectedClick() {
      if (this.metricsService.type == 'single' && this.metricsService.object!='global') {
      return
    }
    let nn = this.selectedLine
    //console.log("selected="+nn+" object="+this.object+" ref="+this.lineLabels[nn])
    if (nn>=0) {
      this.removeSelected()
      this.metricsService.route(this.lines[nn].name)
    }
  }

  isVisible(ref : string) : boolean {
    if (this.metricsService.lineVisibleMap[ref] === undefined) {
      return true
    }
    return this.metricsService.lineVisibleMap[ref]
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
