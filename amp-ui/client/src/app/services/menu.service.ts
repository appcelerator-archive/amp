import { Injectable, Output, EventEmitter } from '@angular/core';
import { Router } from '@angular/router';

@Injectable()
export class MenuService {
  currentMenuItem = ""
  itemDescription = {
    "dashboard" : "Home",
    "nodes" : "Node list",
    "stacks" : "Stack list",
    "password" : "Settings",
    "users" : "Users management",
    "endpoints" : "Endpoints management"
  }

  @Output() onMenuItemSelected = new EventEmitter<string>();

  constructor(private router : Router) { }

  getCurrentItemDescription() {
    return this.itemDescription[this.currentMenuItem]
  }

  setItemMenu(item : string) {
    this.currentMenuItem = item
    this.onMenuItemSelected.emit(item)
    this.router.navigate(["/amp", item])
  }
}
