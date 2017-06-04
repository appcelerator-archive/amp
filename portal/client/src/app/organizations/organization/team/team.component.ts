import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import { OrganizationsService } from '../../../organizations/services/organizations.service'
import { ListService } from '../../../services/list.service';
import { UsersService } from '../../../services/users.service';
import { User } from '../../../models/user.model';
import { Member } from '../../../organizations/models/member.model';
import { Team } from '../../../organizations/models/team.model';
import { Organization } from '../../../organizations/models/organization.model';
import { OrganizationResource } from '../../../organizations/models/organization-resource.model';
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
  teamMode = 'user'
  public listUserService : ListService = new ListService()
  public listUserAddedService : ListService = new ListService()
  public listResourceService : ListService = new ListService()
  public listResourceAddedService : ListService = new ListService()
  organization : Organization = new Organization("", "")
  team : Team = new Team("")
  addedUsers : Member[] = []
  users : Member[] = []
  initialUserList : Member[] = []
  addedResources : OrganizationResource[] = []
  resources : OrganizationResource[] = []
  initialResourceList : OrganizationResource[] = []
  permisionLabel : string[] = ['read', 'write', 'admin']
  userUpdated = false;
  resUpdated = false
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
    this.listResourceAddedService.setFilterFunction(this.matchResource)
  }

  ngOnInit() {
    this.menuService.setItemMenu('organization', 'Team edit')
    this.routeSub = this.route.params.subscribe(params => {
      let orgName = this.organizationsService.currentOrganization.name;
      this.organization = this.organizationsService.currentOrganization
      let name = params['teamName'];
      for (let team of this.organization.teams) {
        if (team.name == name) {
          this.team = team
        }
      }
      if (this.team) {
        //Users
        this.initialUserList = this.organizationsService.getAllNoMembers(this.team.members)
        this.addedUsers = this.team.members.slice()
        this.initialResourceList = this.organizationsService.getAllNoResources(this.team.resources)
        this.addedResources = this.team.resources.slice()
        //
        this.users = this.initialUserList.slice()
        this.listUserAddedService.setData(this.addedUsers)
        this.listUserService.setData(this.users)
        //
        this.resources = this.initialResourceList.slice()
        this.listResourceAddedService.setData(this.addedResources)
        this.listResourceService.setData(this.resources)
      }
    })
  }

  ngOnDestroy() {
    this.routeSub.unsubscribe();
  }

  userManagement() {
    this.teamMode='user'
  }

  resourceManagement(){
    this.teamMode='resource'
  }

  authorizationManagement(){
    this.teamMode='authorization'
  }

  returnBack() {
    this.menuService.navigate(['/amp', 'organizations', this.organization.name])
  }

  removeTeam() {
    if(confirm("Are you sure to delete the team: "+this.team.name)) {
      this.menuService.waitingCursor(true)
      this.httpService.deleteTeam(this.organization.name, this.team.name).subscribe(
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
  }

  //----------------------------------------------------------------------------
  // Users add-remove management

  matchUser(item : Member, value : string) : boolean {
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
    this.userUpdated=this.isUserUpdated()
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
    this.userUpdated=this.isUserUpdated()
  }

  addAllUsers() {
    for (let user of this.users) {
      this.addUser(user)
    }
    this.userUpdated=this.isUserUpdated()
  }

  removeAllUsers() {
    for (let user of this.addedUsers) {
      this.removeUser(user)
    }
    this.userUpdated=this.isUserUpdated()
  }

  isUserUpdated() : boolean {
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
    this.nbSaveInProgress=0
    this.menuService.waitingCursor(true)
    for (let user of this.users) {
      user.saved=false
      user.saveError=""
      if (user.status == -1) {
        this.nbSaveInProgress++
        this.httpService.removeUserFromTeam(this.organization.name, this.team.name, user.userName).subscribe(
          () => {
            user.saved=true
            user.status=0
            user.saveError=""
            this.decrSaveUsersInProgress()
          },
          (error) => {
            console.log(error)
            try {
              let data = JSON.parse(error._body)
                user.saveError=data.error
            } catch (errorj) {
              console.log(errorj)
            }
            this.addUser(user)
            user.saved=true
            user.saveError=error
            this.decrSaveUsersInProgress()
          }
        )
      }
    }
    for (let user of this.addedUsers) {
      user.saved=false
      user.saveError=""
      if (user.status == 1) {
        this.nbSaveInProgress++
        this.httpService.addUserToTeam(this.organization.name, this.team.name, user.userName).subscribe(
          () => {
            user.saved=true
            user.saveError=""
            user.status=0
            this.decrSaveUsersInProgress()
          },
          (error) => {
            console.log(error)
            try {
              let data = JSON.parse(error._body)
                user.saveError=data.error
            } catch (errorj) {
              console.log(errorj)
            }
            this.removeUser(user)
            user.saved=true
            this.decrSaveUsersInProgress()
          }
        )
      }
    }
    if (this.nbSaveInProgress === 0) {
      this.menuService.waitingCursor(false)
    }
    this.team.members = this.addedUsers.slice()
    this.userUpdated=false
  }

  private decrSaveUsersInProgress() {
    this.nbSaveInProgress--
    if (this.nbSaveInProgress === 0) {
      this.menuService.waitingCursor(false)
    }
  }

  //----------------------------------------------------------------------------
  // resource add-remove management

  matchResource(item : OrganizationResource, value : string) : boolean {
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

  addResource( res : OrganizationResource) {
    res.status++;
    res.saved = false
    let list : OrganizationResource[] = []
    for (let item of this.resources) {
      if (item.id !== res.id) {
        list.push(item)
      }
    }
    this.resources=list
    this.listResourceService.setData(this.resources)
    this.addedResources.push(res)
    this.listResourceAddedService.setData(this.addedResources)
    this.resUpdated=this.isResourceUpdated()
  }

  removeResource( res : OrganizationResource) {
    res.status--;
    res.saved = false
    let list : OrganizationResource[] = []
    for (let item of this.addedResources) {
      if (item.id !== res.id) {
        list.push(item)
      }
    }
    this.addedResources=list
    this.listResourceAddedService.setData(this.addedResources)
    this.resources.push(res)
    this.listResourceService.setData(this.resources)
    this.resUpdated=this.isResourceUpdated()
  }

  addAllResources() {
    for (let res of this.resources) {
      this.addResource(res)
    }
    this.resUpdated=this.isResourceUpdated()
  }

  removeAllResources() {
    for (let res of this.addedResources) {
      this.removeResource(res)
    }
    this.resUpdated=this.isResourceUpdated()
  }

  isResourceUpdated() : boolean {
    for (let res of this.resources) {
      if (res.status !== 0) {
        return true
      }
    }
    for (let res of this.addedResources) {
      if (res.status !== 0) {
        return true
      }
    }
    return false
  }

  applyResources() {
    console.log("apply resources")
    this.nbSaveInProgress=0
    this.menuService.waitingCursor(true)
    for (let res of this.resources) {
      res.saved=false
      res.saveError=""
      if (res.status == -1) {
        console.log("removing resources "+res.name)
        this.nbSaveInProgress++
        this.httpService.removeResourceFromTeam(this.organization.name, this.team.name, res.id).subscribe(
          () => {
            res.saved=true
            res.status=0
            res.saveError=""
            this.decrSaveResourcesInProgress()
            console.log("done")
          },
          (error) => {
            console.log(error)
            try {
              let data = JSON.parse(error._body)
                res.saveError=data.error
            } catch (errorj) {
              console.log(errorj)
            }
            this.addResource(res)
            res.saved=true
            res.saveError=error
            this.decrSaveResourcesInProgress()
          }
        )
      }
    }
    console.log("apply addedResources")
    for (let res of this.addedResources) {
      res.saved=false
      res.saveError=""
      if (res.status == 1) {
        console.log("adding resource "+res.name)
        this.nbSaveInProgress++
        this.httpService.addResourceToTeam(this.organization.name, this.team.name, res.id).subscribe(
          () => {
            console.log("done")
            res.saved=true
            res.saveError=""
            res.status=0
            this.decrSaveResourcesInProgress()
          },
          (error) => {
            console.log(error)
            try {
              let data = JSON.parse(error._body)
                res.saveError=data.error
            } catch (errorj) {
              console.log(errorj)
            }
            this.removeResource(res)
            res.saved=true
            this.decrSaveResourcesInProgress()
          }
        )
      }
    }
    if (this.nbSaveInProgress === 0) {
      this.menuService.waitingCursor(false)
    }
    this.team.resources = this.addedResources.slice()
    this.resUpdated=false
  }

  private decrSaveResourcesInProgress() {
    this.nbSaveInProgress--
    if (this.nbSaveInProgress === 0) {
      this.menuService.waitingCursor(false)
    }
  }

  //----------------------------------------------------------------------------
  // Authorization management

  setAuthorization(res : OrganizationResource, level : number) {
    this.httpService.changeTeamResourcePermissionLevel(
      this.organization.name, this.team.name, res.id, level).subscribe(
        () => {
          res.changeAuth = true
          res.changeAuthError = ''
          res.setAuthorization(level)
        },
        (err) => {
          let error = err.json()
          res.changeAuth = true
          res.changeAuthError = error.error
        }
      )
  }

}
