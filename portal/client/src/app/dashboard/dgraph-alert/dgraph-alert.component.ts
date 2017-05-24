import { Component, OnInit, Input } from '@angular/core';
import { Graph } from '../../models/graph.model';
import { DashboardService } from '../services/dashboard.service'


@Component({
  selector: 'app-dgraph-alert',
  templateUrl: './dgraph-alert.component.html',
  styleUrls: ['./dgraph-alert.component.css']
})
export class DgraphAlertComponent implements OnInit {
  @Input() public graph : Graph;
  slider : any
  max: number=1000
  min: number=0


  constructor(public dashboardService: DashboardService) { }

  ngOnInit() {
  }


}
