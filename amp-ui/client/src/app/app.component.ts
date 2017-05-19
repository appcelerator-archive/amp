import { Component, OnInit } from '@angular/core';
import { MenuService } from './services/menu.service';



@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})

export class AppComponent {

  constructor(public menuService : MenuService) { }

  ngOnInit() {
    this.menuService.waitingCursor(false)
  }
}
