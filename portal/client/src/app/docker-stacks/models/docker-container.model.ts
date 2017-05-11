export class DockerContainer {
  public id: string
  public shortId: string
  public image: string
  public state: string
  public desiredState: string
  public nodeId: string

  constructor(id : string, image : string, state : string, desiredState : string, nodeId : string ) {
    this.id = id
    this.shortId = id.substring(0, 12)
    this.image = image
    this.state = state
    this.desiredState = desiredState
    this.nodeId = nodeId
  }
}
