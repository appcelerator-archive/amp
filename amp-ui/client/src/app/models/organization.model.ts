
import { Team } from './team.model'
import { User } from './user.model'
import { TeamResource } from './team-resource.model'

export class Organization {
  public name: string;
  public email: string;
  public members: User[]
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
