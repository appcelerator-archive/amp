

export class Dashboard {
  public id : string
  public name : string
  public ownerName : string
  public ownerType : string
  public data : string
  public date : Date

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

  set(ownerName : string, ownerType : any, date : number) {
    this.ownerName = ownerName
    this.setOwnerType(ownerType)
    this.date = new Date(date * 1000)
  }

  formatedDate() : string {
    return  this.date.getDate()  + "-" +
    (this.date.getMonth()+1) + "-" +
    this.date.getFullYear() + " " +
    this.date.getHours() + ":" +
    this.date.getMinutes();
  }

}
