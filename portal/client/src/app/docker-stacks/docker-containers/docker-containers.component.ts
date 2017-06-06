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
  currentContainer : DockerContainer = new DockerContainer("", "", "", "", "")


  constructor(
    public listService : ListService,
    public dockerContainersService : DockerContainersService,
    public dockerServicesService : DockerServicesService,
    public dockerStacksService : DockerStacksService,
    public menuService : MenuService,
    private route: ActivatedRoute,
    private metricsService : MetricsService) {
      listService.setFilterFunction(dockerContainersService.match)
  }

  ngOnInit() {
    this.menuService.setItemMenu('containers', 'List')
    this.routeSub = this.route.params.subscribe(params => {
      let stackName = params['stackName'];
      this.dockerStacksService.setCurrentStack(stackName)
      let serviceId = params['serviceId'];
      this.dockerServicesService.setCurrentServiceById(serviceId)

      this.dockerContainersService.onContainersLoaded.subscribe(
        () => {
          this.listService.setData(this.dockerContainersService.containers)
        }
      )
      this.menuService.onRefreshClicked.subscribe(
        () => {
          this.dockerContainersService.loadContainers(true)
        }
      )
      this.dockerContainersService.loadContainers(false)
    })
    this.menuService.setCurrentTimer(setInterval(
      () => {
      this.dockerContainersService.loadContainers(true)
      },
      3000
    ))
  }

  ngOnDestroy() {
    this.routeSub.unsubscribe();
    this.menuService.clearCurrentTimer()
  }

  getColor(cont : DockerContainer) : string {
    if (cont.state == 'running') {
      return 'limegreen'
    }
    if (cont.state == 'failed') {
      return 'red'
    }
    if (cont.state == 'shutdown') {
      return 'lightgrey'
    }
    return 'black';
  }

  selectContainer(id : string) {
    this.dockerContainersService.setCurrentContainer(id)
  }

  returnBack() {
    this.menuService.returnToPreviousPath()
  }

  metrics(taskId : string) {
    this.menuService.navigate(['/amp', 'metrics', 'task', 'single', taskId])
  }

  logs(taskId : string) {
    console.log(taskId)
    this.menuService.navigate(['/amp', 'logs', 'task', taskId])
  }
}
