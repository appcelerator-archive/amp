import { Injectable } from '@angular/core';
import { Http, Headers, Response } from '@angular/http';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import { Subject } from 'rxjs/Subject';
import {Observable} from 'rxjs/Observable';
import { User } from '../models/user.model';
import { Team } from '../models/team.model';
import { Organization } from '../models/organization.model';
import { Member } from '../models/member.model';
import { TeamResource } from '../models/team-resource.model';
import { DockerStack } from '../models/docker-stack.model';


@Injectable()
export class HttpService {
  private token = ""
  onHttpError = new Subject();
  constructor(private http : Http) {}

  setToken(token : string) {
    this.token=token
  }

  private setHeaders() {
    var headers = new Headers
    headers.set('TokenKey', this.token)
    return headers
  }

  users() {
    return this.http.get("/api/v1/users", { headers: this.setHeaders() })
      .map((res : Response) => {
        const data = res.json()
        let list : User[] = []
        for (let item of data.users) {
          let user = new User(item.name, item.email, "User")
          user.verified = item.is_verified
          list.push(user)
        }
        return list
      })
  }

  createOrganization(org : Organization) {
    return this.http.post("/api/v1/organization/create", {name: org.name, email: org.email}, { headers: this.setHeaders() })
  }

  deleteOrganization(org : Organization) {
    return this.http.post("/api/v1/organization/remove", {data: org.name}, { headers: this.setHeaders() })
  }

  addUserToOrganization(org : Organization, member : Member) {
    return this.http.post("/api/v1/organization/user/add", {organization: org.name, name: member.userName}, { headers: this.setHeaders() })
  }

  removeUserFromOrganization(org : Organization, member : Member) {
    return this.http.post("/api/v1/organization/user/remove", {organization: org.name, name: member.userName}, { headers: this.setHeaders() })
  }

  createTeam(org : Organization, team : Team) {
    return this.http.post("/api/v1/team/create", {organization: org.name, name: team.name}, { headers: this.setHeaders() })
  }

  deleteTeam(org : Organization, team : Team) {
    return this.http.post("/api/v1/team/remove", {organization: org.name, name: team.name}, { headers: this.setHeaders() })
  }

  addUserToTeam(org : Organization, team : Team, member : Member) {
    return this.http.post("/api/v1/team/user/add", {organization: org.name, team: team.name, name: member.userName}, { headers: this.setHeaders() })
  }

  removeUserFromTeam(org : Organization, team : Team, member : Member) {
    return this.http.post("/api/v1/team/user/remove", {organization: org.name, team: team.name, name: member.userName}, { headers: this.setHeaders() })
  }

  userOrganization(user : User) {
    return this.http.post("/api/v1/user/organizations", {data: user.name}, { headers: this.setHeaders() })
      .map((res : Response) => {
        const data = res.json()
        //console.log("data")
        //console.log(data)
        let list : Organization[] = []
        for (let org of data.organizations) {
          let newOrg = new Organization(org.name, org.email)
          for (let mem of org.members) {
            newOrg.members.push(new Member(mem.name, mem.role))
          }
          for (let team of org.teams) {
            let newTeam = new Team(team.name)
            for (let mname of team.members) {
              newTeam.members.push(new Member(mname, 0))
            }
            newOrg.teams.push(newTeam)
          }
          list.push(newOrg)
        }
        console.log(list)
        return list
      })
  }

  login(user : User, pwd : string) {
    return this.http.post("/api/v1/login", {name: user.name, pwd: pwd}, { headers: this.setHeaders() });
  }

  stacks() {
    return this.http.get("/api/v1/stacks", { headers: this.setHeaders() })
    .map((res : Response) => {
      const data = res.json()
      let list : DockerStack[] = []
      if (data.stacks) {
        for (let item of data.stacks) {
          let stack = new DockerStack(
            item.stack.id,
            item.stack.name,
            item.service,
            item.stack.owner.name,
            item.stack.owner.type
          )
          list.push(stack)
        }
      }
      return list
    })
  }

  organizationRessources() {
    return this.http.get("/api/v1/...", { headers: this.setHeaders() })
    .map((res : Response) => {
      let list : Organization[] = []
      //
      return list
    })
  }

}
