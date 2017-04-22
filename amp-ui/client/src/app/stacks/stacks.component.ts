import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Stack } from '../models/stack.model';
import { StacksService } from '../services/stacks.service';
import { ListService } from '../services/list.service';

@Component({
  selector: 'app-stacks',
  templateUrl: './stacks.component.html',
  styleUrls: ['./stacks.component.css'],
  providers: [ ListService ]
})
export class StacksComponent implements OnInit {
  currentStack : Stack

  constructor(
    private route : ActivatedRoute,
    public stacksService : StacksService,
    public listService : ListService) {
      listService.setFilterFunction(stacksService.match)
      listService.setData(this.stacksService.stacks)
    }

  ngOnInit() {
    let name = this.route.snapshot.params['name']
    //this.route.snapshot.queryParams
    //this.route.snapshot.queryFragment
    this.currentStack = new Stack('', name, 0, '', '')
    //this.route.queryParams.subscribe()
    //this.route.queryFragment.subscribe()
    this.route.params.subscribe( //automatically unsubscribed by A on component destroy
      (params : Params) => {
        this.currentStack = new Stack('', name, 0, '', '')
      }
    );
  }
}
