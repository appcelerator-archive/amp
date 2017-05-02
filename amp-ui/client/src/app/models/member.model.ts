export class Member {
  userName: string
  role: number
  status: number //pure UI: use to detect with one is removed or added 0: not moved, 1:added, -1: removed
  saved: boolean
  saveError: string

  constructor(name: string, role: number) {
    this.userName = name
    if (role === undefined) { //shit grpc
      this.role=0
    } else {
      this.role = role
    }
    this.status = 0
    this.saved = false
    this.saveError = ""
  }

  public getLabeledName() : string {
    if (this.role == 0) {
      return this.userName
    } else {
      return this.userName+" (owner)"
    }
  }
}
