import { Component, OnInit } from '@angular/core';
import { MenuService } from '../services/menu.service';

@Component({
  selector: 'app-swarms',
  templateUrl: './swarms.component.html',
  styleUrls: ['./swarms.component.css']
})
export class SwarmsComponent implements OnInit {

  constructor(private menuService : MenuService) { }

  ngOnInit() {
    this.menuService.setItemMenu('swarms', 'List')
  }

}
