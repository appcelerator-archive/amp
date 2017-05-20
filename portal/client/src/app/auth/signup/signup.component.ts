import { Component, OnInit } from '@angular/core';
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
  message = ""
  messageError = ""
  submitCaption = "Submit"

  constructor(
    public usersService : UsersService,
    private menuService : MenuService,
    private httpService : HttpService) { }

  ngOnInit() {
    this.menuService.setItemMenu('users', 'sign up')
  }

  onSignup(event : NgForm) {
    if (this.submitCaption == "Done") {
      this.menuService.returnToPreviousPath()
      return
    }
    if (event.form.value.password != event.form.value.passwordConfirm) {
        this.messageError = "your password must match"
        return
    }
    this.httpService.signup(event.form.value.username, event.form.value.password, event.form.value.email).subscribe(
      data => {
        this.message = "Your account is created, you are going to receive an email to validate your account"
        this.submitCaption = "Done"
      },
      error => {
        console.log(error)
        let data = error.json()
        if (!data.error) {
          this.messageError = "Certificat issue: You need to import amp certificate in your browser. See documentation"
        }
        this.message = data.error
      }
    )
  }

  returnBack() {
    this.menuService.returnToPreviousPath()
  }

}
