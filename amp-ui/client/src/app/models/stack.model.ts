export class Stack {
  public id: string;
  public name: string;
  public services: number;
  public ownerName: string;
  public ownerType: string;

  constructor(id: string, name: string, services: number, ownerName : string, ownerType : string) {
    this.id = id,
    this.name = name;
    this.services = services;
    this.ownerName = ownerName;
    this.ownerType = ownerType
  }

  match(stack : Stack, value : string) : boolean {
    if (stack.id.includes(value)) {
      return true
    }
    if (stack.name.includes(value)) {
      return true
    }
    if (stack.services.toString().includes(value)) {
      return true
    }
    if (stack.ownerName.includes(value)) {
      return true
    }
    if (stack.ownerType.includes(value)) {
      return true
    }
    return false
  }
}
