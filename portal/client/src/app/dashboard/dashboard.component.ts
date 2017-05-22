import { Component, OnInit, OnDestroy } from '@angular/core';
import { MenuService } from '../services/menu.service';
import { DashboardService } from './services/dashboard.service';
import { ActivatedRoute } from '@angular/router';
import { AppWindow } from '../models/app-window.model';
import { Graph } from '../models/graph.model';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})

export class DashboardComponent implements OnInit, OnDestroy {
  dashboardName = "default"
  periodRefreshLabel = "30 seconds"
  graphPanelHeight = 250
  graphPanelWidth = 500


  constructor(
    private menuService : MenuService,
    public dashboardService : DashboardService,
    private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.menuService.setItemMenu('dashboard', 'View')
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

  resizeGraphs(win : AppWindow) {
    let cww = win.width-25-this.menuService.paddingLeftMenu
    let chh = win.height- 220;
    this.graphPanelHeight = chh
    this.graphPanelWidth = cww
  }

}
