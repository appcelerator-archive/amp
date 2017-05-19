import { Component, OnInit } from '@angular/core';
import { ListService } from '../../services/list.service';
import { DockerStacksService } from '../../services/docker-stacks.service';
import { DockerContainersService } from '../../services/docker-containers.service';
import { DockerServicesService } from '../../services/docker-services.service';
import { DockerContainer } from '../../models/docker-container.model';
import { MenuService } from '../../services/menu.service';

@Component({
  selector: 'app-containers',
  templateUrl: './docker-containers.component.html',
  styleUrls: ['./docker-containers.component.css'],
  providers: [ ListService ]
})
export class DockerContainersComponent implements OnInit {

  constructor(
    public listService : ListService,
    public dockerContainersService : DockerContainersService,
    public dockerServicesService : DockerServicesService,
    public menuService : MenuService) {
    listService.setFilterFunction(dockerContainersService.match)
  }

  ngOnInit() {
    this.menuService.setItemMenu('containers', 'List')
    this.listService.setData(this.dockerContainersService.containers)
  }

  selectContainer(id : string) {
    this.dockerContainersService.setCurrentContainer(id)
  }

  returnBack() {
    console.log("return back form container: "+this.dockerServicesService.currentService.id)
    this.menuService.navigate(["/amp", "stacks", this.dockerServicesService.currentService.id])
  }

}
