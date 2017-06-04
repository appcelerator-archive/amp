import { Component, OnInit } from '@angular/core';
import { Organization } from '../../organizations/models/organization.model'
import { OrganizationsService } from '../../organizations/services/organizations.service'
import { NgForm } from '@angular/forms';
import { MenuService } from '../../services/menu.service'
import { HttpService } from '../../services/http.service'

@Component({
  selector: 'app-organization-create',
  templateUrl: './organization-create.component.html',
  styleUrls: ['./organization-create.component.css']
})
export class OrganizationCreateComponent implements OnInit {
  organization : Organization = new Organization("", "")

  constructor(
    private organizationsService : OrganizationsService,
    public menuService : MenuService,
    private httpService : HttpService) {
    }

  ngOnInit() {
    this.menuService.setItemMenu('organization', 'Create')
  }

  create(form : NgForm) {
    this.organization.name = form.value.name
    this.organization.email = form.value.email
    this.menuService.waitingCursor(true)
    this.httpService.createOrganization(this.organization).subscribe(
      () => {
        this.menuService.waitingCursor(false)
        this.organizationsService.organizations.push(this.organization)
        this.menuService.navigate(['/amp', 'organizations', this.organization.name])
      },
      (error) => {
        this.menuService.waitingCursor(false)
        console.log(error)
      }
    )
  }

  cancel() {
    this.menuService.returnToPreviousPath()
  }

}
