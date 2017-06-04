import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Organization } from '../organizations/models/organization.model';
import { MenuService } from '../services/menu.service';
import { OrganizationsService } from '../organizations/services/organizations.service';
import { UsersService } from '../services/users.service';
import { ListService } from '../services/list.service';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../services/http.service';

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
    private menuService : MenuService,
    private httpService : HttpService,
    private usersService : UsersService) {
      listService.setFilterFunction(organizationsService.match)
    }

  ngOnInit() {
    this.menuService.setItemMenu('organizations', 'List')
    this.organizationsService.onOrganizationsLoaded.subscribe(
      () => {
        this.listService.setData(this.organizationsService.organizations)
      }
    )
    this.menuService.onRefreshClicked.subscribe(
      () => {
        this.organizationsService.loadOrganizations(this.usersService.currentUser.name, true)
        this.usersService.loadUsers(true)
      }
    )
    this.organizationsService.loadOrganizations(this.usersService.currentUser.name, false)
  }

  removeOrganization() {
    if(confirm("Are you sure to delete the organization: "+this.organization.name)) {
      this.menuService.waitingCursor(true)
      this.httpService.deleteOrganization(this.organization.name).subscribe(
        () => {
          this.menuService.waitingCursor(false)
          let list : Organization[] = []
          for (let org of this.organizationsService.organizations) {
            if (org.name != this.organization.name) {
              list.push(org)
            }
          }
          this.organizationsService.organizations=list
          this.listService.setData(this.organizationsService.organizations)
        },
        (error) => {
          this.menuService.waitingCursor(false)
          console.log(error)
        }
      )
    }
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
