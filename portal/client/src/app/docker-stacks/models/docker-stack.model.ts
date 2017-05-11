export class DockerStack {
  public id: string;
  public shortId: string;
  public name: string;
  public services: number;
  public ownerName: string;
  public ownerType: string;

  constructor(id: string, name: string, services: number, ownerName : string, ownerType : string) {
    this.id = id,
    this.shortId = id.substring(0, 12);
    this.name = name;
    this.services = services;
    this.ownerName = ownerName;
    this.ownerType = ownerType
  }

}
