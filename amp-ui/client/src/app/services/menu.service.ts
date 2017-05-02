import { Injectable, Output, EventEmitter } from '@angular/core';
import { Router } from '@angular/router';
import { ItemMenu } from '../models/item-menu.model';
import { Subject } from 'rxjs/Subject';

@Injectable()
export class MenuService {
  currentMenuItem : ItemMenu = new ItemMenu("","")
  autoRefresh : boolean = false
  onRefreshClicked = new Subject()

  constructor(private router : Router) { }

  setItemMenu(name : string, description : string) {
    let item = new ItemMenu(name, description)
    this.currentMenuItem = item
  }

  navigate(path) {
    this.router.navigate(path)
  }

  refresh() {
    this.autoRefresh=false
    this.onRefreshClicked.next()
  }

}
