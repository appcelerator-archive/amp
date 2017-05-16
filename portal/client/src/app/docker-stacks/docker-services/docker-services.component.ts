import { Component, OnInit, OnDestroy } from '@angular/core';
import { ListService } from '../../services/list.service';
import { DockerServicesService } from '../services/docker-services.service';
import { DockerStacksService } from '../services/docker-stacks.service';
import { DockerService } from '../models/docker-service.model';
import { DockerStack } from '../models/docker-stack.model';
import { MenuService } from '../../services/menu.service';
import { ActivatedRoute } from '@angular/router';
import { MetricsService } from '../../metrics/services/metrics.service';

@Component({
  selector: 'app-services',
  templateUrl: './docker-services.component.html',
  styleUrls: ['./docker-services.component.css'],
  providers: [ ListService ]
})
export class DockerServicesComponent implements OnInit, OnDestroy {
  routeSub : any
  timer : any

  constructor(
    public listService : ListService,
    public dockerServicesService : DockerServicesService,
    public dockerStacksService : DockerStacksService,
    public menuService : MenuService,
    private route: ActivatedRoute,
    private metricsService : MetricsService) {
    listService.setFilterFunction(dockerServicesService.match)
  }

  ngOnInit() {
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
          this.loadServices()
        }
      )
      this.loadServices()
    })
  }

  ngOnDestroy() {
    this.routeSub.unsubscribe();
  }

  loadServices() {
    this.dockerServicesService.loadServices()
    if (this.menuService.autoRefresh) {
      this.timer = setInterval( () => {
          this.dockerServicesService.loadServices()
        }, 3000
      )
      return;
    }
    if (this.timer) {
      clearInterval(this.timer);
    }
  }

  containerList(serviceId : string) {
    this.dockerServicesService.setCurrentService(serviceId)
    this.menuService.navigate(["/amp", "stacks", this.dockerStacksService.currentStack.name, "services", serviceId, "containers"])
  }

  selectService(name : string) {
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

}
