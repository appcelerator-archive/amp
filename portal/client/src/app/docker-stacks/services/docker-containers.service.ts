import { Injectable } from '@angular/core';
import { DockerContainer } from '../models/docker-container.model';
import { DockerServicesService } from './docker-services.service';
import { HttpService } from '../../services/http.service';
import { Subject } from 'rxjs/Subject'

@Injectable()
export class DockerContainersService {
  containers : DockerContainer[] = []
  emptyContainer : DockerContainer = new DockerContainer("", "", "", "", "")
  currentContainer : DockerContainer = this.emptyContainer;
  onContainersLoaded = new Subject();
  onContainersError = new Subject();
  currentLoadedServiceId = ""

  constructor(
    private dockerServicesService : DockerServicesService,
    private httpService : HttpService) {}


  match(cont : DockerContainer, value : string) : boolean {
    if (value == "") {
      return true
    }
    if (cont.id && cont.id.includes(value)) {
      return true
    }
    if (cont.state && cont.state.includes(value)) {
      return true
    }
    if (cont.desiredState && cont.desiredState.includes(value)) {
      return true
    }
    if (cont.nodeId && cont.nodeId.includes(value)) {
      return true
    }
    if (cont.image && cont.image.toString().includes(value)) {
      return true
    }
    return false
  }

  setCurrentContainer(id) {
    if (this.currentContainer.id == id) {
      this.currentContainer = this.emptyContainer
      return
    }
    this.currentContainer = this.emptyContainer
    for (let cont of this.containers) {
      if (cont.id === id) {
        this.currentContainer = cont
      }
    }
  }

  loadContainers(refresh : boolean) {
    if (!refresh && this.currentLoadedServiceId == this.dockerServicesService.currentService.id) {
      this.onContainersLoaded.next()
      return
    }
    this.httpService.tasks(this.dockerServicesService.currentService.id).subscribe(
      data => {
        this.containers = data
        this.currentLoadedServiceId = this.dockerServicesService.currentService.id
        this.onContainersLoaded.next()
      },
      error => {
        console.log(error)
        this.onContainersError.next(error)
      }
    );
  }

}
