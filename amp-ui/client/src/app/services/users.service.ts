import { Injectable, Output, EventEmitter } from '@angular/core';
import { HttpService } from '../services/http.service';
import { User } from '../models/user.model';
import { Router } from '@angular/router';
import { Subject } from 'rxjs/Subject'

@Injectable()
export class UsersService {
  onUsersLoaded = new Subject();
  onUsersError = new Subject();
  users : User[] = []
  noLoginUser = new User('not signin','','')
  currentUser = this.noLoginUser
  @Output() onUserLogout = new EventEmitter<void>();

  constructor(private router : Router, private httpService : HttpService) {}

  match(item : User, value : string) : boolean {
    if (item.name.includes(value)) {
      return true
    }
    if (item.email.includes(value)) {
      return true
    }
    if (item.role.includes(value)) {
      return true
    }
    return false
  }

  loadUsers() {
    this.httpService.users().subscribe(
      data => {
        this.users = data
        console.log(data)
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
    this.router.navigate(["/auth/signin"])
  }

  login(user : User, pwd : string) {
    this.currentUser = user
    this.httpService.login(user, pwd).subscribe(
      data => {
        let ret = data.json()
        this.httpService.setToken(ret.data)
        localStorage.setItem('currentUser', JSON.stringify({ username: user.name, token: ret.data, pwd: pwd }));
        localStorage.setItem('lastUser', JSON.stringify({ username: user.name}));
        this.router.navigate(["/amp"])
      },
      error => {
        localStorage.removeItem('currentUser');
        this.onUsersError.next(error)
      }
    );
  }

  setCurrentUser(currentUser : {username : string, endpointname: string, token : string}) {
    this.currentUser = new User(currentUser.username, "", "")
    this.httpService.setToken(currentUser.token)
    this.router.navigate(["/amp"]);
  }

  signup(user : User, pwd : string) {
    this.users.push(user)
    //this.onUserEndCreateMode.emit();
    this.router.navigate(["/auth/signin"])
  }

  isAuthenticated() {
    if (this.currentUser === this.noLoginUser) {
      return false
    }
    return true
  }

  returnToCaller() {
    if (this.isAuthenticated()) {
      this.router.navigate(["/amp/users"]);
    } else {
      this.router.navigate(["/auth/signin"]);
    }
  }

}
