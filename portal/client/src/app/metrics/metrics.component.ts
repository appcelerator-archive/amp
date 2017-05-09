import { Component, OnInit, OnDestroy } from '@angular/core';
import { MenuService } from '../services/menu.service';
import { MetricsService } from './services/metrics.service';
import { StatsRequest } from './models/stats-request.model';

@Component({
  selector: 'app-metrics',
  templateUrl: './metrics.component.html',
  styleUrls: ['./metrics.component.css']
})
export class MetricsComponent implements OnInit, OnDestroy {
    periodLabel = "last 15 min, every 30 sec"

  constructor(
    private menuService : MenuService,
    private metricsService : MetricsService) { }

  ngOnInit() {
    this.menuService.setItemMenu('metrics', 'View')
    let req = new StatsRequest()
    req.stats_cpu = true
    req.stats_mem = true
    req.stats_io = true
    req.stats_net = true
    req.period = "now-15m"
    req.time_group = "30s"
    this.metricsService.setHistoricRequest(req, 30)
  }

  ngOnDestroy() {
    this.metricsService.cancelRequests()
  }

  setPeriod(period : string, group : string) {
    this.periodLabel = this.metricsService.setPeriod(period, group)
  }

}
