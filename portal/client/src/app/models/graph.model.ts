
export class Graph {
  public id : string
  public x : number
  public y : number
  public width: number
  public height: number
  public reelX: number
  public reelY: number
  public reelWidth: number
  public reelHeight: number
  public type: string
  public fields: string[]
  public containerAvg: boolean
  public title: string
  public centerTitle: boolean
  public yTitle: string
  public requestId: string
  public modeParameter: boolean
  public object : string
  public field: string
  public border: boolean
  public counterHorizontal: boolean
  public topNumber: number
  public alert : boolean
  public alertMin: number
  public alertMax: number
  public criterion: string
  public criterionValue: string
  public histoPeriod: string //for historic request
  public bubbleXField: string
  public bubbleYField: string
  public bubbleSizeField: string
  public bubbleScale: string
  public stackedAreas: boolean
  public percentAreas: boolean
  public legendNames: string[]
  public legendColors: string[]
  public transparentLegend: boolean
  public removeLocalLegend: boolean

  constructor(id: string, x : number, y : number, w: number, h: number, type: string, title : string) {
    this.id = id
    this.x = x
    this.y = y
    this.width = w
    this.height = h

    this.type = type
    this.fields=[]
    this.title = title
    this.border = true
    this.modeParameter = false
    this.topNumber = 3
    this.alert = false
    this.alertMin = 0
    this.alertMax = 1000
    this.criterion = ""
    this.criterionValue = ""
    this.stackedAreas = true
    this.legendNames = []
    this.legendColors = []
  }

}
