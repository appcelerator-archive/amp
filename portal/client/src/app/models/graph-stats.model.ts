
export class GraphStats {
  public group : string
  public values: { [name:string]: number; }

  constructor(group : string, values : { [name:string]: number; }) {
    this.group = group
    this.values = values
  }
}
