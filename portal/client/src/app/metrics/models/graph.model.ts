
export class Graph {
  public x : number
  public y : number
  public width: number
  public height: number
  public fields: string[]
  public title: string
  public yTitle: string

  constructor(x : number, y : number, w: number, h: number, fields : string[], title : string, yTitle : string) {
    this.x = x
    this.y = y
    this.width = w
    this.height = h
    this.fields = fields
    this.title = title
    this.yTitle = yTitle
  }

}
