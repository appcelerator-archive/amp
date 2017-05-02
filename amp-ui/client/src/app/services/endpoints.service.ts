import { Injectable, Output, EventEmitter } from '@angular/core';
import { Router  } from '@angular/router';
import { HttpService } from '../services/http.service';
import { Subject } from 'rxjs/Subject'
import { UsersService } from '../services/users.service';
import { User } from '../models/user.model';
import { Endpoint } from '../models/endpoint.model';

@Injectable()
export class EndpointsService {
  onEndpointsLoaded = new Subject();
  onEndpointsError = new Subject();
  endpoints = []
  currentEndpoint : Endpoint = null

  constructor(private httpService : HttpService, private usersService : UsersService) {
  }

  loadEndpoints() {
    this.httpService.endpoints().subscribe(
      data => {
        this.convertLocalEndpoint(data);
        this.onEndpointsLoaded.next()
      },
      error => this.onEndpointsError.next(error)
    );
  }

  convertLocalEndpoint(list : string[]) {
    this.endpoints = []
    for (let host of list) {
      let endpoint : Endpoint = null
      if (host == "LocalEndpoint") {
        endpoint = new Endpoint(window.location.hostname, true)
        this.endpoints.push(endpoint);
      } else {
        endpoint = new Endpoint(host, false)
        this.endpoints.push(endpoint)
      }
      if (this.currentEndpoint == null) {
        this.currentEndpoint = endpoint
      }
    }
    console.log(this.endpoints)
  }

  connectAndLogin(user : User, pwd : string) {
    this.httpService.connectEndpoint(this.currentEndpoint).subscribe(
      data => {
        this.usersService.login(user, pwd, this.currentEndpoint.host)
      },
      error => { alert(error) },
    );
  }

  setCurrentEndpoint(name : string) {
    for (let endpoint of this.endpoints) {
      if (endpoint.host == name) {
        this.currentEndpoint = endpoint
      }
    }
  }

  selectEndpoint(endpoint : Endpoint) {

  }

}
