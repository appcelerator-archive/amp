import { Component, OnInit } from '@angular/core';
import { User } from '../models/user.model';
import { UsersService } from '../services/users.service'
import { ListService } from '../services/list.service';
import { MenuService } from '../services/menu.service';
import { HttpService } from '../services/http.service';

@Component({
  selector: 'app-user-list',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.css'],
  providers: [ ListService ]
})
export class UsersComponent implements OnInit {
  currentUser : User
  createMode = false
  emptyUser : User = new User("", "", "")
  selectedUser : User = this.emptyUser
  message = ""

  constructor(
    public usersService : UsersService,
    public listService : ListService,
    private menuService : MenuService,
    private httpService : HttpService) {
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

  selectUser(name : string) {
    this.selectedUser = this.emptyUser
    for (let user of this.usersService.users) {
      if (user.name == name) {
        this.selectedUser = user
      }
    }
  }

  removeUser() {
    this.httpService.removeUser(this.selectedUser.name).subscribe(
      data => {},
      error => {
        let data = error.json()
        this.message = data.error
      }
    )
  }

}
