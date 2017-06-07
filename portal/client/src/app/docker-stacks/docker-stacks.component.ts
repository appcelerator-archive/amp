import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { DockerStack } from './models/docker-stack.model';
import { DockerStacksService } from './services/docker-stacks.service';
import { MenuService } from '../services/menu.service';
import { MetricsService } from '../metrics/services/metrics.service';
import { ListService } from '../services/list.service';
import { HttpService } from '../services/http.service';
import { Observable } from 'rxjs/Observable';

@Component({
  selector: 'app-stacks',
  templateUrl: './docker-stacks.component.html',
  styleUrls: ['./docker-stacks.component.css'],
  providers: [ ListService ]
})
export class DockerStacksComponent implements OnInit {
  deployTitle = "Deploy"
  message = ""

  constructor(
    private route : ActivatedRoute,
    public dockerStacksService : DockerStacksService,
    public listService : ListService,
    private menuService : MenuService,
    private metricsService : MetricsService,
    private httpService : HttpService) {
      listService.setFilterFunction(dockerStacksService.match)
    }

  ngOnInit() {
    this.menuService.setItemMenu('stacks', 'List')
    this.dockerStacksService.onStacksLoaded.subscribe(
      () => {
        this.listService.setData(this.dockerStacksService.stacks)
        let id = this.dockerStacksService.currentStack.id
        if (id == "") {
            this.deployTitle="Deploy"
        } else {
          this.deployTitle="Update"
        }
      }
    )
    this.menuService.onRefreshClicked.subscribe(
      () => {
        this.dockerStacksService.loadStacks(true)
      }
    )
    this.dockerStacksService.loadStacks(false)
    this.menuService.setCurrentTimer( setInterval(
      () => {
        this.dockerStacksService.loadStacks(true)
      },
      3000
    ))
  }

  ngOnDestroy() {
    this.menuService.clearCurrentTimer()
    //this.dockerStacksService.onStacksLoaded.unsubscribe();
  }

  serviceList(stackName : string) {
    this.dockerStacksService.setCurrentStack(stackName)
    this.menuService.navigate(["/amp", "stacks", stackName, "services"])
  }

  getColor(stack : DockerStack) : string {
    if (stack.status == 'running') {
      return 'limegreen'
    }
    if (stack.status == 'starting') {
      return 'orange'
    }
    return 'red';
  }

  deploy() {
    this.menuService.navigate(["/amp", "stacks", "deploy"])
  }

  update() {
    let stackName = this.dockerStacksService.currentStack.name
    this.menuService.navigate(["/amp", "stacks", stackName, "update"])
  }

  remove() {
    this.httpService.removeStack(this.dockerStacksService.currentStack.name).subscribe(
      data => {
        this.dockerStacksService.setCurrentStack("")
        this.dockerStacksService.loadStacks(true)
      },
      error => {
        let data = error.json()
        this.message = data.error
      }
    )
  }

  metrics(stackName : string) {
    this.menuService.navigate(['/amp', 'metrics', 'stack', 'single', stackName])
  }

  logs(stackName : string) {
    this.menuService.navigate(['/amp', 'logs', 'stack', stackName])
  }

}
