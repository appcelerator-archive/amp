
import { Component, HostListener, OnInit, OnDestroy, Input, ElementRef, ViewChild} from '@angular/core';
import { Graph } from '../../models/graph.model';
import { DashboardService } from '../services/dashboard.service'
import { MenuService } from '../../services/menu.service';

@Component({
  selector: 'app-dgraph-editor',
  templateUrl: "./dgraph-editor.component.html",
  styleUrls: ['./dgraph-editor.component.css'],
})

export class DGraphEditorComponent implements OnInit, OnDestroy {
  @Input() public graph : Graph;

  constructor(
    public dashboardService : DashboardService) {
  }

  ngOnInit() {
    this.dashboardService.onGraphSelect.subscribe(
      (graph) => {

      }
    )
  }

  ngOnDestroy() {
  }

}
