export class User {
  public name: string;
  public email: string;
  public role: string;
  public password: string;
  public verified: boolean;
  public checked: boolean;

  constructor(name: string, email: string, pwd: string, role: string) {
    this.name = name;
    this.email = email;
    this.password = pwd;
    this.role = role;
    this.verified = false;
    this.checked = false;
  }
}
