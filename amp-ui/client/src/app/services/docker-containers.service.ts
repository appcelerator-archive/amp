import { Injectable } from '@angular/core';
import { DockerContainer } from '../models/docker-container.model';

@Injectable()
export class DockerContainersService {
  containers : DockerContainer[] = []
  currentContainer : DockerContainer = new DockerContainer("", "")

  constructor() {
    //temporary for debug
    this.containers.push(new DockerContainer("763f7e23e82c", "testContainer"))
  }

  match(cont : DockerContainer, value : string) : boolean {
    if (value == "") {
      return true
    }
    if (cont.id && cont.id.includes(value)) {
      return true
    }
    if (cont.name && cont.name.includes(value)) {
      return true
    }
    if (cont.image && cont.image.toString().includes(value)) {
      return true
    }
    return false
  }

  setCurrentContainer(id) {
    for (let cont of this.containers) {
      if (cont.id === id) {
        this.currentContainer = cont
      }
    }
  }

}
