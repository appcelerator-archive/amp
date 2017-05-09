import { Injectable } from '@angular/core';
import { Swarm } from '../models/swarm.model';

@Injectable()
export class SwarmsService {
  swarms : Swarm[] = []
  allSwarm : Swarm = new Swarm("<all>", "")
  currentSwarm : Swarm = this.allSwarm
  constructor() { }

}
