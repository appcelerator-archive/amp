
import { Team } from './team.model'
import { Member } from './member.model'
import { TeamResource } from './team-resource.model'

export class Organization {
  public name: string;
  public email: string;
  public members: Member[]
  public teams: Team[]
  public resources: TeamResource[]


  constructor(name: string, email: string) {
    this.name = name;
    this.email = email;
    this.members = []
    this.resources = []
    this.teams = []
  }

}
