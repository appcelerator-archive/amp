import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { DockerStack } from '../models/docker-stack.model';
import { DockerStacksService } from '../services/docker-stacks.service';
import { MenuService } from '../services/menu.service';
import { ListService } from '../services/list.service';
import {Observable} from 'rxjs/Observable';

@Component({
  selector: 'app-stacks',
  templateUrl: './docker-stacks.component.html',
  styleUrls: ['./docker-stacks.component.css'],
  providers: [ ListService ]
})
export class DockerStacksComponent implements OnInit {
  currentStack : DockerStack
  deployTitle = "Deploy"
  timer : any = null;

  constructor(
    private route : ActivatedRoute,
    public dockerStacksService : DockerStacksService,
    public listService : ListService,
    private menuService : MenuService) {
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
        this.loadStacks()
      }
    )
    this.loadStacks()
    let name = this.route.snapshot.params['name']
    this.currentStack = new DockerStack('', name, 0, '', '')
    this.route.params.subscribe( //automatically unsubscribed by A on component destroy
      (params : Params) => {
        this.currentStack = new DockerStack('', name, 0, '', '')
      }
    );
  }

  ngOnDestroy() {
    if (this.timer) {
      clearInterval(this.timer);
    }
    //this.dockerStacksService.onStacksLoaded.unsubscribe();
  }

  loadStacks() {
    this.dockerStacksService.loadStacks()
    if (this.menuService.autoRefresh) {
      this.timer = setInterval( () => {
          this.dockerStacksService.loadStacks()
        }, 3000
      )
      return;
    }
    if (this.timer) {
      clearInterval(this.timer);
    }
  }

  serviceList(stackId : string) {
    this.dockerStacksService.setCurrentStack(stackId)
    this.menuService.navigate(["/amp/stacks/", stackId])
  }

  selectStack(id : string) {
    let lastId = this.dockerStacksService.currentStack.id
    if (lastId == id) {
        this.dockerStacksService.setCurrentStack("")
        this.deployTitle="Deploy"
        return
    }
    this.deployTitle="Update"
    this.dockerStacksService.setCurrentStack(id)
  }

  deploy() {
    if (this.deployTitle == "Update") {
      this.update()
      return
    }
    this.menuService.navigate(["/amp/stacks/deploy"])
  }

  update() {
    let stackId = this.dockerStacksService.currentStack.id
    this.menuService.navigate(["/amp", "stacks", stackId, "deploy"])
  }

}
