
export class Graph {
  public id : number
  public x : number
  public y : number
  public width: number
  public height: number
  public type: string
  public fields: string[]
  public title: string
  public yTitle: string
  public requestId: string
  public modeParameter: boolean
  public object : string
  public field: string
  public border: boolean
  public topNumber: number
  public alert : boolean
  public alertMin: number
  public alertMax: number
  public criterion: string//not used
  public criterionValue: string

  constructor(id: number, x : number, y : number, w: number, h: number, type: string, fields : string[], title : string, yTitle : string) {
    this.id = id
    this.x = x
    this.y = y
    this.width = w
    this.height = h
    this.type = type
    this.fields = fields
    this.title = title
    this.yTitle = yTitle
    this.modeParameter = false
    this.border = false
    this.topNumber = 3
    this.alert = false
    this.alertMin = 0
    this.alertMax = 1000
    this.criterion = ""
    this.criterionValue = ""
  }

}
