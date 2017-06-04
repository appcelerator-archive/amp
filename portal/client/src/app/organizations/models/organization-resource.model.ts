const labels : string[] = ['read', 'write', 'admin']

export class OrganizationResource {
  public id: string
  public type: string
  public name: string
  public permissionLevel: number
  public permissionLabel: string
  //UI
  public status: number //pure UI: use to detect with one is removed or added 0: not moved, 1:added, -1: removed
  public saved: boolean
  public saveError: string
  public changeAuth: boolean
  public changeAuthError: string

  constructor(id : string, type : string, name) {
    this.id = id
    this.type = type
    this.name = name
    this.setAuthorization(0)
    this.status = 0
    this.saved = false
    this.saveError = ""
  }

  public setAuthorization(perm : number) {
    this.permissionLevel = perm
    if (perm<0 || perm>=labels.length) {
      this.permissionLevel = 0
    }
    this.permissionLabel = labels[this.permissionLevel]
  }

  public getLabeledName() : string {
    return this.type+":"+this.name
  }
}
