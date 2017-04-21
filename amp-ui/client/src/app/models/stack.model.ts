export class Stack {
  public id: string;
  public name: string;
  public services: number;
  public ownerName: string;

  constructor(id: string, name: string, services: number, ownerName : string) {
    this.id = id,
    this.name = name;
    this.services = services;
    this.ownerName = ownerName;
  }
}
