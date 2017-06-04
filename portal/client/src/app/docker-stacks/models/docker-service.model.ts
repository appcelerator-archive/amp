export class DockerService {
  public id: string
  public name: string
  public shortName: string
  public image: string
  public mode: string
  public replicas: string
  public tag: string
  public status: string
  public readyTasks : number
  public totalTasks : number

  constructor(id : string, name : string, mode : string, image : string, tag : string) {
    this.id = id
    this.name = name
    this.shortName = this.extractShortName(name)
    this.mode = mode
    this.image = image
    this.tag = tag
  }

  set(status : string, totalTasks : number, readyTasks : number) {
    this.status = status
    this.totalTasks = totalTasks
    this.readyTasks = readyTasks
    this.replicas = ""+this.readyTasks + "/" + this.totalTasks
  }

  extractShortName(fullName : string) {
    if (!fullName) {
      return "unknow"
    }
    let ll = fullName.indexOf('_')
    if (ll<0) {
      return fullName
    }
    return fullName.substring(ll+1)
  }
}
