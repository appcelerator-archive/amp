import { Component, OnInit } from '@angular/core';
import { Organization } from '../../models/organization.model'
import { OrganizationsService } from '../../services/organizations.service'
import { NgForm } from '@angular/forms';
import { MenuService } from '../../services/menu.service'

@Component({
  selector: 'app-organization-create',
  templateUrl: './organization-create.component.html',
  styleUrls: ['./organization-create.component.css']
})
export class OrganizationCreateComponent implements OnInit {
  organization : Organization = new Organization("", "")

  constructor(
    private organizationsService : OrganizationsService,
    private menuService : MenuService) {
    }

  ngOnInit() {
    this.menuService.setItemMenu('organization', 'Create')
  }

  create(form : NgForm) {
    this.organization.name = form.value.name
    this.organization.email = form.value.email
    this.organizationsService.organizations.push(this.organization)
    this.menuService.navigate(['/amp', 'organizations', this.organization.name])
  }

  cancel() {
    this.menuService.navigate(['/amp', 'organizations'])
  }

}
