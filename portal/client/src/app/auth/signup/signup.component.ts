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
  validateLink = false

  constructor(
    public usersService : UsersService,
    private menuService : MenuService,
    private httpService : HttpService) { }

  ngOnInit() {
    this.validateLink = false
    this.menuService.setItemMenu('users', 'sign up')
  }

  onSignup(event : NgForm) {
    if (this.submitCaption == "Done") {
      let previousPath = this.menuService.getPreviousPath()
      if (previousPath.indexOf("signup")>=0) {
        this.menuService.navigate(["/auth", "signin"])
      } else {
        this.menuService.returnToPreviousPath()
      }
      return
    }
    if (event.form.value.password != event.form.value.passwordConfirm) {
        this.messageError = "your password must match"
        return
    }
    this.httpService.signup(event.form.value.username, event.form.value.password, event.form.value.email).subscribe(
      data => {
        this.httpService.registration().subscribe(
          rep => {
            let ret = rep.json()
            console.log(ret)
            this.messageError = ""
            if (ret.email_confirmation) {
              this.message = "Your account is created, you are going to receive an email to validate your account"
            } else {
              this.message = "Your account is created"
            }
            this.submitCaption = "Done"
          },
          err => {
            let error = err.json()
            this.messageError = error.error
          }
        )
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

  returnBack() {
    this.menuService.returnToPreviousPath()
  }

}
