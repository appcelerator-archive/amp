import { Injectable } from '@angular/core';
import { HttpService } from '../../services/http.service';
import { MenuService } from '../../services/menu.service';
import { Log } from '../models/log.model';

@Injectable()
export class LogsService {
  private columnVisible = {}

  constructor(
    private httpService : HttpService,
    private menuService : MenuService) {
      this.initColumns()
    }

  initColumns() {
    this.columnVisible['timestamp']=true
    this.columnVisible['stackName']=false
    this.columnVisible['serviceName']=false
    this.columnVisible['containerShortName']=false
    this.columnVisible['containerState']=false
    this.columnVisible['nodeId']=false
    this.columnVisible['serviceId']=false
    this.columnVisible['containerId']=false
    this.columnVisible['containerName']=false
    this.columnVisible['taskId']=false
  }

  public isColVisible(name : string) : boolean {
    return this.columnVisible[name]
  }

  public toggleColVisibility(name : string) {
    this.columnVisible[name] = !this.columnVisible[name]
  }

  public getVisibleColsName() {
    let list = []
    for (let key in this.columnVisible) {
      if (this.columnVisible[key]) {
        list.push(key)
      }
    }
    return list
  }
}
