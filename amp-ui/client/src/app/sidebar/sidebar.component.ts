import { Component, OnInit, EventEmitter, Output } from '@angular/core';
import { MenuService } from '../services/menu.service'
import { OrganizationsService } from '../services/organizations.service'
import { SwarmsService } from '../services/swarms.service'

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.component.html',
  styleUrls: ['./sidebar.component.css']
})
export class SidebarComponent implements OnInit {
  @Output() onMenu = new EventEmitter<string>();

  constructor(
    public menuService : MenuService,
    public organizationsService : OrganizationsService,
    public swarmsService : SwarmsService) { }

  ngOnInit() {
  }

}
