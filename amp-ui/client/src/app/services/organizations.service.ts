import { Injectable, Output, EventEmitter } from '@angular/core';
import { HttpService } from '../services/http.service';
import { Organization } from '../models/organization.model';
import { Router } from '@angular/router';
import { Subject } from 'rxjs/Subject'

@Injectable()
export class OrganizationsService {
  onOrganizationLoaded = new Subject();
  onOrganizationError = new Subject();
  allOrganizations : Organization = new Organization("All", "")
  organizations : Organization[] = [ this.allOrganizations ]
  currentOrganization = this.allOrganizations
  @Output() onUserLogout = new EventEmitter<void>();

  constructor(private router : Router, private httpService : HttpService) {}

  match(item : Organization, value : string) : boolean {
    if (item.name.includes(value)) {
      return true
    }
    if (item.email.includes(value)) {
      return true
    }
    return false
  }

}
