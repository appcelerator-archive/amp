export class ItemMenu {
  public name: string
  public description: string
  public route: string

  constructor(name : string, description : string, route : string) {
    this.name = name
    this.description = description
    this.route = route
  }
}
