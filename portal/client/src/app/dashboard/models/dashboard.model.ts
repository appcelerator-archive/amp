

export class Dashboard {
  public id : string
  public name : string
  public ownerName : string
  public ownerType : string
  public data : string
  public date : string

  constructor(id: string, name: string, data: string) {
    this.id = id
    this.name = name
    this.data = data
  }

  setOwnerType(type : any) {
    if (!type || type == 0) {
      this.ownerType = "USER"
    } else if (type == 1) {
      this.ownerType = "ORGANIZATION"
    } else {
      this.ownerType = type
    }
  }

  set(ownerName : string, ownerType : any, sdate : string) {
    this.ownerName = ownerName
    this.setOwnerType(ownerType)
    this.date = sdate
  }

}
