import { Component, OnInit, OnDestroy, ElementRef, ViewChild, Renderer2 } from '@angular/core';
import { MenuService } from '../services/menu.service';
import { HttpService } from '../services/http.service';
import { DashboardService } from './services/dashboard.service';
import { ActivatedRoute } from '@angular/router';
import { AppWindow } from '../models/app-window.model';
import { Graph } from '../models/graph.model';
import { Dashboard } from './models/dashboard.model';
import { NgForm } from '@angular/forms';
import * as $ from 'jquery';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})

export class DashboardComponent implements OnInit, OnDestroy {
  @ViewChild('container') private container: ElementRef;
  saveAsDashboardName = ""
  periodRefreshLabel = "30 secs"
  periodLabel = '2 min'
  graphPanelHeight = 250
  graphPanelWidth = 500
  offsetTop0 : number
  offsetLeft0 : number
  graph00 : Graph = new Graph("graph00", 0, 0, 10, 10, 'text', undefined)
  dialogGraph : Graph = new Graph("dialog", 0, 0, 10, 10, 'text', undefined)
  dashboards : Dashboard[] = []
  private noDashboard = new Dashboard("", "", "")
  public currentDashboard : Dashboard = this.noDashboard
  public selected : Dashboard = this.noDashboard
  public messageError = ""
  public message = ""
  public dialogx = 0
  public dialogy = 0
  public dialogHidden = true
  public dialogMode = ""

  constructor(
    private menuService : MenuService,
    public dashboardService : DashboardService,
    private httpService : HttpService,
    private route: ActivatedRoute,
    private elementRef : ElementRef,
    private renderer : Renderer2) {
      this.graph00.border = false
  }

  ngOnInit() {
    this.dashboardService.isVisible = true
    this.menuService.setItemMenu('dashboard', 'View')
    this.offsetTop0 = this.container.nativeElement.getBoundingClientRect().top
    this.offsetLeft0 = this.container.nativeElement.getBoundingClientRect().left
    this.resizeGraphs(this.menuService.appWindow)
    this.menuService.onWindowResize.subscribe(
      (win) => {
        this.resizeGraphs(win)
      }
    )
    if (this.dashboardService.graphs.length == 0) {
      let dashboardId = localStorage.getItem('dashboard')
      if (dashboardId) {
        this.openOneDashboard(dashboardId)
      } else {
        this.openOneDashboard("")
      }
    }
  }

  ngOnDestroy() {
    this.dashboardService.isVisible = false
    this.dashboardService.cancelRequests()
  }

  onMouseUp($event) {
    if (!this.dashboardService.showEditor) {
      if ($event.target.className == 'panel-body') {
        this.dashboardService.selected = this.dashboardService.notSelected;
      }
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
    let chh = win.height- 190;
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
    this.message = ""
    this.dialogx = this.menuService.appWindow.width/2 - 250/2;
    //this.dialogy = this.menuService.appWindow.height/2
    this.dialogy = 180
  }

  create() {
    this.currentDashboard = this.noDashboard
    this.dashboardService.clear()
  }

  save() {
    if (this.currentDashboard.id == "") {
      this.saveAs()
      return
    }
    let data = this.dashboardService.getData()
    this.httpService.updateDashboard(this.currentDashboard.id, data).subscribe(
      () => {
        this.dialogHidden = true
      },
      (err) => {
        let error = err.json()
        this.setMessage(error.error)
      }
    )
  }

  saveAs() {
    this.moveDialog()
    this.dialogMode = "SaveAs"
    this.dialogHidden = false
  }

  open() {
    this.httpService.listDashboard().subscribe(
      (data) => {
        console.log(data)
        this.dashboards = data
        //console.log(this.dashboards)
        this.moveDialog()
        this.dialogMode = "Open"
        this.dialogHidden = false
        //console.log(data)
      },
      (err) => {
        let error = err.json()
        this.message=error.error
      }
    )
  }

  saveNewDashboard() {
    let data = this.dashboardService.getData()
    this.httpService.createDashboard(this.saveAsDashboardName, data).subscribe(
      (id) => {
        this.dialogHidden = true
        this.currentDashboard = new Dashboard(id, this.saveAsDashboardName, '')
        localStorage.setItem('dashboard', id);
      },
      (err) => {
        let error = err.json()
        this.message = error.error
      }
    )
  }

  openOneDashboard(id : string) {
    this.httpService.getDashboard(id).subscribe(
      (dashboard) => {
        this.currentDashboard = dashboard
        this.dialogHidden = true
        this.dashboardService.setData(dashboard.data)
        localStorage.setItem('dashboard', dashboard.id);
      },
      (err) => {
        console.log(err)
        let error = err.json()
        this.message = error.error
        let dashboard = new Dashboard("", "default", this.dashboardService.defaultDefaultDashboard())
        this.currentDashboard = dashboard
        this.dialogHidden = true
        this.dashboardService.setData(dashboard.data)
        localStorage.setItem('dashboard', dashboard.id);
      }
    )
  }

  delete() {
    this.httpService.removeDashboard(this.currentDashboard.id).subscribe(
      () => {
        this.dashboardService.clear()
      },
      (err) => {
        console.log(err)
        let error = err.json()
        this.setMessage(error.error)
      }
    )
  }

  setMessage(mes : string) {
    this.messageError = mes
    setTimeout(() => { this.messageError=""}, 3000)
  }

}
