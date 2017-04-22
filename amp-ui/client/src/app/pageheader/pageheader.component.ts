import { Component, OnInit } from '@angular/core';
import { User } from '../models/user.model';
import { UsersService } from '../services/users.service';
import { MenuService } from '../services/menu.service';

@Component({
  selector: 'app-pageheader',
  templateUrl: './pageheader.component.html',
  styleUrls: ['./pageheader.component.css']
})

export class PageheaderComponent implements OnInit {
  menuTitle = "title"
  menuItem = "item"

  constructor(public usersService : UsersService, public menuService : MenuService) {}

  ngOnInit() {
  }

  

}
