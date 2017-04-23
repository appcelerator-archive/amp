export class Endpoint {
  public host: string;
  public local: boolean;


  constructor(host : string, local : boolean) {
    this.host = host
    this.local = local
  }
}
