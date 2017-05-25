import { Component, OnInit, OnDestroy, ElementRef, ViewChild, Renderer2 } from '@angular/core';
import { MenuService } from '../services/menu.service';
import { DashboardService } from './services/dashboard.service';
import { ActivatedRoute } from '@angular/router';
import { AppWindow } from '../models/app-window.model';
import { Graph } from '../models/graph.model';
import { NgForm } from '@angular/forms';
import * as $ from 'jquery';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})

export class DashboardComponent implements OnInit, OnDestroy {
  @ViewChild('container') private container: ElementRef;
  dashboardName = "default"
  periodRefreshLabel = "30 seconds"
  periodLabel = '2 min'
  graphPanelHeight = 250
  graphPanelWidth = 500
  offsetTop0 : number
  offsetLeft0 : number
  graph00 : Graph = new Graph("graph00", 0, 0, 10, 10, 'text', undefined)
  public dialogx = 0
  public dialogy = 0
  public dialogHidden = true
  public dialogMode = ""

  constructor(
    private menuService : MenuService,
    public dashboardService : DashboardService,
    private route: ActivatedRoute,
    private elementRef : ElementRef,
    private renderer : Renderer2) {
      this.graph00.border = false
  }

  ngOnInit() {
    this.menuService.setItemMenu('dashboard', 'View')
    this.offsetTop0 = this.container.nativeElement.getBoundingClientRect().top
    this.offsetLeft0 = this.container.nativeElement.getBoundingClientRect().left
    this.resizeGraphs(this.menuService.appWindow)
    this.menuService.onWindowResize.subscribe(
      (win) => {
        this.resizeGraphs(win)
      }
    )
  }

  ngOnDestroy() {
    this.dashboardService.cancelRequests()
  }

  onMouseUp($event) {
    if ($event.target.className == 'panel-body') {
      this.dashboardService.selected = this.dashboardService.notSelected;
    }
  }

  addGraph(type : string) {
    let offtop = this.container.nativeElement.getBoundingClientRect().top - this.offsetTop0
    let offleft = this.container.nativeElement.getBoundingClientRect().left - this.offsetLeft0
    this.dashboardService.addGraph(type, offtop, offleft)
  }

  toggleEditor() {
    let offtop = this.container.nativeElement.getBoundingClientRect().top - this.offsetTop0
    let offleft = this.container.nativeElement.getBoundingClientRect().left - this.offsetLeft0
    this.dashboardService.toggleEditor(offtop, offleft)
  }

  resizeGraphs(win : AppWindow) {
    let cww = win.width-25-this.menuService.paddingLeftMenu
    let chh = win.height- 210;
    this.graphPanelHeight = chh
    this.graphPanelWidth = cww
  }

  setRefreshPeriod(refresh : number, label : string) {
    this.periodRefreshLabel = label
    this.dashboardService.setRefreshPeriod(refresh)
  }

  setPeriod(period : string, label : string) {
    this.periodLabel = label
    this.dashboardService.setPeriod(period)
  }

  moveDialog() {
    let offtop = this.container.nativeElement.getBoundingClientRect().top - this.offsetTop0
    let offleft = this.container.nativeElement.getBoundingClientRect().left - this.offsetLeft0
    let ww = this.container.nativeElement.getBoundingClientRect().width
    this.dialogx = offleft + ww/2
    this.dialogy = offtop + 100
  }

  load() {

  }

  save() {
    
  }

}
