import { Injectable, Output, EventEmitter } from '@angular/core';
import { HttpService } from '../services/http.service';
import { OrganizationsService } from '../organizations/services/organizations.service';
import { User } from '../models/user.model';
import { Member } from '../models/member.model';
import { Organization } from '../models/organization.model';
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
    private menuService : MenuService) {}

  match(item : User, value : string) : boolean {
    if (item.name && item.name.includes(value)) {
      return true
    }
    if (item.email && item.email.includes(value)) {
      return true
    }
    if (item.role && item.role.includes(value)) {
      return true
    }
    return false
  }

  loadUsers() {
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
    localStorage.removeItem('currentUser');
    this.menuService.navigate(["/auth", "signin"])
  }

  login(username : string, pwd : string) {
    this.httpService.login(username, pwd).subscribe(
      data => {
        let ret = data.json()
        localStorage.setItem('currentUser', JSON.stringify({ username: username, token: ret.auth }));
        localStorage.setItem('lastUser', JSON.stringify({ username: username}));
        this.setCurrentUser(username, ret.auth, true)
      },
      error => {
        localStorage.removeItem('currentUser');
        this.onUsersError.next(error)
      }
    );
  }

  setCurrentUser(username : string, token : string, nav : boolean) {
    this.httpService.setToken(token)
    this.currentUser = new User(username, "", "")
    this.httpService.userOrganization(this.currentUser).subscribe(
      data => {
        this.organizationsService.organizations = data
        this.httpService.users().subscribe(
          data => {
            this.users = data
            if (nav) {
              this.menuService.navigate(['/amp', 'dashboard'])
            }
          },
          error => {
            console.log(error)
            this.onUsersError.next(error)
          }
        );
      },
      error => {
        console.log(error)
        this.onUsersError.next(error)
      }
    )
  }

  signup(user : User, pwd : string) {
    this.users.push(user)
    //this.onUserEndCreateMode.emit();
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
        list.push(new Member(user.name, 0))
      }
    }
    return list
  }

}
