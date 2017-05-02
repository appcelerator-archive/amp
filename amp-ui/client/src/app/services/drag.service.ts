import { Injectable } from '@angular/core';

@Injectable()
export class DragService {
  private zone: string;
  
  startDrag(zone: string) {
    this.zone = zone;
  }
  
  accepts(zone: string): boolean {
    return zone == this.zone;
  }
}