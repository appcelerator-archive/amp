import { OrganizationResource } from './organization-resource.model'
import { Member } from './member.model'

export class Team {
  public name: string
  public members: Member[]
  public resources: OrganizationResource[]

  constructor(name : string) {
    this.name = name
    this.members = []
    this.resources = []
  }
}
