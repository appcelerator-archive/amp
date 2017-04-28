import { TeamResource } from './team-resource.model'

export class Team {
  public name: string
  public members: string[]
  public resources: TeamResource[]

  constructor(name : string) {
    this.name = name
  }
}
