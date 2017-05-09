export class DockerContainer {
  public id: string
  public name: string
  public image: string
  public status: string
  public command: string

  constructor(id : string, name : string) {
    this.id = id
    this.name = name
  }
}
