import { CanActivate, CanActivateChild, ActivatedRouteSnapshot, RouterStateSnapshot, Router } from '@angular/router';
import { Injectable } from '@angular/core';

import { UsersService } from './users.service'

@Injectable()
export class AuthGuard implements CanActivate, CanActivateChild {
  constructor(private usersService : UsersService, private router : Router ) {}

  canActivate(route : ActivatedRouteSnapshot, state : RouterStateSnapshot) {
    if (localStorage.getItem('currentUser')) {
      let currentUser = JSON.parse(localStorage.getItem('currentUser'));
      if (this.usersService.currentUser.name !== currentUser.username) {
        console.log(this.usersService.currentUser.name + "<>" + currentUser.username + "-> reload user")
        this.usersService.setCurrentUser(currentUser.username, currentUser.token, false)
      }
      return true;
    }
    this.router.navigate(['/auth", "signin'])
    return false
  }

  canActivateChild(route : ActivatedRouteSnapshot, state : RouterStateSnapshot) {
    return this.canActivate(route, state);
  }
}
