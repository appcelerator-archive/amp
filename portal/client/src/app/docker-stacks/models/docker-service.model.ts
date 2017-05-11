export class DockerService {
  public id: string
  public name: string
  public image: string
  public mode: string
  public replicas: string
  public ports: string

  constructor(id : string, name : string, mode : string, replicas : string, image : string) {
    this.id = id
    this.name = name
    this.mode = mode
    this.replicas = replicas
    this.image = image
  }
}
