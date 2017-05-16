export class LogsRequest {
  public container: string
  public message: string
  public node: string
  public size: number
  public service: string
  public stack: string
  public task: string
  public infra: boolean

  constructor() {
  }
}
