import { OrganizationMember } from './organization-member.model'
import { Team } from './team.model'

export class Organization {
  public name: string;
  public email: string;
  public members: OrganizationMember[]
  public teams: Team[]


  constructor(name: string, email: string) {
    this.name = name;
    this.email = email;
  }

}
