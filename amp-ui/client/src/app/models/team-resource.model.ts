export class TeamResource {
  public id: string
  public permissionLevel: number

  constructor(id : string, permissionLevel : number) {
    this.id = id
    this.permissionLevel = permissionLevel
  }
}
