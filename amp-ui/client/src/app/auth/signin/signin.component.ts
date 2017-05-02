import { Component, OnInit, OnDestroy } from '@angular/core';
import { User } from '../../models/user.model';
import { UsersService } from '../../services/users.service';
import { NgForm } from '@angular/forms';
import { MenuService } from '../../services/menu.service';

@Component({
  selector: 'app-signin',
  templateUrl: './signin.component.html',
  styleUrls: ['./signin.component.css']
})
export class SigninComponent implements OnInit, OnDestroy {
  message = ""
  messageError = ""
  constructor(
    public usersService : UsersService,
    private menuService : MenuService) { }

  ngOnInit() {
    let currentUser = JSON.parse(localStorage.getItem('currentUser'));
    if (currentUser) {
      this.usersService.loadUsers()
      this.usersService.setCurrentUser(currentUser)
    }
    this.menuService.navigate(['amp', 'dashboard'])
  }

  ngOnDestroy() {
    //this.endpointsService.onEndpointsLoaded.unsubscribe()
    //this.endpointsService.onEndpointsError.unsubscribe()
  }

  signin(form : NgForm) {
    let user = new User(form.value.username, '', '')
    this.usersService.login(user, form.value.password)
  }

}
