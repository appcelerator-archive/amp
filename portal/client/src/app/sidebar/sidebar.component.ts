import { Component, OnInit, EventEmitter, Output } from '@angular/core';
import { MenuService } from '../services/menu.service'
import { OrganizationsService } from '../organizations/services/organizations.service'
import { SwarmsService } from '../services/swarms.service'
import * as $ from 'jquery';

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.component.html',
  styleUrls: ['./sidebar.component.css']
})
export class SidebarComponent implements OnInit {
  @Output() onMenu = new EventEmitter<string>();
  sidebarDisplay = "normal";

  constructor(
    public menuService : MenuService,
    public organizationsService : OrganizationsService,
    public swarmsService : SwarmsService) { }

  ngOnInit() {
    this.sidebarDisplay = localStorage.getItem('sidebar')
    if (!this.sidebarDisplay) {
      this.sidebarDisplay = 'normal'
    }
    this.resize()
  }

  minimize() {
    if (this.sidebarDisplay == 'normal') {
      this.sidebarDisplay = 'mini'
    } else {
      this.sidebarDisplay = 'normal'
    }
    localStorage.setItem('sidebar', this.sidebarDisplay);
    this.resize()
  }

  resize() {
    if (this.sidebarDisplay == 'mini') {
      this.menuService.paddingLeftMenu=70
      $('.sidebar-list').width(70)
    } else {
        this.menuService.paddingLeftMenu=250
        $('.sidebar-list').width(250)
    }
    this.menuService.onWindowResize.next(this.menuService.appWindow)
  }

}
