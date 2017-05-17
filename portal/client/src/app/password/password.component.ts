import { Component, OnInit } from '@angular/core';
import { MenuService } from '../services/menu.service'
import { UsersService } from '../services/users.service'
import { NgForm } from '@angular/forms';
import { HttpService } from '../services/http.service'

@Component({
  selector: 'app-password',
  templateUrl: './password.component.html',
  styleUrls: ['./password.component.css']
})
export class PasswordComponent implements OnInit {
  message = ""
  submitCaption = "Submit"

  constructor(
    public usersService : UsersService,
    private menuService : MenuService,
    private httpService : HttpService) {
  }

  ngOnInit() {
    this.menuService.setItemMenu('Password', 'change')
  }

  returnBack() {
    this.menuService.returnToPreviousPath()
  }

  onValidation(event : NgForm) {
    this.message = ""
    if (this.submitCaption == "Done") {
      this.menuService.returnToPreviousPath()
      return
    }
    console.log(event.form.value.newPassword +","+ event.form.value.newPasswordVerif)
    if (event.form.value.newPassword != event.form.value.newPasswordVerif) {
      this.message = "The new password is not the same than the new password verification"
      return
    }
    this.httpService.changePassword(event.form.value.currentPassword, event.form.value.newPassword).subscribe(
      data => {
        this.message="Your password has been changed"
        this.submitCaption = "Done"
      },
      error => {
        let data = error.json()
        this.message = data.error
      }
    )
  }

}
