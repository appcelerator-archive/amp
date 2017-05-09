import { Component, OnInit } from '@angular/core';
import { MenuService } from '../services/menu.service';

@Component({
  selector: 'app-nodes',
  templateUrl: './nodes.component.html',
  styleUrls: ['./nodes.component.css']
})
export class NodesComponent implements OnInit {

  constructor(private menuService : MenuService) { }

  ngOnInit() {
    this.menuService.setItemMenu('nodes', 'List')
  }

}
