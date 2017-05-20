import { Injectable } from '@angular/core';
import { Node } from '../models/node.model';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { Subject } from 'rxjs/Subject'

@Injectable()
export class NodesService {
  nodes : Node[] = []
  emptyNode : Node = new Node("")
  currentNode : Node = this.emptyNode
  onNodesLoaded = new Subject();
  onNodesError = new Subject<string>();

  constructor(
    private httpService : HttpService,
    private menuService : MenuService) {
  }


  match(item : Node, value : string) : boolean {
    if (item.id && item.id.includes(value)) {
      return true
    }
    if (item.shortId && item.shortId.includes(value)) {
      return true
    }
    if (item.name && item.name.includes(value)) {
      return true
    }
    if (item.role && item.role.includes(value)) {
      return true
    }
    if (item.hostname && item.hostname.includes(value)) {
      return true
    }
    if (item.architecture && item.architecture.includes(value)) {
      return true
    }
    if (item.os && item.os.includes(value)) {
      return true
    }
    if (item.engine && item.engine.includes(value)) {
      return true
    }
    if (item.status && item.status.includes(value)) {
      return true
    }
    if (item.availability && item.availability.includes(value)) {
      return true
    }
    if (item.addr && item.addr.includes(value)) {
      return true
    }
    if (item.reachability && item.reachability.includes(value)) {
      return true
    }

    return false
  }

  loadNodes(refresh : boolean) {
    if (!refresh && this.nodes.length>0) {
      this.onNodesLoaded.next()
      return  
    }
    this.httpService.nodes().subscribe(
      data => {
        this.nodes = data;
        this.onNodesLoaded.next()
      },
      error => {
        let data = error.json()
        this.onNodesError.next(data.error)
      }
    );
  }

  setCurrentNode(id : string) {
    if (this.currentNode.id == id) {
      this.currentNode = this.emptyNode
      return
    }
    this.currentNode = this.emptyNode
    for (let node of this.nodes) {
      if (node.id === id) {
        this.currentNode = node
      }
    }
  }

}
