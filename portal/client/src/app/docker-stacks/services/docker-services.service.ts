import { Injectable } from '@angular/core';
import { DockerService } from '../models/docker-service.model';
import { DockerStacksService } from './docker-stacks.service';
import { HttpService } from '../../services/http.service';
import { Subject } from 'rxjs/Subject'

@Injectable()
export class DockerServicesService {
  services : DockerService[] = []
  currentService : DockerService = new DockerService("", "", "", "", "")
  onServicesLoaded = new Subject();
  onServicesError = new Subject();

  constructor(
    private dockerStacksService : DockerStacksService,
    private httpService : HttpService
  ) {}

  match(serv : DockerService, value : string) : boolean {
    if (value == "") {
      return true
    }
    if (serv.id && serv.id.includes(value)) {
      return true
    }
    if (serv.name && serv.name.includes(value)) {
      return true
    }
    if (serv.image && serv.image.toString().includes(value)) {
      return true
    }
    if (serv.mode && serv.mode.includes(value)) {
      return true
    }
    return false
  }

  setCurrentService(id) {
    for (let service of this.services) {
      if (service.id === id) {
        this.currentService = service
      }
    }
  }

  loadServices() {
    this.httpService.services(this.dockerStacksService.currentStack.name).subscribe(
      data => {
        this.services = data
        this.onServicesLoaded.next()
      },
      error => {
        //console.log("loadStacksError")
        //console.log(error)
        this.onServicesError.next(error)
      }
    );
  }

}
