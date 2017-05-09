import { Component, OnInit, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import { User } from '../../models/user.model';
import { UsersService } from '../../services/users.service'
import { MenuService } from '../../services/menu.service'

@Component({
  selector: 'app-signup',
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.css']
})
export class SignupComponent implements OnInit {
  @ViewChild ('f') form: NgForm;

  constructor(
    public usersService : UsersService,
    private menuService : MenuService) { }

  ngOnInit() {
    this.menuService.setItemMenu('users', 'sign up')
  }

  onSignup() {
    let user = new User(this.form.value.username, this.form.value.email, '')
    this.usersService.signup(user, this.form.value.password)
  }

}
