import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import { OrganizationsService } from '../../services/organizations.service'
import { ListService } from '../../services/list.service';
import { UsersService } from '../../services/users.service';
import { User } from '../../models/user.model';
import { Member } from '../../models/member.model';
import { Organization } from '../../models/organization.model';
import { Team } from '../../models/team.model';
import { ActivatedRoute } from '@angular/router';
import { MenuService } from '../../services/menu.service';
import { HttpService } from '../../services/http.service';

@Component({
  selector: 'app-organization',
  templateUrl: './organization.component.html',
  styleUrls: ['./organization.component.css']


})
export class OrganizationComponent implements OnInit, OnDestroy {
  organization : Organization
  name = ""
  routeSub : any
  modeCreation : boolean = false
  public listUserService : ListService = new ListService()
  public listUserAddedService : ListService = new ListService()
  addedUsers : Member[] = []
  users : Member[] = []
  initialUserList : Member[] = []
  updated = false
  nbSaveInProgress = 0
  @ViewChild ('f') form: NgForm;

  constructor(
    private route: ActivatedRoute,
    private usersService: UsersService,
    public organizationsService : OrganizationsService,
    private menuService : MenuService,
    private httpService : HttpService) {
    this.listUserService.setFilterFunction(this.match)
    this.listUserAddedService.setFilterFunction(this.match)
  }

  ngOnInit() {
    this.menuService.setItemMenu('organization', 'Edit')
    this.routeSub = this.route.params.subscribe(params => {
      this.name = params['orgName'];
      for (let org of this.organizationsService.organizations) {
        if (org.name == this.name) {
          this.organization = org
        }
      }
      if (this.organization) {
        this.usersService.onUsersLoaded.subscribe(
          () => {
            this.initialUserList = this.usersService.getAllNoMembers(this.organization.members)
            this.addedUsers = this.organization.members.slice()
            this.users = this.initialUserList.slice()
            this.listUserAddedService.setData(this.addedUsers)
            this.listUserService.setData(this.users)
          }
        )
        this.usersService.loadUsers()
      }
    })
  }

  ngOnDestroy() {
    this.routeSub.unsubscribe();
  }

  match(item : Member, value : string) : boolean {
    if (value == '' || item.userName==='') {
      return true
    }
    if (item.userName && item.userName.includes(value)) {
      return true
    }
    return false
  }

  addUser( user : Member) {
    user.status++;
    user.saved = false
    let list : Member[] = []
    for (let item of this.users) {
      if (item.userName !== user.userName) {
        list.push(item)
      }
    }
    this.users=list
    this.listUserService.setData(this.users)
    this.addedUsers.push(user)
    this.listUserAddedService.setData(this.addedUsers)
    this.updated=this.isUpdated()
  }

  removeUser( user : Member) {
    user.status--;
    user.saved = false
    let list : Member[] = []
    for (let item of this.addedUsers) {
      if (item.userName !== user.userName) {
        list.push(item)
      }
    }
    this.addedUsers=list
    this.listUserAddedService.setData(this.addedUsers)
    this.users.push(user)
    this.listUserService.setData(this.users)
    this.updated=this.isUpdated()
  }

  addAll() {
    for (let user of this.users) {
      this.addUser(user)
    }
    this.updated=this.isUpdated()
  }

  removeAll() {
    for (let user of this.addedUsers) {
      this.removeUser(user)
    }
    this.updated=this.isUpdated()
  }

  isUpdated() : boolean {
    for (let user of this.users) {
      if (user.status !== 0) {
        return true
      }
    }
    for (let user of this.addedUsers) {
      if (user.status !== 0) {
        return true
      }
    }
    return false
  }

  applyUsers() {
    console.log("apply users")
    this.nbSaveInProgress=0
    this.menuService.waitingCursor(true)
    for (let user of this.users) {
      console.log(user)
      user.saved=false
      user.saveError=""
      if (user.status == -1) {
        console.log("removing user "+user.userName)
        this.nbSaveInProgress++
        this.httpService.removeUserFromOrganization(this.organization, user).subscribe(
          () => {
            user.saved=true
            user.status=0
            user.saveError=""
            this.decrSaveInProgress()
            console.log("done")
          },
          (error) => {
            this.addUser(user)
            user.saved=true
            user.saveError="save error"
            this.decrSaveInProgress()
            console.log(error)
          }
        )
      }
    }
    console.log("apply addedUsers")
    for (let user of this.addedUsers) {
      console.log(user)
      user.saved=false
      user.saveError=""
      if (user.status == 1) {
        console.log("adding user "+user.userName)
        this.nbSaveInProgress++
        this.httpService.addUserToOrganization(this.organization, user).subscribe(
          () => {
            console.log("done")
            user.saved=true
            user.saveError=""
            user.status=0
            this.decrSaveInProgress()
          },
          (error) => {
            console.log(error)
            this.removeUser(user)
            user.saved=true
            user.saveError="save error"
            this.decrSaveInProgress()
          }
        )
      }
    }
    if (this.nbSaveInProgress === 0) {
      this.menuService.waitingCursor(false)
    }
    this.organization.members = this.addedUsers.slice()
    this.updated=false
  }

  private decrSaveInProgress() {
    this.nbSaveInProgress--
    if (this.nbSaveInProgress === 0) {
      this.menuService.waitingCursor(false)
    }
  }

}
