import { Component, OnInit, OnDestroy, ElementRef, AfterViewChecked, ViewChild } from '@angular/core';
import { MenuService } from '../services/menu.service';
import { HttpService } from '../services/http.service';
import { LogsService } from './services/logs.service';
import { ActivatedRoute } from '@angular/router';
import { Log } from './models/log.model';
import { LogsRequest } from './models/logs-request.model';
import { AppWindow } from '../models/app-window.model';
import * as $ from 'jquery';

@Component({
  selector: 'app-logs',
  templateUrl: './logs.component.html',
  styleUrls: ['./logs.component.css']
})
export class LogsComponent implements OnInit, AfterViewChecked, OnDestroy {
  @ViewChild('logScroll') private scrollContainer: ElementRef;

  object = ""
  ref = ""
  routeSub : any
  periodRefreshLabel = "30 sec"
  logs : Log[] = []
  logsPanelHeight = 400
  req : LogsRequest = new LogsRequest()
  filter = ""
  isReturnBack = false
  timer : any
  autoRefreshToggle = false
  lastTimestamp = ""

  constructor(
    private menuService : MenuService,
    private route: ActivatedRoute,
    private httpService : HttpService,
    public logsService : LogsService) {
  }

  ngOnInit() {
    this.menuService.setItemMenu('logs', 'View')
    this.resizeLogs(this.menuService.appWindow)
    this.menuService.onWindowResize.subscribe(
      (win) => {
        this.resizeLogs(win)
      }
    )
    this.routeSub = this.route.params.subscribe(params => {
      let object = params['object']//'global' or 'stack' or 'service' or 'container'
      let ref = params['ref'];//stackName or serviceId or containerId or 'all'
      console.log("object="+object+" ref="+ref)
      this.computeNames(object, ref)
      if (ref) {
        this.isReturnBack = true
      }

      this.req.infra = false
      this.req.size = 1000

      if (object=='stack') {
        this.req.stack = ref
        this.menuService.setItemMenu('logs', 'Vie stack')
      }
      if (object=='service') {
        this.req.service = ref
        this.menuService.setItemMenu('logs', 'View service')
      }
      if (object=='task') {
        this.req.task = ref
        this.menuService.setItemMenu('logs', 'View container')
      }
      if (object=='node') {
        this.req.node = ref
        this.menuService.setItemMenu('logs', 'View node')
      }
      console.log(this.req)
      this.executeRequest()

    })
  }

  ngAfterViewChecked() {
    if (this.logs.length>0) {
      let last = this.logs[this.logs.length-1].timestamp
      if (this.lastTimestamp !== last) {
        this.lastTimestamp = last
        this.scrollToBottom()
      }
    }
    this.alignColumnTitleToContent()
  }

  ngOnDestroy() {
    if (this.timer) {
      clearInterval(this.timer);
    }
  }

  scrollToBottom() {
      try {
        this.scrollContainer.nativeElement.scrollTop = this.scrollContainer.nativeElement.scrollHeight;
      } catch(err) { }
  }

  executeRequest() {
    this.httpService.logs(this.req).subscribe(
      data => {
        this.logs = data
      },
      error => {
        console.log(error)
      }
    )
  }

  autoRefresh() {
    this.autoRefreshToggle = !this.autoRefreshToggle;
    if (this.autoRefresh) {
      this.menuService.setCurrentTimer(setInterval(() => this.executeRequest(), 3000))
    } else {
      if (this.timer) {
        this.menuService.clearCurrentTimer()
      }
    }
  }

  returnBack() {
    this.menuService.returnToPreviousPath()
  }

  resizeLogs(win : AppWindow) {
    this.logsPanelHeight = win.height-240
  }

  alignColumnTitleToContent() {
    let list = this.logsService.getVisibleColsName()
    for (let name of list) {
      $('.'+name+'Title').width($('.'+name+'Cell').width());
    }
  }

  setFilter(filter : string) {
    this.req.message = filter
    console.log(filter)
    this.executeRequest()
  }

  computeNames(object : string, ref : string) {
    if (object == "") {
      this.object = "global"
      return
    }
    this.object = object
    this.ref = ref
  }

}
