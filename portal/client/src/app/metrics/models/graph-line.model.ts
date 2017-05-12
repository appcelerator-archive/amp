export class GraphLine {
  public name: string
  public displayedName: string
  public color: string

  constructor(name : string, color : string) {
    this.name = name
    let list = name.split('_')
    if (list.length>1) {
      this.displayedName = list[1]
      for (let ii=2;ii<list.length;ii++) {
        this.displayedName+='_'+list[ii]
      }
    } else {
      this.displayedName = name
    }
    this.color = color
  }
}
