
export class DashboardInnerStats {

  public refreshRequestNb: number
  public totalRequestTime: number
  public totalRefreshTime: number
  public avgRefreshTime: number
  public maxRequestTime: number
  public maxGraphTitle: string
  public requestErrorNb: number
  public refreshNb: number

  constructor() {
    this.refreshNb=0
    this.totalRefreshTime=0
  }

  public initNewRefresh() {
    this.refreshRequestNb=0
    this.totalRequestTime=0
    this.maxRequestTime=0
    this.requestErrorNb=0
    this.maxGraphTitle=""
    this.refreshNb++
  }

  public setRequestTime(t0: number, title: string) {
    let requestTime = ( new Date().getTime()) - t0
    this.refreshRequestNb++
    this.totalRequestTime += requestTime
    if (requestTime > this.maxRequestTime) {
      this.maxRequestTime = requestTime
      this.maxGraphTitle = title
    }
    this.totalRefreshTime+= requestTime
    this.avgRefreshTime = this.totalRefreshTime / this.refreshNb
  }

  public setRequestError() {
    this.requestErrorNb++
  }
}
