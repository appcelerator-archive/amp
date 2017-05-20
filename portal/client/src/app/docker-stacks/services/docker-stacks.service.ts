import { Injectable } from '@angular/core';
import { DockerStack } from '../models/docker-stack.model';
import { DockerService } from '../models/docker-service.model';
import { DockerContainer } from '../models/docker-container.model';
import { Subject } from 'rxjs/Subject'
import { HttpService } from '../../services/http.service';

@Injectable()
export class DockerStacksService {
  stacks : DockerStack[] = []
  emptyStack : DockerStack = new DockerStack("", "", 0, "", "")
  currentStack : DockerStack = this.emptyStack
  onStacksLoaded = new Subject();
  onStacksError = new Subject();

  constructor(private httpService : HttpService) {
  }

  match(stack : DockerStack, value : string) : boolean {
    if (value == "") {
      return true
    }
    if (stack.id && stack.id.includes(value)) {
      return true
    }
    if (stack.name && stack.name.includes(value)) {
      return true
    }
    if (stack.services && stack.services.toString().includes(value)) {
      return true
    }
    if (stack.ownerName && stack.ownerName.includes(value)) {
      return true
    }
    if (stack.ownerType && stack.ownerType.includes(value)) {
      return true
    }
    return false
  }

  setCurrentStack(name) {
    if (this.currentStack.name == name) {
      this.currentStack = this.emptyStack
      return
    }
    this.currentStack = this.emptyStack
    for (let stack of this.stacks) {
      if (stack.name === name) {
        this.currentStack = stack
      }
    }
  }

  loadStacks(refresh : boolean) {
    if (!refresh && this.stacks.length > 0) {
      this.onStacksLoaded.next()
      return
    }
    this.httpService.stacks().subscribe(
      data => {
        this.stacks = data
        this.onStacksLoaded.next()
      },
      error => {
        //console.log("loadStacksError")
        //console.log(error)
        this.onStacksError.next(error)
      }
    );
  }

}
