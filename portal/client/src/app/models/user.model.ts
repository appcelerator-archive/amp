import { Member } from '../organizations/models/member.model'

export class User {
  public name: string;
  public email: string;
  public member: Member;
  public verified: boolean;
  public createDate: string;
  public checked: boolean;
  public pendingOrganizations: string[];
  public label: string
  public tokenUsed: boolean;

  constructor(name: string, email: string, role: string) {
    this.name = name;
    this.email = email;
    this.member = undefined
    this.verified = false;
    this.checked = false;
    this.pendingOrganizations = [];
    this.label=name
    if (role == "owner") {
      this.label+=" (owner)"
    }
  }

}
