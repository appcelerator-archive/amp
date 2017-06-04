
export class GraphCurrentData {
  public group : string
  public values: { [name:string]: number; }
  public valueUnit: number
  public valueUnitx: number
  public valueUnity: number

  constructor(group : string, values : { [name:string]: number; }) {
    this.group = group
    this.values = values
    this.valueUnit = 0
    this.valueUnitx = 0
    this.valueUnity = 0
  }
}
