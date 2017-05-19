const labels : string[] = ['none', 'read', 'write', 'admin']

export class TeamResource {
  public id: string
  public type: string
  public name: string
  public permissionLevel: number
  public permissionLabel: string

  constructor(id : string, type : string, name : string, permissionLevel : number) {
    this.id = id
    this.type = type
    this.name = name
    this.setPermission(permissionLevel)
  }

  setPermission(perm : number) {
      this.permissionLevel = perm
    if (perm<0 || perm>=labels.length) {
      this.permissionLevel = 0
    }
    this.permissionLabel = labels[this.permissionLevel]
  }
}
