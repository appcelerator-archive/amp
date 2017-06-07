import { Injectable, Output, EventEmitter } from '@angular/core';
import { Router } from '@angular/router';
import { ItemMenu } from '../models/item-menu.model';
import { AppWindow } from '../models/app-window.model';
import { Subject } from 'rxjs/Subject';


@Injectable()
export class MenuService {
  currentMenuItem : ItemMenu = new ItemMenu("","")
  tooltipLabel = ""
  autoRefresh : boolean = false
  onRefreshClicked = new Subject()
  public onWindowResize = new Subject<AppWindow>();
  private paths : string[] = []
  private lastPath = ""
  private cursorClass = ""
  paddingLeftMenu = 250
  currentTimer : any
  public appWindow : AppWindow = new AppWindow(document.documentElement.clientWidth, document.documentElement.clientHeight)

  constructor(private router : Router) {
  }

  resize(evt : UIEvent) {
    this.appWindow = new AppWindow((<Window>event.target).innerWidth, (<Window>event.target).innerHeight)
    this.onWindowResize.next(this.appWindow);
  }

  setItemMenu(name : string, description : string) {
    this.tooltipLabel = ""
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

  getPreviousPath() : string {
    return this.paths.pop()
  }

  returnToPreviousPath() {
    let path = this.paths.pop()
    path = this.paths.pop()
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

  public setCurrentTimer(timer : any) {
    if (this.currentTimer) {
      clearInterval(this.currentTimer)
    }
    this.currentTimer = timer
  }

  public clearCurrentTimer() {
    if (this.currentTimer) {
      clearInterval(this.currentTimer)
    }
    this.currentTimer = undefined
  }

}
