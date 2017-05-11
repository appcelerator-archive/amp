import { Component, Input, OnInit, OnDestroy } from '@angular/core';
import { MenuService } from '../services/menu.service';
import { MetricsService } from './services/metrics.service';
import { StatsRequest } from './models/stats-request.model';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-dashboard',
  templateUrl: './metrics.component.html',
  styleUrls: ['./metrics.component.css']
})
export class MetricsComponent implements OnInit, OnDestroy {
    periodLabel = "last 15 min, every 30 sec"
    routeSub : any
    stackName = ""
    serviceId = ""
    containerId = ""

  constructor(
    private menuService : MenuService,
    private metricsService : MetricsService,
    private route: ActivatedRoute) { }

  ngOnInit() {
    this.menuService.setItemMenu('metrics', 'View')
    this.routeSub = this.route.params.subscribe(params => {
      this.stackName = params['stackName'];
      this.serviceId = params['serviceId'];
      this.containerId = params['containerId'];

      let req = new StatsRequest()
      req.stats_cpu = true
      req.stats_mem = true
      req.stats_io = true
      req.stats_net = true
      req.period = "now-15m"
      req.time_group = "30s"
      if (this.stackName) {
        req.filter_stack_name = this.stackName
        this.menuService.setItemMenu('metrics', 'View stack')
      }
      if (this.serviceId) {
        req.filter_service_id = this.serviceId
        this.menuService.setItemMenu('metrics', 'View service')
      }
      if (this.containerId) {
        req.filter_task_id = this.containerId
        this.menuService.setItemMenu('metrics', 'container')
      }
      this.metricsService.setHistoricRequest(req, 30)
    })
  }

  ngOnDestroy() {
    this.metricsService.cancelRequests()
  }

  setPeriod(period : string, group : string) {
    this.periodLabel = this.metricsService.setPeriod(period, group)
  }

}
