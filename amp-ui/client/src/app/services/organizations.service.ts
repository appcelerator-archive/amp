import { Injectable, Output, EventEmitter } from '@angular/core';
import { HttpService } from '../services/http.service';
import { MenuService } from '../services/menu.service';
import { Organization } from '../models/organization.model';
import { OrganizationMember } from '../models/organization-member.model';
import { Router } from '@angular/router';
import { Subject } from 'rxjs/Subject';

@Injectable()
export class OrganizationsService {
  onOrganizationLoaded = new Subject();
  onOrganizationError = new Subject();
  noOrganization : Organization = new Organization("", "")
  organizations : Organization[] = []
  currentOrganization = this.noOrganization
  @Output() onUserLogout = new EventEmitter<void>();

  constructor(private router : Router, private httpService : HttpService, private menuService : MenuService) {
    this.organizations.push(new Organization("appcelerator", "amp@appcelerator.fr"))
    this.organizations.push(new Organization("myOrganization", "amp@myorg.fr"))
  }

  match(item : Organization, value : string) : boolean {
    if (item.name && item.name.includes(value)) {
      return true
    }
    if (item.email && item.email.includes(value)) {
      return true
    }
    return false
  }

  setOrganization(org : Organization) {
    console.log(org.name)
    this.currentOrganization = org
  }

  edit() {
    this.menuService.setItemMenu("organization", "edit", "/amp/organizations/"+this.currentOrganization.name)
    //this.router.navigate(["/amp/organizations/", this.currentOrganization.name])
  }

}
