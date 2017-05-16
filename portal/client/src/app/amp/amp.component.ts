import { Component, OnInit } from '@angular/core';
import { MenuService } from '../services/menu.service'

@Component({
  selector: 'app-amp',
  templateUrl: './amp.component.html',
  styleUrls: ['./amp.component.css']
})
export class AmpComponent implements OnInit {

  constructor(public menuService : MenuService) { }

  ngOnInit() {
  }

}
