import { Injectable } from '@angular/core';
import { GraphColor } from '../models/graph-color.model'
import { Graph } from '../../models/graph.model'
import * as d3 from 'd3';

class ObjectColors {
  private index: number
  private refColor: string[]
  private object: string
  private colorList: GraphColor[]
  private nameMap: {}
  private graphMap: {}
  private colorMap : { [name:string]: GraphColor; } = {}

  constructor(object: string, refColor: string[]) {
    this.index = 0
    this.refColor =refColor
    this.object = object
    this.colorList = []
    this.colorMap = {}
    this.graphMap = {}
    this.nameMap = {}
  }

  getColor(name: string, graphId: string) : string {
    this.graphMap[name+'-'+graphId] = true
    let col = this.colorMap[name]
    if (col) {
      let exist = this.nameMap[name]
      if (!exist) {
        this.nameMap[name] = "."
        this.colorList.push(col)
      }
      return col.color;
    }
    if (this.index >= this.refColor.length) {
      this.index = 10
    }
    col = new GraphColor(name, graphId, this.refColor[this.index])
    this.index++
    this.colorMap[name] = col
    this.graphMap[name+'-'+graphId] = true
    this.nameMap[name] = "."
    this.colorList.push(col)
    return col.color
  }

  getColorList(graphId : string) {
    let list : GraphColor[] = []
    for (let col of this.colorList) {
      if (!graphId || this.graphMap[col.name+'-'+graphId]) {
        list.push(col)
      }
    }
    return list
  }

  refresh() {
    this.colorList = []
    this.nameMap = {}
    this.graphMap = {}
    this.index = 0
  }
}

export class ColorsService {
  private defaultColor = 'magenta'
  private refColors : string[] = ['#396AB1', '#DA7C30', '#3E9651', '#CC2529', '#535154', '#6B4C9A', '#922428', '#948B3D']
  private objectColorsMap : { [name:string]: ObjectColors; } = {}

  constructor() {
    for (let i=0;i<100;i++) {
      this.refColors.push(d3.interpolateCool(Math.random()))
    }
    this.objectColorsMap['stack'] = new ObjectColors('stack', this.refColors)
    this.objectColorsMap['service'] = new ObjectColors('service', this.refColors)
    this.objectColorsMap['container'] = new ObjectColors('container', this.refColors)
    this.objectColorsMap['node'] = new ObjectColors('node', this.refColors)
  }

  public getColor(graph : Graph, name: string) {
    let objectColors = this.objectColorsMap[graph.object];
    if (!objectColors) {
      return this.defaultColor;
    }
    return objectColors.getColor(name, graph.id)
  }

  public getColors(object: string, graphId: string): GraphColor[] {
    let objectColors = this.objectColorsMap[object];
    if (!objectColors) {
      return [];
    }
    return objectColors.getColorList(graphId)
  }

  public refresh() {
    for (let key in this.objectColorsMap) {
      let objectColors = this.objectColorsMap[key]
      objectColors.refresh()
    }
  }



}
