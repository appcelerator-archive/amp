export class DockerStack {
  public id: string;
  public shortId: string;
  public name: string;
  public createDate: string
  public services: string;
  public ownerName: string;
  public ownerType: string;
  public status: string;
  public runningService: number;
  public totalService: number;

  constructor(id: string, name: string, services: number, ownerName : string, ownerType : string) {
    this.id = id,
    this.shortId = id.substring(0, 12);
    this.name = name;
    this.ownerName = ownerName;
    this.ownerType = ownerType
  }

  set(status : string, date : string, totalService : number, runningService : number) {
    this.status = status
    this.createDate = date
    this.totalService = totalService
    this.runningService = runningService
    this.services = ""+this.runningService +"/"+this.totalService
  }

}
