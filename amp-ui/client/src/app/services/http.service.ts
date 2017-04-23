import { Injectable } from '@angular/core';
import { Http, Headers, Response } from '@angular/http';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import { User } from '../models/user.model'
import { Subject } from 'rxjs/Subject'
import {Observable} from 'rxjs/Observable';

@Injectable()
export class HttpService {
  onHttpError = new Subject();
  baseURL = window.location.protocol + window.location.hostname +":" + window.location.port
  constructor(private http : Http) {}

  getEndpoints() {
    return this.http.get("/api/v1/endpoints")
      .map((res:Response) => res.json())
      //.catch((error : any) => Observable.throw(error.json().error || error))
  }

  getUsers() {
    return this.http.get(this.baseURL+"/api/v1/users")
      .map((res:Response) => {
        const users : User[] = res.json()
        //..
        return users
      })
  }
}
