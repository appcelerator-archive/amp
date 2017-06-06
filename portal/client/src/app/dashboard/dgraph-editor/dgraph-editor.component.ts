
import { Component, HostListener, OnInit, OnDestroy, Input, ElementRef, ViewChild} from '@angular/core';
import { Graph } from '../../models/graph.model';
import { DashboardService } from '../services/dashboard.service'
import { MenuService } from '../../services/menu.service';

@Component({
  selector: 'app-dgraph-editor',
  templateUrl: "./dgraph-editor.component.html",
  styleUrls: ['./dgraph-editor.component.css'],
})

export class DGraphEditorComponent implements OnInit, OnDestroy {
  @Input() public graph : Graph;
  all="[centertitle][request][setting][top][histoperiod][bubble][areas][alert]"
  visibility : { [name:string]: string; } = {}
  messageLegend = ""

  constructor(
    public dashboardService : DashboardService) {
    this.visibility['text']=""
    this.visibility['lines']="[request][centertitle][setting][top][histoperiod]"
    this.visibility['areas']="[request][centertitle][setting][top][histoperiod][areas]"
    this.visibility['bars']="[request][centertitle][setting][top][removelegend]"
    this.visibility['pie']="[request][setting][top]"
    this.visibility['bubbles']="[request][centertitle][setting][top][bubble]"
    this.visibility['counterSquare']="[horizontal][request][alert]"
    this.visibility['counterCircle']="[horizontal][request][alert]"
    this.visibility['legend']="[centertitle][setting][legend]"
    this.visibility['innerStats']="[centertitle][setting][legend]"
  }

  ngOnInit() {
  }

  ngOnDestroy() {
  }

  isVisible(name : string) : boolean {
    let visibility = this.visibility[this.dashboardService.selected.type]
    if (!visibility) {
      return false
    }
    if (visibility.indexOf("["+name+"]")>=0) {
      return true
    }
    return false
  }

}
