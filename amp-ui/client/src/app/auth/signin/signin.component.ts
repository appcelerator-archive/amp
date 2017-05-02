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
    let currentUser = JSON.parse(localStorage.getItem('currentUser'));
    if (currentUser) {
      this.usersService.setCurrentUser(currentUser)
    }
    if (this.endpointsService.currentEndpoint == null ) {
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
      this.usersService.onUsersError.subscribe(
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
    let user = new User(form.value.username, '', '')
    this.endpointsService.currentEndpoint = form.value.endpoint
    this.endpointsService.connectAndLogin(user, form.value.password)
  }

}
