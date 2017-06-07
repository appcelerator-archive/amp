import { Component, OnInit, OnDestroy } from '@angular/core';
import { ListService } from '../../services/list.service';
import { DockerServicesService } from '../services/docker-services.service';
import { DockerStacksService } from '../services/docker-stacks.service';
import { DockerService } from '../models/docker-service.model';
import { DockerStack } from '../models/docker-stack.model';
import { MenuService } from '../../services/menu.service';
import { ActivatedRoute } from '@angular/router';
import { MetricsService } from '../../metrics/services/metrics.service';
import { HttpService } from '../../services/http.service';

@Component({
  selector: 'app-services',
  templateUrl: './docker-services.component.html',
  styleUrls: ['./docker-services.component.css'],
  providers: [ ListService ]
})
export class DockerServicesComponent implements OnInit, OnDestroy {
  routeSub : any
  scaleMode = false
  message = ""

  constructor(
    public listService : ListService,
    public dockerServicesService : DockerServicesService,
    public dockerStacksService : DockerStacksService,
    public menuService : MenuService,
    private route: ActivatedRoute,
    private metricsService : MetricsService,
    private httpService : HttpService) {
    listService.setFilterFunction(dockerServicesService.match)
  }

  ngOnInit() {
    this.scaleMode = false
    this.message = ""
    this.menuService.setItemMenu('services', 'List')
    this.routeSub = this.route.params.subscribe(params => {
      let stackName = params['stackName'];
      for (let stack of this.dockerStacksService.stacks) {
        if (stack.name == stackName) {
          this.dockerStacksService.currentStack = stack
        }
      }
      this.dockerServicesService.onServicesLoaded.subscribe(
        () => {
          this.listService.setData(this.dockerServicesService.services)
        }
      )
      this.menuService.onRefreshClicked.subscribe(
        () => {
          this.dockerServicesService.loadServices(true)
        }
      )
      this.dockerServicesService.loadServices(false)
    })
    this.menuService.setCurrentTimer(setInterval(
      () => {
        this.dockerServicesService.loadServices(true)
      },
      3000
    ))
  }

  ngOnDestroy() {
    this.routeSub.unsubscribe();
    this.menuService.clearCurrentTimer()
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

  containerList(serviceId : string) {
    this.dockerServicesService.setCurrentService(serviceId)
    this.menuService.navigate(["/amp", "stacks", this.dockerStacksService.currentStack.name, "services", serviceId, "containers"])
  }

  selectService(name : string) {
    this.scaleMode = false
    if (this.dockerServicesService.currentService.name== name) {
      this.dockerServicesService.setCurrentService("")
      return
    }
    this.dockerServicesService.setCurrentService(name)
  }

  returnBack() {
    this.menuService.returnToPreviousPath()
  }

  metrics(serviceName : string) {
    this.menuService.navigate(['/amp', 'metrics', 'service', 'single', serviceName])
  }

  logs(serviceName : string) {
    this.menuService.navigate(['/amp', 'logs', 'service', serviceName])
  }

  scaleModeToggle() {
    this.scaleMode = !this.scaleMode
    if (this.dockerServicesService.currentService.name=='') {
      this.scaleMode = false
    }
  }

  scale(number : string) {
    let num = +number
    this.httpService.serviceScale(this.dockerServicesService.currentService.id, num).subscribe(
      () => {
        this.scaleMode = false
        this.dockerServicesService.loadServices(true)
      },
      (err) => {
        this.scaleMode = false
        let error = err.json()
        this.message = error.errror
      }
    )
  }

}
