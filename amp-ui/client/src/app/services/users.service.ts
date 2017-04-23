import { Injectable, Output, EventEmitter } from '@angular/core';
import { HttpService } from '../services/http.service';
import { User } from '../models/user.model';
import { Router } from '@angular/router';
import { Subject } from 'rxjs/Subject'

@Injectable()
export class UsersService {
  users : User[] = []
  noLoginUser = new User('not signin','','','')
  currentUser = this.noLoginUser
  onUsersLoaded = new Subject();
  @Output() onUserLogout = new EventEmitter<void>();

  constructor(private router : Router, private httpService : HttpService) {
    this.users.push(new User('freignat', 'freignat@axway.com', '', 'USER'))
    this.users.push(new User('bquenin', 'bquenin@axway.com', '', 'USER'))
  }

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
    this.httpService.getUsers().subscribe(
      data => {
        this.users = data
        this.onUsersLoaded.next()
      },
      //error => {}
    );
  }

  logout() {
    this.currentUser = this.noLoginUser
    this.router.navigate(["/auth/signin"])
  }

  login(user : User) {
    this.currentUser = user
    //const headers = new Headers({'Content-Type' ,'application/json'})
    //obs = this.http.post("http://...", user, {headers: headers})
    this.router.navigate(["/amp"])
  }

  signup(user : User) {
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


/*
this.usersService.storeServers(this.servers) {
  .subcribe(
    (response) => console.log(response)
    (error) => console.log(error)
  )
}
*/
