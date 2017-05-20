import { Component, OnInit } from '@angular/core';
import { MenuService } from '../services/menu.service';
import { ListService } from '../services/list.service';
import { Node } from './models/node.model';
import { HttpService } from '../services/http.service';
import { NodesService } from './services/nodes.service';

@Component({
  selector: 'app-nodes',
  templateUrl: './nodes.component.html',
  styleUrls: ['./nodes.component.css'],
  providers: [ ListService ]
})
export class NodesComponent implements OnInit {
  message : ""

  constructor(
    public listService : ListService,
    public nodesService : NodesService,
    private menuService : MenuService,
    private httpService : HttpService,) {
      listService.setFilterFunction(this.nodesService.match)
    }

  ngOnInit() {
    this.menuService.setItemMenu('nodes', 'List')
    this.nodesService.onNodesLoaded.subscribe(
      () => {
        this.listService.setData(this.nodesService.nodes)
      }
    )
    this.menuService.onRefreshClicked.subscribe(
      () => {
      this.nodesService.loadNodes(true)
      }
    )
    this.nodesService.loadNodes(false)
  }

  getColor(node : Node) {
    if (node.reachability == 'unreachable' || node.status == 'disconnected' ) {
      return 'red'
    } else if (node.reachability == 'unknown' || node.status == 'unknown') {
      return 'orange'
    } else if (node.status == 'down') {
      return 'lightgrey';
    } else if (node.role == 'manager' || node.role == 'leader') {
      return 'limegreen'
    }
    return 'green'
  }

  metrics(id : string) {
    this.menuService.navigate(['/amp', 'metrics', 'node', 'single', id])
  }

  logs(id : string) {
    this.menuService.navigate(['/amp', 'logs', 'node', id])
  }

}
