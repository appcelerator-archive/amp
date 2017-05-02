import { Injectable } from '@angular/core';
import { Http, Headers, Response } from '@angular/http';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import { Subject } from 'rxjs/Subject';
import {Observable} from 'rxjs/Observable';
import { User } from '../models/user.model';
import { DockerStack } from '../models/docker-stack.model';
import { Organization } from '../models/organization.model';

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

  userOrganization(userName : string) {
    return this.http.get("/api/v1/account/users/organization/"+userName, { headers: this.setHeaders() })
    .map((res : Response) => {
      const data = res.json()
      let list : Organization[] = []
      if (data.organizations) {
        for (let item of data.organization) {
          let orga = new Organization(
            '',
            ''
          )
          list.push(orga)
        }
      }
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

}
