import { Component, OnInit } from '@angular/core';
import { User } from '../models/user.model';
import { UsersService } from '../services/users.service'
import { ListService } from '../services/list.service';
import { MenuService } from '../services/menu.service';

@Component({
  selector: 'app-user-list',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.css'],
  providers: [ ListService ]
})
export class UsersComponent implements OnInit {
  currentUser : User
  createMode = false

  constructor(
    public usersService : UsersService,
    public listService : ListService,
    private menuService : MenuService) {
    listService.setFilterFunction(usersService.match)
  }


  ngOnInit() {
    this.menuService.setItemMenu('users', 'List')
    this.usersService.onUsersLoaded.subscribe(
      () => {
        this.listService.setData(this.usersService.users)
      }
    )
    this.usersService.loadUsers()
  }

  setCreateMode(mode: boolean) {
    this.createMode = mode
  }

  selectAllItems() {
    for (let user of this.usersService.users) {
      user.checked=true
    }
  }

}
