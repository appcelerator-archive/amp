import { Component, OnInit, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import { User } from '../../models/user.model';
import { UsersService } from '../../services/users.service'
import { MenuService } from '../../services/menu.service'
import { HttpService } from '../../services/http.service'

@Component({
  selector: 'app-signup',
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.css']
})
export class SignupComponent implements OnInit {
  @ViewChild ('f') form: NgForm;
  message = ""
  submitCaption = "Submit"

  constructor(
    public usersService : UsersService,
    private menuService : MenuService,
    private httpService : HttpService) { }

  ngOnInit() {
    this.menuService.setItemMenu('users', 'sign up')
  }

  onSignup() {
    if (this.submitCaption == "Done") {
      this.menuService.returnToPreviousPath()
      return
    }
    if (this.form.value.password != this.form.value.passwordVerif) {
        this.message = "your password is not the same than your password confirmation"
        return
    }
    this.httpService.signup(this.form.value.username, this.form.value.password, this.form.value.email).subscribe(
      data => {
        this.message = "Your account is created, you are going to receive an email. Confirmed you received it, before sign in"
        this.submitCaption = "Done"
      },
      error => {
        let data = error.json()
        this.message = data.error
      }
    )
  }

  returnBack() {
    this.menuService.returnToPreviousPath()
  }

}
