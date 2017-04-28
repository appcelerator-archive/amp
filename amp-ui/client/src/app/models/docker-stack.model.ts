export class DockerStack {
  public id: string;
  public name: string;
  public services: number;
  public ownerName: string;
  public ownerType: string;

  constructor(id: string, name: string, services: number, ownerName : string, ownerType : string) {
    this.id = id,
    this.name = name;
    this.services = services;
    this.ownerName = ownerName;
    this.ownerType = ownerType
  }

}
