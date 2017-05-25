

export class GraphColors {
  public object: string
  private index: number
  private names: string[]
  private colorMap: { [name:string]: string; }


  constructor(object : string) {
    this.object = object
    this.index = 0
    this.colorMap = {}
    this.names = []
  }

  getColor(name : string) : string {
    return this.colorMap[name]
  }

  setColor(name : string, color : string) {
    let col = this.colorMap[name]
    if (!col) {
      this.names[this.index] = name
      //console.log(this.object+": "+name+" color="+color)
      this.index++
    }
    this.colorMap[name] = color
  }

  getIndex() : number {
    return this.index;
  }

  getName(i : number) : string {
    return this.names[i]
  }

  clear() {
    this.index = 0
    this.colorMap = {}
    this.names = []
  }

}
