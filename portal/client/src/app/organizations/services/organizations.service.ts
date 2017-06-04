import { Injectable, Output, EventEmitter } from '@angular/core';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { Organization } from '../../organizations/models/organization.model';
import { Member } from '../../organizations/models/member.model';
import { OrganizationResource } from '../../organizations/models/organization-resource.model';
import { Router } from '@angular/router';
import { Subject } from 'rxjs/Subject'

@Injectable()
export class OrganizationsService {
  noOrganization : Organization = new Organization("", "")
  organizations : Organization[] = []
  currentOrganization = this.noOrganization
  onOrganizationsLoaded = new Subject();
  onOrganizationsError = new Subject();
  @Output() onUserLogout = new EventEmitter<void>();
  currentLoadedUser = ""

  constructor(
    private router : Router,
    private httpService : HttpService,
    private menuService : MenuService) {
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

  loadOrganizations(userName : string, refresh : boolean) {
    if (!refresh && this.currentLoadedUser == userName) {
      this.onOrganizationsLoaded.next()
      return
    }
    this.httpService.userOrganization(userName).subscribe(
      data => {
        this.organizations = data
        this.currentLoadedUser = userName
        this.onOrganizationsLoaded.next()
      },
      error => {
        this.onOrganizationsError.next(error)
        console.log(error)
      }
    )
  }

  setCurrentOrganization(org : Organization) {
    this.httpService.switchOrganization(org.name).subscribe(
      (rep) => {
        let data = rep.json()
        let token = data.auth
        this.httpService.setToken(token)
        localStorage.setItem('token', JSON.stringify({ token: token }));
        this.currentOrganization = org
        this.httpService.organizationRessources().subscribe(
          (data) => {
            org.resources = data
            for (let team of org.teams) {
              for (let res of team.resources) {
                for (let ref of data) {
                  if (ref.id == res.id) {
                    res.name = ref.name
                    res.type = ref.type
                    break;
                  }
                }
              }
            }
          },
          (err) => {
            console.log(err)
          }
        )
      }
    )
  }

  edit() {
    this.menuService.navigate(["/amp", "organizations", this.currentOrganization.name])
  }

  //to be refactor with associative array
  getAllNoMembers(members : Member[]) {
    let list : Member [] = []
    for (let user of this.currentOrganization.members) {
      let found = false
      for (let member of members) {
        if (member.userName == user.userName) {
          found= true
          break;
        }
      }
      if (!found) {
        list.push(new Member(user.userName, undefined))
      }
    }
    return list
  }

  //to be refactor with associative array
  getAllNoResources(resources : OrganizationResource[]) {
    let list : OrganizationResource [] = []
    for (let ref of this.currentOrganization.resources) {
      let found = false
      for (let res of resources) {
        if (ref.id == res.id) {
          found= true
          break;
        }
      }
      if (!found) {
        list.push(new OrganizationResource(ref.id, ref.type, ref.name))
      }
    }
    return list
  }

}
