import { Component, OnInit } from '@angular/core';
import { ListService } from '../../services/list.service';
import { DockerServicesService } from '../../services/docker-services.service';
import { DockerStacksService } from '../../services/docker-stacks.service';
import { DockerService } from '../../models/docker-service.model';
import { MenuService } from '../../services/menu.service';

@Component({
  selector: 'app-services',
  templateUrl: './docker-services.component.html',
  styleUrls: ['./docker-services.component.css'],
  providers: [ ListService ]
})
export class DockerServicesComponent implements OnInit {
  currentService : DockerService
  constructor(
    public listService : ListService,
    public dockerServicesService : DockerServicesService,
    public dockerStacksService : DockerStacksService,
    public menuService : MenuService,) {
    listService.setFilterFunction(dockerServicesService.match)
  }

  ngOnInit() {
    this.menuService.setItemMenu('services', 'List')
    this.listService.setData(this.dockerServicesService.services)

  }

  containerList(serviceId : string) {
    this.dockerServicesService.setCurrentService(serviceId)
    this.menuService.navigate(["/amp/services", serviceId])
  }

  selectService(id : string) {
    console.log("select service: "+id)
    this.dockerServicesService.setCurrentService(id)
  }

  returnBack() {
    console.log("return back form service: "+this.dockerStacksService.currentStack.id)
    this.menuService.navigate(["/amp", "stacks"])
  }

}
