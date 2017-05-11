import { Component, OnInit, OnDestroy } from '@angular/core';
import { ListService } from '../../services/list.service';
import { DockerServicesService } from '../services/docker-services.service';
import { DockerStacksService } from '../services/docker-stacks.service';
import { DockerService } from '../models/docker-service.model';
import { DockerStack } from '../models/docker-stack.model';
import { MenuService } from '../../services/menu.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-services',
  templateUrl: './docker-services.component.html',
  styleUrls: ['./docker-services.component.css'],
  providers: [ ListService ]
})
export class DockerServicesComponent implements OnInit, OnDestroy {
  routeSub : any
  timer : any
  currentService : DockerService = new DockerService("", "", "", "", "")
  stack : DockerStack = new DockerStack("", "", 0, "", "")

  constructor(
    public listService : ListService,
    public dockerServicesService : DockerServicesService,
    private dockerStacksService : DockerStacksService,
    public menuService : MenuService,
    private route: ActivatedRoute) {
    listService.setFilterFunction(dockerServicesService.match)
  }

  ngOnInit() {
    this.menuService.setItemMenu('services', 'List')
    this.routeSub = this.route.params.subscribe(params => {
      let stackName = params['stackName'];
      for (let stack of this.dockerStacksService.stacks) {
        if (stack.name == stackName) {
          this.stack = stack
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
    this.menuService.navigate(["/amp", "stacks", this.stack.name, "services", serviceId, "containers"])
  }

  selectService(id : string) {
    console.log("select service: "+id)
    this.dockerServicesService.setCurrentService(id)
  }

  returnBack() {
    console.log("return back form service: "+this.dockerStacksService.currentStack.id)
    this.menuService.navigate(["/amp", "stacks"])
  }

  metrics(serviceId : string) {
    this.menuService.navigate(['/amp', 'metrics', 'service', serviceId])
  }


}
