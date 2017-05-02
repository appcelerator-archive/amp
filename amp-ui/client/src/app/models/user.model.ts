export class User {
  public name: string;
  public email: string;
  public role: string;
  public verified: boolean;
  public checked: boolean;
  public pendingOrganizations: string[];
  public label: string

  constructor(name: string, email: string, role: string) {
    this.name = name;
    this.email = email;
    this.role = role;
    this.verified = false;
    this.checked = false;
    this.pendingOrganizations = [];
    this.label=name
    if (role == "owner") {
      this.label+=" (owner)"
    }
  }

}
