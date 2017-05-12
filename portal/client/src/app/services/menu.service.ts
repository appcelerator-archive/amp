import { Injectable, Output, EventEmitter } from '@angular/core';
import { Router } from '@angular/router';
import { ItemMenu } from '../models/item-menu.model';
import { AppWindow } from '../models/app-window.model';
import { Subject } from 'rxjs/Subject';


@Injectable()
export class MenuService {
  currentMenuItem : ItemMenu = new ItemMenu("","")
  autoRefresh : boolean = false
  onRefreshClicked = new Subject()
  public onWindowResize = new Subject<AppWindow>();
  private paths : string[] = []
  private lastPath = ""
  private cursorClass = ""
  public appWindow : AppWindow = new AppWindow(document.documentElement.clientWidth, document.documentElement.clientHeight)

  constructor(private router : Router) {
  }

  resize(evt : UIEvent) {
    this.appWindow = new AppWindow((<Window>event.target).innerWidth, (<Window>event.target).innerHeight)
    this.onWindowResize.next(this.appWindow);
  }

  setItemMenu(name : string, description : string) {
    let item = new ItemMenu(name, description)
    this.currentMenuItem = item
  }

  pushPath(path : string) {
    if (!path || this.lastPath==path) {
      return
    }
    this.lastPath = path
    //console.log("push path: "+path)
    this.paths.push(path)
  }

  returnToPreviousPath() {
    let path = this.paths.pop()
    //console.log("return path1: "+path)
    path = this.paths.pop()
    //console.log("return path2: "+path)
    if (path) {
      this.navigate(path.split('/'))
    }
  }

  navigate(path : string[]) {
    this.router.navigate(path)
  }

  refresh() {
    this.autoRefresh=false
    this.onRefreshClicked.next()
  }

  public waitingCursor(mode : boolean) {
      if (mode) {
        this.cursorClass='waiting';
      } else {
        this.cursorClass='';
      }
  }

  public getCursorClass() : string {
    return this.cursorClass
  }
}
