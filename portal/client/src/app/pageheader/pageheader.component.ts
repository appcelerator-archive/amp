import { Component, OnInit } from '@angular/core';
import { MenuService } from '../services/menu.service';
import { UsersService } from '../services/users.service';
import { OrganizationsService } from '../organizations/services/organizations.service';
import { User } from '../models/user.model';
import { Organization } from '../organizations/models/organization.model';

@Component({
  selector: 'app-pageheader',
  templateUrl: './pageheader.component.html',
  styleUrls: ['./pageheader.component.css']
})

export class PageheaderComponent implements OnInit {
  menuTitle = "title"
  menuItem = "item"
  menuOver = ""


  constructor(
    public menuService : MenuService,
    public usersService : UsersService,
    public organizationsService : OrganizationsService) {}

  ngOnInit() {
  }

  createOrganization() {
    this.menuService.navigate(['/amp', 'organizations', 'create'])
  }

}
