import { Injectable, Output, EventEmitter } from '@angular/core';
import { Http, Headers, Response } from '@angular/http';
import { User } from '../models/user.model';
import { Router } from '@angular/router';

@Injectable()
export class UsersService {
  users : User[] = []
  noLoginUser = new User('not signin','','','')
  currentUser = this.noLoginUser
  //@Output() onUserEndCreateMode = new EventEmitter<void>();
  @Output() onUserLogout = new EventEmitter<void>();

  constructor(private router : Router, private http : Http) {
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
  /*
    this.http.get("http://...?auth=token").subcribe(
      (error) => console.log(error),
      (response : Response) => {
        const data = response.json()
      }
    )
  */
  }

  createModeOn(from : string) {
    this.router.navigate(["auth/signup"])
    //this.router.navigate(["auth/signup", var, 'bouturl'], { queryPrarams: { myParam : '1'}, fragment: "myFragment"});
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
