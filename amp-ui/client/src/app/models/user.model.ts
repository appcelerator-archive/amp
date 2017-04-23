export class User {
  public name: string;
  public email: string;
  public role: string;
  public verified: boolean;
  public checked: boolean;

  constructor(name: string, email: string, role: string) {
    this.name = name;
    this.email = email;
    this.role = role;
    this.verified = false;
    this.checked = false;
  }

}
