import { Injectable, Output, EventEmitter } from '@angular/core';
import { HttpService } from '../services/http.service';
import { MenuService } from '../services/menu.service';
import { Organization } from '../models/organization.model';
import { Router } from '@angular/router';

@Injectable()
export class OrganizationsService {
  noOrganization : Organization = new Organization("", "")
  organizations : Organization[] = []
  currentOrganization = this.noOrganization
  @Output() onUserLogout = new EventEmitter<void>();

  constructor(private router : Router, private httpService : HttpService, private menuService : MenuService) {
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
    this.menuService.navigate(["/amp/organizations/", this.currentOrganization.name])
  }

}
