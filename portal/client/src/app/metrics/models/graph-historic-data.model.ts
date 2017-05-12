
export class GraphHistoricData {
  public date: Date
  public name: string
  public values: { [name:string]: number; }
  public graphValues: number[]

  constructor(date : Date, name: string, values : { [name:string]: number; }) {
    this.date = date
    this.name = name
    this.values = values
  }
}
