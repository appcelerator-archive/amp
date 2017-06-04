export class Member {
  public userName: string
  public role: string
  //UI
  public status: number //pure UI: use to detect with one is removed or added 0: not moved, 1:added, -1: removed
  public saved: boolean
  public saveError: string

  constructor(name: string, role: string) {
    this.userName = name
    if (role === undefined) { //shit grpc
      this.role="Member"
    } else if (role == "ORGANIZATION_MEMBER") {
      this.role = "Member"
    } else if (role == "ORGANIZATION_OWNER") {
      this.role = "Owner"
    } else {
      this.role = "."
    }
    this.status = 0
    this.saved = false
    this.saveError = ""
  }

  public getLabeledName() : string {
    if (this.role != "Owner") {
      return this.userName
    } else {
      return this.userName+" (owner)"
    }
  }
}
