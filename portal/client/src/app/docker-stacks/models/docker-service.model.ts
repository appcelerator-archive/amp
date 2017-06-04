export class DockerService {
  public id: string
  public name: string
  public image: string
  public mode: string
  public replicas: string
  public tag: string
  public status: string
  public readyTasks : number
  public totalTasks : number

  constructor(id : string, name : string, mode : string, replicas : string, image : string) {
    this.id = id
    this.name = name
    this.mode = mode
    this.replicas = replicas
    this.image = image
  }

  set(status : string, totalTasks : number, readyTasks : number) {
    this.status = status
    this.totalTasks = totalTasks
    this.readyTasks = readyTasks
    this.replicas = ""+this.readyTasks + "/" + this.totalTasks
  }
}
