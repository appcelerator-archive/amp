import { Injectable, Output, EventEmitter } from '@angular/core';
import { Router } from '@angular/router';
import { ItemMenu } from '../models/item-menu.model'
import { Subject } from 'rxjs/Subject'

@Injectable()
export class MenuService {
  currentMenuItem : ItemMenu = new ItemMenu("","","")
  autoRefresh : boolean = false
  onRefreshClicked = new Subject()

  constructor(private router : Router) { }

  setItemMenu(name : string, description : string, routePath : string) {
    if (routePath === "") {
      routePath = "/amp/"+name
    }
    let item = new ItemMenu(name, description, routePath)
    //this.onMenuItemSelected.emit(Item)
    this.currentMenuItem = item
    this.router.navigate([routePath])
  }

  refresh() {
    this.autoRefresh=false
    this.onRefreshClicked.next()
  }

}
