export class DockerService {
  public id: string
  public name: string
  public image: string
  public mode: string
  public replicas: string
  public ports: string

  constructor(id : string, name : string) {
    this.id = id
    this.name = name
  }
}
