import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import { OrganizationsService } from '../../../services/organizations.service'
import { ListService } from '../../../services/list.service';
import { UsersService } from '../../../services/users.service';
import { User } from '../../../models/user.model';
import { Member } from '../../../models/member.model';
import { Team } from '../../../models/team.model';
import { Organization } from '../../../models/organization.model';
import { TeamResource } from '../../../models/team-resource.model';
import { ActivatedRoute } from '@angular/router';
import { MenuService } from '../../../services/menu.service';
import { HttpService } from '../../../services/http.service';

@Component({
  selector: 'app-organization',
  templateUrl: './team.component.html',
  styleUrls: ['./team.component.css']
})

export class TeamComponent implements OnInit, OnDestroy {
  routeSub : any
  modeCreation : boolean = false
  userResourceToggle = false
  public listUserService : ListService = new ListService()
  public listUserAddedService : ListService = new ListService()
  public listResourceService : ListService = new ListService()
  organization : Organization
  team : Team
  addedUsers : Member[] = []
  users : Member[] = []
  initialUserList : Member[] = []
  resources : TeamResource[] = []
  permisionLabel : string[] = ['node', 'read', 'write', 'admin']
  updated = false;
  nbSaveInProgress = 0;
  @ViewChild ('f') form: NgForm;

  constructor(
    private route: ActivatedRoute,
    public organizationsService : OrganizationsService,
    private usersService : UsersService,
    private menuService : MenuService,
    private httpService : HttpService) {
    this.listUserService.setFilterFunction(this.matchUser)
    this.listUserAddedService.setFilterFunction(this.matchUser)
    this.listResourceService.setFilterFunction(this.matchResource)
  }

  ngOnInit() {
    this.menuService.setItemMenu('organization', 'Team edit')
    this.routeSub = this.route.params.subscribe(params => {
      let orgName = params['orgName'];
      let name = params['teamName'];
      for (let org of this.organizationsService.organizations) {
        if (org.name == orgName) {
          this.organization = org
        }
      }
      if (this.organization) {
        for (let team of this.organization.teams) {
          if (team.name == name) {
            this.team = team
          }
        }
        if (this.team) {
          this.initialUserList = this.usersService.getAllNoMembers(this.team.members)
          this.addedUsers = this.team.members.slice()
          this.resources = this.organization.resources.slice()
          //
          this.users = this.initialUserList.slice()
          this.listUserAddedService.setData(this.addedUsers)
          this.listUserService.setData(this.users)
          this.listResourceService.setData(this.resources)
        }
      }
    })
  }

  ngOnDestroy() {
    this.routeSub.unsubscribe();
  }

  matchUser(item : Member, value : string) : boolean {
    if (value == '' || item.userName==='') {
      return true
    }
    if (item.userName && item.userName.includes(value)) {
      return true
    }
    return false
  }

  matchResource(item : TeamResource, value : string) : boolean {
    if (value == '') {
      return true
    }
    if (item.id && item.id.includes(value)) {
      return true
    }
    if (item.type && item.type.includes(value)) {
      return true
    }
    if (item.name && item.name.includes(value)) {
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

  userManagement() {
    this.userResourceToggle=false
  }

  resourceManagement(){
    this.userResourceToggle=true
  }

  returnBack() {
    this.menuService.navigate(['/amp', 'organizations', this.organization.name])
  }

  removeTeam() {
    this.menuService.waitingCursor(true)
    this.httpService.deleteTeam(this.organization, this.team).subscribe(
      () => {
        this.menuService.waitingCursor(false)
        let list : Team[] = []
        for (let team of this.organization.teams) {
          if (team.name !== this.team.name) {
            list.push(team)
          }
        }
        this.organization.teams=list
        this.menuService.navigate(['/amp', 'organizations', this.organization.name])
      },
      (error) => {
        this.menuService.waitingCursor(false)
        console.log(error)
      }
    )
  }

  setPermission(res : TeamResource, level : number) {
    res.setPermission(level)
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
        this.httpService.removeUserFromTeam(this.organization, this.team, user).subscribe(
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
        this.httpService.addUserToTeam(this.organization, this.team, user).subscribe(
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
    this.team.members = this.addedUsers.slice()
    this.updated=false
  }

  private decrSaveInProgress() {
    this.nbSaveInProgress--
    if (this.nbSaveInProgress === 0) {
      this.menuService.waitingCursor(false)
    }
  }
}
