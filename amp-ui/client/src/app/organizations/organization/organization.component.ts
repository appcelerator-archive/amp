import { Component, OnInit, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import { OrganizationsService } from '../../services/organizations.service'
import { ListService } from '../../services/list.service';
import { User } from '../../models/user.model';

@Component({
  selector: 'app-organization',
  templateUrl: './organization.component.html',
  styleUrls: ['./organization.component.css']


})
export class OrganizationComponent implements OnInit {
  modeCreation : boolean = false
  public listUserService : ListService = new ListService()
  public listUserAddedService : ListService = new ListService()
  addedUsers : User[] = []
  users : User[] = []
  @ViewChild ('f') form: NgForm;

  constructor(public organizationsService : OrganizationsService) {
    this.listUserService.setFilterFunction(this.match)
    this.listUserAddedService.setFilterFunction(this.match)
  }

  ngOnInit() {
    this.users.push(new User("user1",'','User'))
    this.users.push(new User("user2",'','Owner'))
    this.users.push(new User("user3",'','User'))
    this.users.push(new User("user4",'','User'))
    this.users.push(new User("user5",'','User'))
    this.listUserAddedService.setData(this.addedUsers)
    this.listUserService.setData(this.users)
  }

  match(item : User, value : string) : boolean {
    if (value == '') {
      return true
    }
    if (item.name && item.name.includes(value)) {
      return true
    }
    return false
  }

  onDrop( user : User) {
    let list : User[] = []
    for (let item of this.users) {
      if (item.name !== user.name) {
        list.push(item)
      }
    }
    this.users=list
    this.listUserService.setData(this.users)
    this.addedUsers.push(user)
    this.listUserAddedService.setData(this.addedUsers)
  }
}
