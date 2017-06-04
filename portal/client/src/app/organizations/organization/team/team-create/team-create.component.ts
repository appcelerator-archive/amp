import { Component, OnInit } from '@angular/core';
import { Organization } from '../../../../organizations/models/organization.model'
import { Team } from '../../../../organizations/models/team.model'
import { OrganizationsService } from '../../../../organizations/services/organizations.service'
import { NgForm } from '@angular/forms';
import { MenuService } from '../../../../services/menu.service';
import { ActivatedRoute } from '@angular/router';
import { HttpService } from '../../../../services/http.service';

@Component({
  selector: 'app-team-create',
  templateUrl: './team-create.component.html',
  styleUrls: ['./team-create.component.css']
})
export class TeamCreateComponent implements OnInit {
  team : Team = new Team("")
  organization : Organization
  routeSub : any
  messageError = ""

  constructor(
    private organizationsService : OrganizationsService,
    private menuService : MenuService,
    private route: ActivatedRoute,
    private httpService : HttpService) {
    }

  ngOnInit() {
    this.messageError = ""
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
    this.team.name = form.value.name
    this.organization.teams.push(this.team)
    this.menuService.waitingCursor(true)
    this.httpService.createTeam(this.organization.name, this.team.name).subscribe(
      () => {
        this.menuService.waitingCursor(false)
        this.menuService.navigate(['/amp', 'organizations', this.organization.name, 'team', this.team.name])
      },
      (err) => {
        this.menuService.waitingCursor(false)
        let error = err.json()
        this.messageError = error.error
        console.log(error)
      }
    )
  }

  cancel() {
    this.menuService.returnToPreviousPath()
  }

}
