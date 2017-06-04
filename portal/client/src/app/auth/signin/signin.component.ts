import { Component, OnInit, OnDestroy } from '@angular/core';
import { User } from '../../models/user.model';
import { UsersService } from '../../services/users.service';
import { NgForm } from '@angular/forms';
import { MenuService } from '../../services/menu.service';
import { HttpService } from '../../services/http.service';

@Component({
  selector: 'app-signin',
  templateUrl: './signin.component.html',
  styleUrls: ['./signin.component.css']
})
export class SigninComponent implements OnInit, OnDestroy {
  message = ""
  messageError = ""
  byPass = false
  login = ""
  validateLink = false
  constructor(
    public usersService : UsersService,
    private menuService : MenuService,
    private httpService : HttpService) { }

  ngOnInit() {
    this.validateLink = false
    let currentUser = JSON.parse(localStorage.getItem('currentUser'));
    if (currentUser) {
      this.login = currentUser.username
      let token = JSON.parse(localStorage.getItem('token'));
      if (token) {
        this.byPass = true
        this.usersService.setCurrentUser(token.token, true)
      }
    }
  }

  ngOnDestroy() {
    //this.endpointsService.onEndpointsLoaded.unsubscribe()
    //this.endpointsService.onEndpointsError.unsubscribe()
  }

  signin(form : NgForm) {
    this.httpService.login(form.value.username, form.value.password).subscribe(
      data => {
        let ret = data.json()
        this.usersService.login(ret.auth)
      },
      error => {
        let data = error.json()
        if (!data.error) {
          this.validateGtw()
          return
        }
        this.messageError = data.error
      }
    )
  }

  validateGtw() {
    this.validateLink = true
    this.messageError = "First time: Certificat issue: Please, clic on the link below and accept the connection"
  }

}
