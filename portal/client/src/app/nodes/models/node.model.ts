export class Node {
  public id: string
  public shortId: string
  public name: string
  public role: string
  public hostname: string
  public architecture: string
  public os: string
  public engine: string
  public status: string
  public availability: string
  public leader: boolean
  public addr: string
  public reachability: string

  constructor(id: string) {
    this.id = id
    this.shortId = id.substring(0, 12)
  }
}
