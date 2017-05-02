import { Component, OnInit, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import { OrganizationsService } from '../../../services/organizations.service'
import { ListService } from '../../../services/list.service';
import { User } from '../../../models/user.model';
import { Organization } from '../../../models/organization.model';
import { TeamResource } from '../../../models/team-resource.model';
import { ActivatedRoute } from '@angular/router';
import { MenuService } from '../../../services/menu.service';

@Component({
  selector: 'app-organization',
  templateUrl: './team.component.html',
  styleUrls: ['./team.component.css']
})

export class TeamComponent implements OnInit {
  orgName = ""
  name = ""
  routeSub : any
  modeCreation : boolean = false
  userResourceToggle = false
  public listUserService : ListService = new ListService()
  public listUserAddedService : ListService = new ListService()
  public listResourceService : ListService = new ListService()
  organization : Organization
  addedUsers : User[] = []
  users : User[] = []
  initialUserList : User[] = []
  resources : TeamResource[] = []
  permisionLabel : string[] = ['node', 'read', 'write', 'admin']
  updated = false
  @ViewChild ('f') form: NgForm;

  constructor(
    private route: ActivatedRoute,
    public organizationsService : OrganizationsService,
    private menuService : MenuService) {
    this.listUserService.setFilterFunction(this.matchUser)
    this.listUserAddedService.setFilterFunction(this.matchUser)
    this.listResourceService.setFilterFunction(this.matchResource)
  }

  ngOnInit() {
    this.menuService.setItemMenu('organization', 'Team edit')
    this.routeSub = this.route.params.subscribe(params => {
      this.orgName = params['orgName'];
      this.name = params['teamName'];
      for (let org of this.organizationsService.organizations) {
        if (org.name == this.orgName) {
          this.organization = org
        }
      }
      if (this.organization) {
        console.log("Team: "+this.orgName+":"+this.name)
        this.initialUserList = this.organization.members.slice()
        this.addedUsers = []
        this.resources = this.organization.resources.slice()
        //
        this.users = this.initialUserList.slice()
        this.listUserAddedService.setData(this.addedUsers)
        this.listUserService.setData(this.users)
        this.listResourceService.setData(this.resources)
      }
    })
  }

  matchUser(item : User, value : string) : boolean {
    if (value == '' || item.name==='') {
      return true
    }
    if (item.name && item.name.includes(value)) {
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

  addUser( user : User) {
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
    this.updated=true
  }

  removeUser( user : User) {
    let list : User[] = []
    for (let item of this.addedUsers) {
      if (item.name !== user.name) {
        list.push(item)
      }
    }
    this.addedUsers=list
    this.listUserAddedService.setData(this.addedUsers)
    this.users.push(user)
    this.listUserService.setData(this.users)
    if (this.users.length === this.initialUserList.length) {
      this.updated=false
    }
  }

  addAll() {
    this.addedUsers=this.initialUserList.slice()
    this.users=[]
    this.listUserAddedService.setData(this.addedUsers)
    this.listUserService.setData(this.users)
    this.updated=true
  }

  removeAll() {
    this.updated=false
    this.users = this.initialUserList.slice()
    this.addedUsers=[]
    this.listUserAddedService.setData(this.addedUsers)
    this.listUserService.setData(this.users)
  }

  applyUsers() {
    console.log("apply")
    //apply
    this.updated=false
  }

  userManagement() {
    this.userResourceToggle=false
  }

  resourceManagement(){
    this.userResourceToggle=true
  }

  returnBack() {
    console.log(this.orgName)
    this.menuService.navigate(['/amp/organizations', this.orgName])
  }

  setPermission(res : TeamResource, level : number) {
    res.setPermission(level)
  }
}
