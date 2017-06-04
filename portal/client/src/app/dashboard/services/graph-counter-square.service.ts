import { Injectable } from '@angular/core';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { DashboardService } from './dashboard.service'
import { Subject } from 'rxjs/Subject'
import { Graph } from '../../models/graph.model';
import { GraphCurrentData } from '../../models/graph-current-data.model';
import * as d3 from 'd3';

@Injectable()
export class GraphCounterSquare {
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
    this.margin.top = 0
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
    if (this.data.length == 0) {
      return
    }

    this.svg.selectAll("*").remove();

    let dx = this.margin.left
    let dy = this.margin.top
    let val = this.data.length
    let uval = val
    let sval = ""+val
    if (this.data.length>0) {
      if (graph.field != 'number') {
        val = 0
        for (let dat of this.data) {
          val += dat.values[graph.field]
        }
        let unit=this.dashboardService.computeUnit(graph.field, val, "")
        uval = unit.val
        sval = unit.sval
      }
    }

    let title = graph.title+" "+sval
    if (title == " ") {
      return
    }

    this.svg.append("text")
     .attr("id", "title")
     .attr("class", "wtitle")
     .style("text-anchor", "middle")
     .text(title);

    let padding = 10;
    let titleBox : any
    this.svg.select("#title").each(function(d, i) {
      titleBox = this.getBBox();
    });
    if (titleBox.width == 0 || titleBox.height == 0) {
      return
    }

    let fontSize = this.computeFontSize(titleBox, padding)
    let dty = this.computeDty(titleBox, fontSize)

    this.svg.select("#title")
        .attr("transform", "translate("+ [(this.width-padding)/2+dx, (this.height-padding)/2+dy+dty] + ")")
        .style("font-size", fontSize+'px')

    let color="green"
    let alertMin = this.getRealNumber(graph.alertMin)
    let alertMax = this.getRealNumber(graph.alertMax)
    if (graph.alert) {
      if (!graph.alertMax ) {
        if (val < alertMin) {
          color="orange"
        } else if (val < alertMin/2) {
          color="red"
        }
      } else if (!graph.alertMin) {
        if (val < alertMax) {
          color="orange"
        } else if (val < alertMax/2) {
          color="red"
        }
      } else if (alertMin < alertMax) {
        if (val>=alertMin) {
          color="orange"
          if (val>=alertMax) {
            color="red"
          }
        }
      } else {
        if (val<=alertMin) {
          color="orange"
          if (val<=alertMax) {
            color="red"
          }
        }
      }
      this.svg.append("rect")
        .attr('width', this.width+this.margin.left+this.margin.right)
        .attr('height', this.height+this.margin.top+this.margin.bottom)
        .attr("transform", "translate(" + [0,0] +")")
        .attr('stroke', 'lightgrey')
        .style('fill', color)
        .attr('fill-opacity', 0.4)
        .attr("rx", 10)
        .attr("ry", 10)
    }

  }

  getRealNumber(val : string) : number {
    if (!val) {
      return 0
    }
    let ret = parseInt(val)
    if (val.length<=2) {
      return ret
    }
    if (val.substring(val.length-2)=='KB') return ret * 1024;
    if (val.substring(val.length-2)=='MB') return ret * 1048576;
    if (val.substring(val.length-2)=='GB') return ret * 1073741824;
    return ret
  }

  //box1 mandatory, box2 optionnal
  computeFontSize(box1 : any, padding : number) : number {
    let nn : number
    let dd : number
    nn = Math.max(this.width - padding, this.height - padding)
    dd = Math.max(box1.width, box1.height)
    //console.log(nn+","+dd)
    return nn / dd *12
  }

  computeDty(box : any, fontSize : number) : number {
    return fontSize / 12 * box.height / 3;
  }


}
