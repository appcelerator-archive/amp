import { Component, OnInit, OnDestroy} from '@angular/core';
import { ListService } from '../../services/list.service';
import { DockerStacksService } from '../services/docker-stacks.service';
import { DockerContainersService } from '../services/docker-containers.service';
import { DockerServicesService } from '../services/docker-services.service';
import { DockerContainer } from '../models/docker-container.model';
import { DockerService } from '../models/docker-service.model';
import { MenuService } from '../../services/menu.service';
import { ActivatedRoute } from '@angular/router';
import { MetricsService } from '../../metrics/services/metrics.service';

@Component({
  selector: 'app-containers',
  templateUrl: './docker-containers.component.html',
  styleUrls: ['./docker-containers.component.css'],
  providers: [ ListService ]
})
export class DockerContainersComponent implements OnInit, OnDestroy {
  routeSub : any
  timer : any
  currentContainer : DockerContainer = new DockerContainer("", "", "", "", "")
  service = new DockerService("", "", "", "", "")
  stackName = ""

  constructor(
    public listService : ListService,
    public dockerContainersService : DockerContainersService,
    private dockerServicesService : DockerServicesService,
    public menuService : MenuService,
    private route: ActivatedRoute,
    private metricsService : MetricsService) {
      listService.setFilterFunction(dockerContainersService.match)
  }

  ngOnInit() {
    this.menuService.setItemMenu('containers', 'List')
    this.routeSub = this.route.params.subscribe(params => {
      this.stackName = params['stackName'];
      let serviceId = params['serviceId'];
      for (let serv of this.dockerServicesService.services) {
        if (serv.id == serviceId) {
          this.service = serv
          this.dockerServicesService.currentService = serv
        }
      }
      this.dockerContainersService.onContainersLoaded.subscribe(
        () => {
          this.listService.setData(this.dockerContainersService.containers)
        }
      )
      this.menuService.onRefreshClicked.subscribe(
        () => {
          this.loadContainers()
        }
      )
      this.loadContainers()
    })
  }

  ngOnDestroy() {
    this.routeSub.unsubscribe();
  }

  loadContainers() {
    this.dockerContainersService.loadContainers()
    if (this.menuService.autoRefresh) {
      this.timer = setInterval( () => {
          this.dockerContainersService.loadContainers()
        }, 3000
      )
      return;
    }
    if (this.timer) {
      clearInterval(this.timer);
    }
  }

  selectContainer(id : string) {
    this.dockerContainersService.setCurrentContainer(id)
  }

  returnBack() {
    this.menuService.returnToPreviousPath()
  }

  metrics(containerId : string) {
    this.menuService.navigate(['/amp', 'metrics', 'task', 'single', containerId])
  }

}
