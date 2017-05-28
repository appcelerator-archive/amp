
export class GraphHistoricData {
  public date: Date
  public sdate: string
  public name: string
  public max: number[]
  public values: { [name:string]: number; }
  public graphValues: number[]
  public graphValuesUnit: number[]

  constructor(date : Date) {
    this.date = date
    this.graphValues = []
    this.graphValuesUnit = []
  }
}
