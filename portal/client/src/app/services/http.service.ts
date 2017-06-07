import { Injectable } from '@angular/core';
import { Http, Headers, Response } from '@angular/http';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'RxJS/Rx';
import { User } from '../models/user.model';
import { Team } from '../organizations/models/team.model';
import { Organization } from '../organizations/models/organization.model';
import { Member } from '../organizations/models/member.model';
import { OrganizationResource } from '../organizations/models/organization-resource.model';
import { DockerStack } from '../docker-stacks/models/docker-stack.model';
import { DockerService } from '../docker-stacks/models/docker-service.model';
import { DockerContainer } from '../docker-stacks/models/docker-container.model';
import { StatsRequest } from '../models/stats-request.model';
import { GraphHistoricData } from '../models/graph-historic-data.model';
import { GraphCurrentData } from '../models/graph-current-data.model';
import { LogsRequest } from '../logs/models/logs-request.model';
import { Log } from '../logs/models/log.model';
import { Node } from '../nodes/models/node.model';
import { Dashboard } from '../dashboard/models/dashboard.model'
import * as d3 from 'd3';
import 'rxjs/add/operator/retrywhen';
import 'rxjs/add/operator/scan';
import 'rxjs/add/operator/delay';

const httpRetryDelay = 200
const httpRetryNumber = 3

@Injectable()
export class HttpService {
  private token = ""
  onHttpError = new Subject();
  parseTime = d3.timeParse("%Y-%m-%dT%H:%M:%S");
  //default dev debug url
  addr = "https://gw.local.appcelerator.io/v1"


  constructor(private http : Http) {
    if (window.location.host.substring(0,9)=="localhost") {
      console.log("Dev mode, Gateway url: "+this.addr)
      return
    }
    let host = "gw."+window.location.host
    this.addr=window.location.protocol +"//"+host+"/v1"
    console.log("Gateway url: "+this.addr)
  }

  users() {
    return this.httpGet("/users")
      .map((res : Response) => {
        const data = res.json()
        let list : User[] = []
        if (data.users) {
          for (let item of data.users) {
            let user = new User(item.name, item.email, "User")
            user.verified = item.is_verified
            user.tokenUsed = item.token_used
            user.createDate = this.formatedDate(item.create_dt)
            list.push(user)
          }
        }
        return list
      }
    );
  }

  getVersion() {
    return this.httpGet("/version")
  }

  changePassword(currentPwd : string, newPwd : string) {
    return this.httpPut("/users/password/change", { existingPassword: currentPwd, newPassword: newPwd});
  }

  retrieveLogin(email : string) {
    return this.httpPost("/users/"+email+"/reminder", { email : email});
  }

  resetPassword(name : string) {
    return this.httpPost("/users/"+name+"/password/reset", {name : name});
  }

  createOrganization(orgName : string, orgEmail : string) {
    return this.httpPost("/organizations", {name: orgName, email: orgEmail});
  }

  switchOrganization(orgName : string) {
    return this.httpPost("/switch", { account: orgName });
  }

  deleteOrganization(orgName : string) {
    return this.httpDelete("/organizations/"+orgName);
  }

  addUserToOrganization(orgName : string, memberName : string) {
    return this.httpPost("/organizations/"+orgName+"/members", {organization_name: orgName, user_name: memberName});
  }

  removeUserFromOrganization(orgName : string, memberName : string) {
    return this.httpDelete("/organizations/"+orgName+"/members/"+memberName);
  }

  createTeam(orgName : string, teamName : string) {
    return this.httpPost("/organizations/"+orgName+"/teams", {organization_name: orgName, team_name: teamName});
  }

  deleteTeam(orgName : string, teamName : string) {
    return this.httpDelete("/organizations/"+orgName+"/teams/"+teamName);
  }

  addUserToTeam(orgName : string, teamName : string, memberName : string) {
    return this.httpPost("/organizations/"+orgName+"/teams/"+teamName+"/members", {organization_name: orgName, team_name: teamName, user_name: memberName});
  }

  removeUserFromTeam(orgName : string, teamName : string, memberName : string) {
    return this.httpDelete("/organizations/"+orgName+"/teams/"+teamName+"/members/"+memberName);
  }

  userOrganization(userName : string) {
    //return this.http.get(this.addr+"/users/"+user.name+"/organizations", { headers: this.setHeaders() })
    return this.httpGet("/users/"+userName+"/organizations")
      .map((res : Response) => {
        const data = res.json()
        //console.log("data")
        //console.log(data)
        let list : Organization[] = []
        if (data.organizations) {
          for (let org of data.organizations) {
            let newOrg = new Organization(org.name, org.email)
            if (org.members) {
              for (let mem of org.members) {
                newOrg.members.push(new Member(mem.name, mem.role))
              }
            }
            if (org.teams) {
              for (let team of org.teams) {
                let newTeam = new Team(team.name)
                if (team.members) {
                  for (let mname of team.members) {
                    newTeam.members.push(new Member(mname, undefined))
                  }
                }
                if (team.resources) {
                  for (let res of team.resources) {
                    newTeam.resources.push(new OrganizationResource(res.id, "", ""))
                  }
                }
                newOrg.teams.push(newTeam)
              }
            }
            list.push(newOrg)
          }
        }
        //console.log(list)
        return list
      }
    )
  }

  //role=0 Member
  //role=1 Owner
  changeOrganizationMemberRole(orgName : string, userName : string, role : number) {
    return this.httpPut("/organizations/"+orgName+"/members/"+userName, {
      organization_name: orgName, user_name: userName, role: role
    })
  }

  login(username : string, pwd : string) {
    return this.httpPost("/login", {name: username, password: pwd});
  }

  signup(username : string, pwd : string, email : string) {
    return this.httpPost("/signup", {name: username, password: pwd, email: email, url: window.location.protocol +"//"+window.location.host});
  }

  verify(token : string) {
  return this.httpPost("/verify/"+token, { token: token});
}

  registration() {
    return this.httpGet("/clusters/registration");
  }

  removeUser(username : string) {
    return this.httpDelete("/users/"+username)
  }

  stacks() {
    return this.httpGet("/stacks")
    .map((res : Response) => {
      const data = res.json()
      let list : DockerStack[] = []
      if (data.entries) {
        for (let item of data.entries) {
          let stack = new DockerStack(
            item.stack.id,
            item.stack.name,
            item.services,
            item.stack.owner.name,
            item.stack.owner.type
          )
          stack.set(
            this.convertStatus(item.status),
            this.formatedDate(item.stack.create_dt),
            item.total_services,
            item.running_services
          )
          list.push(stack)
        }
      }
      return list
    })
  }

  deployStack(stackName : string, fileContent : string) {
    var data = [];
    for (let i = 0; i < fileContent.length; i++){
        data.push(fileContent.charCodeAt(i));
    }
    return this.httpPost("/stacks", { name : stackName, compose: data});
  }

  removeStack(stackName : string) {

    return this.httpDelete("/stacks/"+stackName)
  }

  services(stackName : string) {
    return this.httpGet("/services/"+stackName)
    .map((res : Response) => {
      const data = res.json()
      console.log(data)
      let list : DockerService[] = []
      if (data.entries) {
        for (let item of data.entries) {
          if (item.service && item.service.id) {
            let serv = new DockerService(
              this.shortcutId(item.service.id),
              item.service.name,
              item.service.mode,
              item.service.image,
              item.service.tag
            )
            serv.set(this.convertStatus(item.status), item.total_tasks, item.ready_tasks)
            list.push(serv)
          }
        }
      }
      console.log(list)
      return list
    })
  }

  tasks(serviceId : string) {
    return this.httpGet("/tasks/"+serviceId)
    .map((res : Response) => {
      const data = res.json()
      console.log(data)
      let list : DockerContainer[] = []
      if (data.tasks) {
        for (let item of data.tasks) {
          if (item.id) {
            let cont = new DockerContainer(
              this.shortcutId(item.id),
              item.image,
              this.convertStatus(item.current_state),
              this.convertStatus(item.desired_state),
              this.shortcutId(item.node_id)
            )
            list.push(cont)
          }
        }
      }
      return list
    })
  }

  serviceScale(serviceId : string, replicas : number) {
    return this.httpPut("/scale/"+serviceId+"/"+replicas,
    { service_id: serviceId, replicas_number: replicas})
  }

  organizationRessources() {
    return this.httpGet("/resources")
    .map((res : Response) => {
      let data = res.json()
      //console.log(data)
      let list : OrganizationResource[] = []
      if (data.resources) {
        for (let item of data.resources) {
          let type="unknow:"+item.type
          if (!item.type || item.type == 'RESOURCE_STACK') {
            type = "Stack"
          } else if (item.type == 'RESOURCE_DASHBOARD') {
            type = "Dashboard"
          }
          let res = new OrganizationResource(item.id, type, item.name)
          list.push(res)
        }
      }
      return list
    })
  }

  addResourceToTeam(orgName : string, teamName : string, resourceId : string) {
    return this.httpPost("/organizations/"+orgName+"/teams/"+teamName+"/resources",
      { organization_name: orgName, team_name: teamName, resource_id: resourceId}
    )
  }

  removeResourceFromTeam(orgName : string, teamName : string, resourceId : string) {
    return this.httpDelete("/organizations/"+orgName+"/teams/"+teamName+"/resources/"+resourceId)
  }

  changeTeamResourcePermissionLevel(orgName : string, teamName : string, resourceId : string, level : number) {
    return this.httpPut("/organizations/"+orgName+"/teams/"+teamName+"/resources/"+resourceId,
      { organization_name: orgName, team_name: teamName, resource_id: resourceId, permission_level: level }
    )
  }

  nodes() {
    return this.httpGet("/nodes")
      .map((res : Response) => {
        let data = res.json()
        let list : Node[] = []
        if (data.entries) {
          for (let item of data.entries) {
            let node = new Node(item.id)
            node.name = item.name;
            node.hostname = item.hostname;
            node.role = item.role;
            node.architecture = item.architecture
            node.os = item.os
            node.engine = item.engine
            node.status = item.status
            node.availability = item.availability
            node.leader = item.leader
            if (node.leader) {
              node.role = 'leader'
            }
            node.addr = item.addr
            node.reachability = item.reachability
            list.push(node);
          }
        }
        return list
      }
    )
  }

  logs(req : LogsRequest) {
    return this.httpPost("/logs", req)
      .map((res : Response) => {
        let data = res.json()
        let list : Log[] = []
        if (data.entries) {
          for (let item of data.entries) {
            let log = new Log(item.timestamp, item.msg)
            log.container_id = item.container_id
            log.container_name = item.container_mame
            log.container_short_name = item.container_short_name
            log.service_name = item.service_name
            log.service_id = item.service_id
            log.task_id = item.task_id
            log.stack_name = item.stack_name
            log.node_id = item.node_id
            list.push(log);
          }
        }
        return list
      }
    )
  }

  stats(req : StatsRequest) {
    if (req.time_group == "") {
      return this.statsCurrent(req);
    }
    return this.statsHistoric(req);
  }

  statsHistoric(request : StatsRequest) {
    return this.httpPost("/stats", request)
      .map((res : Response) => {
        let data = res.json()
        //console.log(data)
        let list : GraphHistoricData[] = []
        if (data.entries) {
          for (let item of data.entries) {
            let datal : { [name:string]: number; } = {}
            if (request.stats_cpu) {
              this.setValue(datal, 'cpu-usage', item.cpu.total_usage, 1, 1)
            }
            if (request.stats_io) {
              this.setValue(datal, 'io-total', item.io.total, 1, 1)
              this.setValue(datal, 'io-write', item.io.write, 1, 1)
              this.setValue(datal, 'io-read', item.io.read, 1, 1)
            }
            if (request.stats_mem) {
              this.setValue(datal, 'mem-limit', item.mem.limit, 1, 1)
              this.setValue(datal, 'mem-maxusage', item.mem.maxusage, 1, 1)
              if (request.format) {
                this.setValue(datal, 'mem-usage', item.mem.usage, 1, 1024*1024)
              } else {
                this.setValue(datal, 'mem-usage', item.mem.usage, 1, 1)
              }
              this.setValue(datal, 'mem-usage-p', item.mem.usage_p, 100, 1)
            }
            if (request.stats_net) {
              this.setValue(datal, 'net-rx-bytes', item.net.rx_bytes, 1, 1)
              this.setValue(datal, 'net-rx-packets', item.net.rx_packets, 1, 1)
              this.setValue(datal, 'net-tx-bytes', item.net.tx_bytes, 1, 1)
              this.setValue(datal, 'net-tx-packets', item.net.tx_packets, 1, 1)
              this.setValue(datal, 'net-total-bytes', item.net.total_bytes, 1, 1)
            }
            let hgraph = new GraphHistoricData(this.parseTime(item.group))
            hgraph.name =item.sgroup
            hgraph.values = datal
            hgraph.sdate = item.group
            list.push(hgraph)
          }
        }
        return list
      }
    );
  }

  statsCurrent(request : StatsRequest) {
    return this.httpPost("/stats", request)
      .map((res : Response) => {
        let data = res.json()
        //console.log(data)
        let list : GraphCurrentData[] = []
        if (data.entries) {
          for (let item of data.entries) {
            let datal : { [name:string]: number; } = {}
            if (request.stats_cpu) {
              this.setValue(datal, 'cpu-usage', item.cpu.total_usage, 1, 1)
            }
            if (request.stats_io) {
              this.setValue(datal, 'io-total', item.io.total, 1, 1)
              this.setValue(datal, 'io-write', item.io.write, 1, 1)
              this.setValue(datal, 'io-read', item.io.read, 1, 1)
            }
            if (request.stats_mem) {
              this.setValue(datal, 'mem-limit', item.mem.limit, 1, 1)
              this.setValue(datal, 'mem-maxusage', item.mem.maxusage, 1, 1)
              this.setValue(datal, 'mem-usage', item.mem.usage, 1, 1)
              this.setValue(datal, 'mem-usage-p', item.mem.usage_p, 100, 1)
            }
            if (request.stats_net) {
              this.setValue(datal, 'net-rx-bytes', item.net.rx_bytes, 1, 1)
              this.setValue(datal, 'net-rx-packets', item.net.rx_packets, 1, 1)
              this.setValue(datal, 'net-tx-bytes', item.net.tx_bytes, 1, 1)
              this.setValue(datal, 'net-tx-packets', item.net.tx_packets, 1, 1)
              this.setValue(datal, 'net-total-bytes', item.net.total_bytes, 1, 1)
            }
            list.push(new GraphCurrentData(item.group, datal))
          }
        }
        return list
      }
    );
  }

  setValue(datal :{ [name:string]: number; }, name : string, val : number, mul : number, div : number) {
    datal[name] = this.getValue(val, mul, div)
  }

  getValue(val: number, mul: number, div: number) {
    if (val) {
      return (val*mul)/div
    }
    return 0
  }

  createDashboard(name : string, data : string) {
    return this.httpPost("/dashboards", { name: name, data: data})
      .map((res : Response) => {
        let data = res.json()
        if (data.dashboard) {
          return data.dashboard.id
        }
        return undefined
      }
    )
  }

  getDashboard(id : string) {
    return this.httpGet("/dashboards/"+id)
      .map((res : Response) => {
        let data = res.json()
        if (data.dashboard) {
          let dashboard = new Dashboard(data.dashboard.id, data.dashboard.name, data.dashboard.data)
          let sdate = this.formatedDate(data.dashboard.create_dt)
          dashboard.set(data.dashboard.owner.name, data.dashboard.owner.type, sdate)
          return dashboard
        }
        return undefined
      }
    )
  }

  listDashboard() {
    return this.httpGet("/dashboards")
      .map((res : Response) => {
        let data = res.json()
        let list : Dashboard[] = []
        if (data.dashboards) {
            for (let dash of data.dashboards) {
              let dashboard = new Dashboard(dash.id, dash.name, dash.data)
              let sdate = this.formatedDate(dash.create_dt)
              dashboard.set(dash.owner.name, dash.owner.type, sdate)
              list.push(dashboard)
            }
        }
        list.sort((a ,b) => {
          if (a.date < b.date) {
            return 1
          } else {
            return -1
          }
        })
        return list
      }
    )
  }

  updateDashboardName(id : string, name : string) {
    return this.httpPut("/dashboards/"+id+"/name/"+name, {});
  }

  updateDashboard(id : string, data : string) {
    return this.httpPut("/dashboards/"+id+"/data", {id: id, data: data});
  }

  removeDashboard(id : string) {
    return this.httpDelete("/dashboards/"+id)
  }


//--------------------------------------------------------------------------------------
// http core functions
//--------------------------------------------------------------------------------------

  private formatedDate(daten : number) : string {
    let date = new Date(daten * 1000)
    let num = ""+date.getDate()
    if (date.getDate()<10) {
      num='0'+num
    }
    let month = ""+(date.getMonth()+1)
    if (date.getMonth()+1<10) {
      month = '0'+month
    }
    return date.getFullYear()  + "-" +
    month + "-" +
    num + " " +
    date.getHours() + ":" +
    date.getMinutes();
  }

  convertStatus(status: string) : string {
    if (status) {
      return status.toLowerCase()
    }
    return "unknow"
  }

  shortcutId(id : string) {
    if (!id) {
      return "unknow"
    }
    return id.substring(0, 12)
  }

  private setHeaders() {
    var headers = new Headers
    headers.set('Authorization', 'amp '+this.token)
    return headers
  }

  setToken(token : string) {
    this.token = token
  }

  httpGet(url : string) : Observable<any> {
    let headers = this.setHeaders()
    return this.http.get(this.addr+url, { headers: this.setHeaders() })
      .retryWhen(e => e.scan<number>((errorCount, err) => {
        console.log("retry: "+(errorCount+1))
        if (errorCount >= httpRetryNumber-1) {
            throw err;
        }
        return errorCount + 1;
      }, 0).delay(httpRetryDelay)
    )
  }

  httpDelete(url : string) : Observable<any> {
    let headers = this.setHeaders()
    return this.http.delete(this.addr+url, { headers: this.setHeaders() })
      .retryWhen(e => e.scan<number>((errorCount, err) => {
        console.log("retry: "+(errorCount+1))
        if (errorCount >= httpRetryNumber-1) {
            throw err;
        }
        return errorCount + 1;
      }, 0).delay(httpRetryDelay)
    )
  }

  httpPost(url : string, data : any) : Observable<any> {
    let headers = this.setHeaders()
    return this.http.post(this.addr+url, data, { headers: this.setHeaders() })
      .retryWhen(e => e.scan<number>((errorCount, err) => {
        console.log("retry: "+(errorCount+1))
        if (errorCount >= httpRetryNumber-1) {
            throw err;
        }
        return errorCount + 1;
      }, 0).delay(httpRetryDelay)
    )
  }

  httpPut(url : string, data : any) : Observable<any> {
    let headers = this.setHeaders()
    return this.http.put(this.addr+url, data, { headers: this.setHeaders() })
      .retryWhen(e => e.scan<number>((errorCount, err) => {
        console.log("retry: "+(errorCount+1))
        if (errorCount >= httpRetryNumber-1) {
            throw err;
        }
        return errorCount + 1;
      }, 0).delay(httpRetryDelay)
    )
  }

}
