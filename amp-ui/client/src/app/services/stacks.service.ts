import { Injectable } from '@angular/core';
import { Stack } from '../models/stack.model';

@Injectable()
export class StacksService {
  stacks : Stack[] = []
  stack : Stack = null
  constructor() {
    this.stacks.push(new Stack("907345782", "haproxy", 1, "freignat", "USER"))
    this.stacks.push(new Stack("234523455", "pinger", 3, "freignat", "USER"))
    this.stacks.push(new Stack("832323642", "funchttp", 2, "bquenin", "USER"))
  }

  match(stack : Stack, value : string) : boolean {
    if (stack.id.includes(value)) {
      return true
    }
    if (stack.name.includes(value)) {
      return true
    }
    if (stack.services.toString().includes(value)) {
      return true
    }
    if (stack.ownerName.includes(value)) {
      return true
    }
    if (stack.ownerType.includes(value)) {
      return true
    }
    return false
  }
}
