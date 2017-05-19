import { Injectable } from '@angular/core';
import { DockerService } from '../models/docker-service.model';

@Injectable()
export class DockerServicesService {
  services : DockerService[] = []
  currentService : DockerService = new DockerService("", "")

  constructor() {
    //temporary for debug
    this.services.push(new DockerService("3ab23e42cd8", "testService1"))
    this.services.push(new DockerService("4534a34bc29", "testService2"))
  }

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

}
