import { Injectable, Output, EventEmitter } from '@angular/core';
import { Router  } from '@angular/router';
import { HttpService } from '../services/http.service';
import { Subject } from 'rxjs/Subject'

@Injectable()
export class EndpointsService {
  onEndpointsLoaded = new Subject();
  onEndpointsError = new Subject();
  endpoints = []//this.convertLocalEndpoint([ "LocalEndpoint", "cluster.atomiq.io"])
  currentEndpoint = this.endpoints[0]

  constructor(private httpService : HttpService) {
  }

  loadEndpoints() {
    this.httpService.getEndpoints().subscribe(
      data => {
        this.endpoints = this.convertLocalEndpoint(data);
        if ((data.length) > 1) {
          this.currentEndpoint = data[0]
        }
        this.onEndpointsLoaded.next()
      },
      error => this.onEndpointsError.next(error)
    );
  }

  convertLocalEndpoint(list : string[]) : string[] {
    for (let i=0; i<list.length; i++) {
      if (list[i] == "LocalEndpoint") {
        list[i] = window.location.hostname
      }
    }
    return list
  }

}
