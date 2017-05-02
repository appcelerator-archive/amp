import { Component, OnInit } from '@angular/core';
import { Organization } from '../../../../models/organization.model'
import { Team } from '../../../../models/team.model'
import { OrganizationsService } from '../../../../services/organizations.service'
import { NgForm } from '@angular/forms';
import { MenuService } from '../../../../services/menu.service'
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-team-create',
  templateUrl: './team-create.component.html',
  styleUrls: ['./team-create.component.css']
})
export class TeamCreateComponent implements OnInit {
  team : Team = new Team("")
  organization : Organization
  routeSub : any

  constructor(
    private organizationsService : OrganizationsService,
    private menuService : MenuService,
    private route: ActivatedRoute) {
    }

  ngOnInit() {
    this.menuService.setItemMenu('organization', 'Team creation')
    this.routeSub = this.route.params.subscribe(params => {
      let name = params['orgName'];
      for (let org of this.organizationsService.organizations) {
        if (org.name == name) {
          this.organization = org
        }
      }
    })
  }

  create(form : NgForm) {
    if (form.valid) {
      this.team.name = form.value.name
      this.organization.teams.push(this.team)
      this.menuService.navigate(['amp', 'organizations', this.organization.name, 'team', this.team.name])
    }
  }

  cancel() {
    this.menuService.navigate(['amp', 'organizations'])
  }

}
