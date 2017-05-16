export class Log {
  public timestamp: string
  public container_id: string
  public container_name: string
  public container_short_name: string
  public container_state: string
  public service_name: string
  public service_id: string
  public task_id: string
  public stack_name: string
  public node_id: string
  public msg: string

  constructor(timestamp : string, msg : string) {
    this.timestamp = this.dateFormat(timestamp)
    this.msg = msg
  }

  private dateFormat(date : string) : string {
    return date.replace('T', ' ').replace('Z', '')
  }
}
