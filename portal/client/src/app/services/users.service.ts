import { Injectable, Output, EventEmitter } from '@angular/core';
import { HttpService } from '../services/http.service';
import { OrganizationsService } from '../organizations/services/organizations.service';
import { User } from '../models/user.model';
import { Member } from '../organizations/models/member.model';
import { Organization } from '../organizations/models/organization.model';
import { Subject } from 'rxjs/Subject'
import { MenuService } from '../services/menu.service'

@Injectable()
export class UsersService {
  onUsersLoaded = new Subject();
  onUsersError = new Subject();
  users : User[] = []
  noLoginUser = new User('not signin','','')
  currentUser = this.noLoginUser
  @Output() onUserLogout = new EventEmitter<void>();

  constructor(
    private httpService : HttpService,
    private organizationsService : OrganizationsService,
    private menuService : MenuService) {
      this.loadUsers(true)
    }

  match(item : User, value : string) : boolean {
    if (item.name && item.name.includes(value)) {
      return true
    }
    if (item.email && item.email.includes(value)) {
      return true
    }
    if (item.member && item.member && item.member.role.includes(value)) {
      return true
    }
    return false
  }

  loadUsers(refresh : boolean) {
    if (!refresh && this.users.length>0) {
      this.onUsersLoaded.next()
      return
    }
    this.httpService.users().subscribe(
      data => {
        this.users = data
        this.onUsersLoaded.next()
      },
      error => {
        this.onUsersError.next(error)
      }
    );
  }

  logout() {
    this.currentUser = this.noLoginUser
    localStorage.removeItem('token');
    this.menuService.navigate(["/auth", "signin"])
  }

  login(token : string) {
    localStorage.setItem('token', JSON.stringify({ token: token }));
    this.setCurrentUser(token, true)
  }

  setCurrentUser(token : string, nav : boolean) {
    let plainToken = this.parseJwt(token)
    if (plainToken.Type != 'login' || plainToken.iss !='amplifier') {
      return
    }
    this.httpService.setToken(token)
    this.currentUser = new User(plainToken.AccountName, "", "")
    localStorage.setItem('currentUser', JSON.stringify({ username: this.currentUser.name, email: this.currentUser.email }));
    this.httpService.userOrganization(this.currentUser.name).subscribe(
      data => {
        this.organizationsService.organizations = data
        this.organizationsService.currentOrganization = this.organizationsService.noOrganization
        for (let org of data) {
          if (org.name == plainToken.ActiveOrganization) {
            this.organizationsService.currentOrganization = org
          }
        }
        this.httpService.users().subscribe(
          data => {
            this.users = data
            if (nav) {
              this.menuService.navigate(['/amp', 'dashboard'])
            }
          },
          error => {
            console.log(error)
            this.logout()
            this.onUsersError.next(error)
          }
        );
      },
      error => {
        console.log(error)
        this.logout()
        this.onUsersError.next(error)
      }
    )
  }

  switchToUserOnly() {
    if (this.currentUser) {
      this.httpService.switchOrganization(this.currentUser.name).subscribe(
        (rep) => {
          let data = rep.json()
          let token = data.auth
          this.httpService.setToken(token)
          localStorage.setItem('token', JSON.stringify({ token: token }));
          this.organizationsService.currentOrganization = this.organizationsService.noOrganization
        }
      )
    }
  }

  parseJwt (token) {
    if (!token) {
      return {}
    }
    let base64Url = token.split('.')[1];
    let base64 = base64Url.replace('-', '+').replace('_', '/');
    let ret = JSON.parse(window.atob(base64));
    console.log("token organization: "+ret.ActiveOrganization)
    return ret
  }

  signup(user : User, pwd : string) {
    this.users.push(user)
    //this.onUserEndCreateMode.emit();
    localStorage.setItem('currentUser', JSON.stringify({ username: user.name, email: user.email }));
    this.menuService.navigate(["/auth", "signin"])
  }

  isAuthenticated() {
    if (this.currentUser === this.noLoginUser) {
      return false
    }
    return true
  }

  returnToCaller() {
    if (this.isAuthenticated()) {
      this.menuService.navigate(["/amp", "users"]);
    } else {
      this.menuService.navigate(["/auth", "signin"]);
    }
  }

  //to be refactor with associative array
  getAllNoMembers(members : Member[]) {
    let list : Member [] = []
    for (let user of this.users) {
      let found = false
      for (let member of members) {
        if (member.userName == user.name) {
          found= true
          break;
        }
      }
      if (!found) {
        list.push(new Member(user.name, undefined))
      }
    }
    return list
  }

  getUserList(orgName : string) : User[] {
    let userList : User[] = []
    for (let org of this.organizationsService.organizations) {
      if (org.name == orgName) {
        for (let member of org.members) {
          for (let user of this.users) {
            if (member.userName == user.name) {
              user.member = member
              userList.push(user)
            }
          }
        }
      }
    }
    return userList
  }


}
