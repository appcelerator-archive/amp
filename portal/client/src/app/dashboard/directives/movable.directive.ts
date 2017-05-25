import { Directive, EventEmitter, ElementRef, OnInit, Input } from '@angular/core';
import { Graph } from '../../models/graph.model'

@Directive({
  selector: '[graphMovable]',
  host: {
		'(mousedown)': 'onMouseDown($event)',
		'(document: mousemove)': 'onMouseMove($event)',
		'(document: mouseup)': 'onMouseUp($event)',
    /*
		'(keydown.ArrowUp)': 'onNudge($event)',
		'(keydown.ArrowRight)': 'onNudge($event)',
		'(keydown.ArrowDown)': 'onNudge($event)',
		'(keydown.ArrowLeft)': 'onNudge($event)',
		'(keyup.ArrowUp)': 'onNudge($event)',
		'(keyup.ArrowRight)': 'onNudge($event)',
		'(keyup.ArrowDown)': 'onNudge($event)',
		'(keyup.ArrowLeft)': 'onNudge($event)'
    */
	}
})

export class MovableDirective implements OnInit {
  @Input() graph: Graph;
  private keys: Array<number> = [37, 38, 39, 40];
  private clienty0 = 0
  private clientx0 = 0
  private graphx0 = 0
  private graphy0 = 0
  private movable = false
  private corner = 0

  constructor(private eRef : ElementRef) {
  }

  ngOnInit() {
    this.eRef.nativeElement.style.position = "absolute"
    this.eRef.nativeElement.style.left = "0px"
    this.eRef.nativeElement.style.top = "0px"
  }

  onMouseDown($event) {
    this.movable = true
    this.graphx0 = parseInt(this.eRef.nativeElement.style.left.replace('px', ''));
    this.graphy0 = parseInt(this.eRef.nativeElement.style.top.replace('px', ''));
    this.clientx0 = $event.clientX
    this.clienty0 = $event.clientY
  }

  onMouseUp($event) {
    this.movable = false
  }

  onMouseMove($event) {
    if (this.movable) {
      if (!this.graph) {
        this.eRef.nativeElement.style.left = (($event.clientX - this.clientx0) + this.graphx0) + 'px';
        this.eRef.nativeElement.style.top = (($event.clientY - this.clienty0) + this.graphy0) + 'px';
      } else {
        this.graph.x = (($event.clientX - this.clientx0) + this.graphx0);
        this.graph.y = (($event.clientY - this.clienty0) + this.graphy0);
      }
    }
  }

  private onNudge($event) {
		this.keys[$event.keyCode] = $event.type === 'keydown' ? 1 : 0;

    let x = parseInt(this.eRef.nativeElement.style.left.replace('px', ''));
    let y = parseInt(this.eRef.nativeElement.style.top.replace('px', ''));

    x = x - this.keys[37] + this.keys[39];
		y = y - this.keys[38] + this.keys[40];

    this.eRef.nativeElement.style.left = x + 'px';
    this.eRef.nativeElement.style.top = y + 'px';

	}


}
