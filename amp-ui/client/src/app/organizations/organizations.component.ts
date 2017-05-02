import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Organization } from '../models/organization.model';
import { MenuService } from '../services/menu.service';
import { OrganizationsService } from '../services/organizations.service';
import { ListService } from '../services/list.service';
import { Observable } from 'rxjs/Observable';

@Component({
  selector: 'app-organizations',
  templateUrl: './organizations.component.html',
  styleUrls: ['./organizations.component.css'],
  providers: [ ListService ]
})
export class OrganizationsComponent implements OnInit {
  noOrganization = new Organization("", "")
  organization : Organization = this.noOrganization


  constructor(
    private route : ActivatedRoute,
    public organizationsService : OrganizationsService,
    public listService : ListService,
    private menuService : MenuService) {
      listService.setFilterFunction(organizationsService.match)
    }

  ngOnInit() {
    this.menuService.setItemMenu('organizations', 'List')
    this.listService.setData(this.organizationsService.organizations)
  }

  removeOrganization() {
    let list : Organization[] = []
    for (let org of this.organizationsService.organizations) {
      if (org.name != this.organization.name) {
        list.push(org)
      }
    }
    this.organizationsService.organizations=list
    this.listService.setData(this.organizationsService.organizations)
  }

  selectOrganization(org : Organization) {
    this.organization = org
  }

  createOrganization() {
    this.menuService.navigate(['/amp', 'organizations', 'create'])
  }

  editOrganization(org : Organization) {
    this.menuService.navigate(['/amp', 'organizations', org.name])
  }

}
