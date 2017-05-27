
import { Team } from './team.model'
import { Member } from './member.model'
import { OrganizationResource } from './organization-resource.model'

export class Organization {
  public name: string;
  public email: string;
  public members: Member[]
  public teams: Team[]
  public resources: OrganizationResource[]


  constructor(name: string, email: string) {
    this.name = name;
    this.email = email;
    this.members = []
    this.resources = []
    this.teams = []
  }

}
