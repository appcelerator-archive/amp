import { TeamResource } from './team-resource.model'
import { User } from './user.model'

export class Team {
  public name: string
  public members: User[]
  public resources: TeamResource[]

  constructor(name : string) {
    this.name = name
    this.members = []
    this.resources = []
  }
}
