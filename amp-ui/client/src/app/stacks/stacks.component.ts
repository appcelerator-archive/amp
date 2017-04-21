import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Stack } from '../models/stack.model';

@Component({
  selector: 'app-stacks',
  templateUrl: './stacks.component.html',
  styleUrls: ['./stacks.component.css']
})
export class StacksComponent implements OnInit {
  currentStack : Stack

  constructor(private route : ActivatedRoute) { }

  ngOnInit() {
    let name = this.route.snapshot.params['name']
    //this.route.snapshot.queryParams
    //this.route.snapshot.queryFragment
    this.currentStack = new Stack('', name, 0, '')
    //this.route.queryParams.subscribe()
    //this.route.queryFragment.subscribe()
    this.route.params.subscribe( //automatically unsubscribed by A on component destroy
      (params : Params) => {
        this.currentStack = new Stack('', name, 0, '')
      }
    );
  }

}
