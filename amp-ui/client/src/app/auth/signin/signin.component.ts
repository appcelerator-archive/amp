import { Component, OnInit, OnDestroy } from '@angular/core';
import { User } from '../../models/user.model';
import { UsersService } from '../../services/users.service'
import { EndpointsService } from '../../services/endpoints.service'
import { NgForm } from '@angular/forms';

@Component({
  selector: 'app-signin',
  templateUrl: './signin.component.html',
  styleUrls: ['./signin.component.css']
})
export class SigninComponent implements OnInit, OnDestroy {
  message = ""
  messageError = ""
  constructor(public usersService : UsersService, public endpointsService : EndpointsService) { }

  ngOnInit() {
    if (this.endpointsService.endpoints.length == 0) {
      this.message="loading endpoints list..."
      this.endpointsService.loadEndpoints()
      this.endpointsService.onEndpointsLoaded.subscribe(
        () => {this.message="";}
      )
      this.endpointsService.onEndpointsError.subscribe(
        (err : any) => {
          console.log(err)
          this.messageError=err.statusText+": "+err._body;
        }
      )
    }
  }

  ngOnDestroy() {
    this.endpointsService.onEndpointsLoaded.unsubscribe()
    this.endpointsService.onEndpointsError.unsubscribe()
  }

  signin(form : NgForm) {
    let user = new User(form.value.username, '', form.value.password, '')
    this.usersService.login(user)
  }
}
