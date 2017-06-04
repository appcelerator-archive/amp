import { Component, OnInit, OnDestroy } from '@angular/core';
import { User } from '../models/user.model';
import { UsersService } from '../services/users.service'
import { OrganizationsService } from '../organizations/services/organizations.service'
import { ListService } from '../services/list.service';
import { MenuService } from '../services/menu.service';
import { HttpService } from '../services/http.service';
import { ActivatedRoute } from '@angular/router';
import { AppWindow } from '../models/app-window.model';

@Component({
  selector: 'app-user-list',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.css'],
  providers: [ ListService ]
})

export class UsersComponent implements OnInit, OnDestroy {
  currentUser : User
  createMode = false
  emptyUser : User = new User("", "", "")
  selectedUser : User = this.emptyUser
  message = ""
  routeSub : any
  title = "users"
  orgName = ""
  changeRoleTitle = "Update role"
  listHeight = 400

  constructor(
    public usersService : UsersService,
    public organizationsService : OrganizationsService,
    public listService : ListService,
    private menuService : MenuService,
    private httpService : HttpService,
    private route: ActivatedRoute) {
    listService.setFilterFunction(usersService.match)
  }


  ngOnInit() {
    this.title=" Users"
    this.menuService.setItemMenu('users', 'List')
    this.resize(this.menuService.appWindow)
    this.menuService.onWindowResize.subscribe(
      (win) => {
        this.resize(win)
      }
    )
    this.routeSub = this.route.params.subscribe(params => {
      this.orgName = params['orgName'];
      if (this.orgName) {
        this.title=" Organization "+this.orgName+": users"
      }
      //update user list if users list reloaded
      this.usersService.onUsersLoaded.subscribe(
        () => {
          if (!this.orgName) {
            this.listService.setData(this.usersService.users)
          }
        }
      )
      //update orgName user list if organization users list reloaded
      this.organizationsService.onOrganizationsLoaded.subscribe(
        () => {
          if (this.orgName) {
            this.listService.setData(this.usersService.getUserList(this.orgName))
          }
        }
      )
      //reload users list on refresh click
      this.menuService.onRefreshClicked.subscribe(
        () => {
          if (this.orgName) {
            this.organizationsService.loadOrganizations(this.usersService.currentUser.name, true)
          } else {
            this.usersService.loadUsers(true)
          }
        }
      )
      //if user list alreday in memory, don't reload
      if (this.orgName) {
        //load orgName users only if not already in memory
        if (this.organizationsService.organizations.length == 0) {
          this.organizationsService.loadOrganizations(this.usersService.currentUser.name, true)
        } else {
          this.listService.setData(this.usersService.getUserList(this.orgName))
        }
      } else {
        //load all users only if not already in memory
        if (this.usersService.users.length == 0) {
          this.usersService.loadUsers(true)
        } else {
          this.listService.setData(this.usersService.users)
        }
      }
    })
  }

  ngOnDestroy() {
    this.routeSub.unsubscribe();
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
    if (this.orgName) {
      this.httpService.removeUserFromOrganization(this.orgName, this.selectedUser.name).subscribe(
        data => {
          this.organizationsService.loadOrganizations(this.usersService.currentUser.name, true)
        },
        error => {
          let data = error.json()
          this.message = data.error
        }
      )
    } else {
      this.httpService.removeUser(this.selectedUser.name).subscribe(
        data => {
          this.usersService.loadUsers(true)
        },
        error => {
          let data = error.json()
          this.message = data.error
        }
      )
    }
  }

  returnBack() {
    this.menuService.returnToPreviousPath()
  }

  resize(win : AppWindow) {
    this.listHeight = win.height - 240;
  }

  setUserRole(user : User, role : string) {
    if (user.member.role != role) {
      let nrole = 0
      if (role == 'Owner') {
        nrole = 1
      }
      this.httpService.changeOrganizationMemberRole(this.orgName, user.name, nrole).subscribe(
        () => {
          user.member.role=role
        },
        (err) => {
          let error = err.json()
          this.message = error.error
        }
      )
    }
  }

}
