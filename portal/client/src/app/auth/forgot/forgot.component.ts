import { Component, OnInit } from '@angular/core';
import { NgForm } from '@angular/forms';
import { User } from '../../models/user.model';
import { UsersService } from '../../services/users.service'
import { MenuService } from '../../services/menu.service'
import { HttpService } from '../../services/http.service'

@Component({
  selector: 'app-forgot',
  templateUrl: './forgot.component.html',
  styleUrls: ['./forgot.component.css']
})
export class ForgotComponent implements OnInit {
  message = ""
  submitCaption = "Submit"
  forgotToggle = true

  constructor(
    public usersService : UsersService,
    private menuService : MenuService,
    private httpService : HttpService) { }

  ngOnInit() {
    this.menuService.setItemMenu('users', 'forgot login/pwd')
  }

  onForgot(event : NgForm) {
    if (this.submitCaption == "Done") {
      this.menuService.returnToPreviousPath()
      return
    }
    if (this.forgotToggle) {
      this.onForgotLogin(event)
    } else {
      this.onForgotPassword(event)
    }
  }

  onForgotLogin(event : NgForm) {
      this.httpService.retrieveLogin(event.form.value.email).subscribe(
      data => {
        this.message = "An email has been sent to you with your login name"
        this.submitCaption = "Done"
      },
      error => {
        let data = error.json()
        this.message = data.error
      }
    )
  }

  onForgotPassword(event : NgForm) {
      this.httpService.resetPassword(event.form.value.username).subscribe(
      data => {
        this.message = "An email has been sent to you to allow you to set a new password"
        this.submitCaption = "Done"
      },
      error => {
        let data = error.json()
        this.message = data.error
      }
    )
  }

  forgotLogin() {
    this.message = ""
    this.forgotToggle = true
  }

  forgotPassword() {
    this.message = ""
    this.forgotToggle = false
  }

  returnBack() {
    this.menuService.returnToPreviousPath()
  }

}
