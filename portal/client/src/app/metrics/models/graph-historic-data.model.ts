
export class GraphHistoricData {
  public date: Date
  public values: { [name:string]: number; }
  public graphValues: number[]

  constructor(date : Date, values : { [name:string]: number; }) {
    this.date = date
    this.values = values
  }
}
