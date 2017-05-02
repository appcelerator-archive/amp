import { Component, OnInit } from '@angular/core';
import { MenuService } from '../services/menu.service';

@Component({
  selector: 'app-logs',
  templateUrl: './logs.component.html',
  styleUrls: ['./logs.component.css']
})
export class LogsComponent implements OnInit {

  constructor(private menuService : MenuService) { }

  ngOnInit() {
    this.menuService.setItemMenu('logs', 'View')
  }

}
