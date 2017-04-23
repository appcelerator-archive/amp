import { Component, OnInit, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';

import { User } from '../../models/user.model';
import { UsersService } from '../../services/users.service'

@Component({
  selector: 'app-signup',
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.css']
})
export class SignupComponent implements OnInit {
  @ViewChild ('f') form: NgForm;

  constructor(public usersService : UsersService) { }

  ngOnInit() {
  }

  onSignup() {
    console.log(this.form)
    let user = new User(this.form.value.username, this.form.value.email, '')
    this.usersService.signup(user, this.form.value.password)
  }

}
