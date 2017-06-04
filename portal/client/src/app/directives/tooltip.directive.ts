import { Directive, EventEmitter, ElementRef, Renderer2, OnInit, Input } from '@angular/core';

@Directive({
  selector: '[tooltip]',
  host: {
		'(mouseenter)': 'onMouseEnter($event)',
		'(mouseleave)': 'onMouseLeave()',
	}
})

export class TooltipDirective implements OnInit {
  tooltip : any

  constructor(
    private eRef : ElementRef,
    private renderer: Renderer2) {
  }

  ngOnInit() {
    let div = this.renderer.createElement('div')
    this.renderer.setAttribute(div, 'id', 'ampTooltip');
    this.renderer.setAttribute(div, 'id', 'ampTooltip');
    this.renderer.setStyle(div, 'height', '20px');
    this.renderer.setStyle(div, 'width', '100px');
    this.renderer.setStyle(div, 'background-color', 'white');
    this.renderer.setStyle(div, 'border', '1px solid black');
    this.renderer.setProperty(div, 'innerHTML', 'test');
    this.renderer.setStyle(div, 'position', 'absolute');
    this.renderer.setProperty(div, 'hidden', true);
    this.renderer.appendChild(this.eRef.nativeElement, div)
    this.tooltip = div
  }

  onMouseEnter($event) {
    //console.log("enter")
    //let x = parseInt(this.eRef.nativeElement.style.left.replace('px', ''));
    //let y = parseInt(this.eRef.nativeElement.style.top.replace('px', ''));
    let x = $event.clientX
    let y = $event.clientY
    let div = document.getElementById('ampTooltip')
    this.renderer.setStyle(div, 'left', x+"px")
    this.renderer.setStyle(div, 'top', y+"px")
    //this.renderer.setProperty(div, 'hidden', false);
  }

  onMouseLeave() {
    //console.log("leave")
    let div = document.getElementById('ampTooltip')
    this.renderer.setProperty(div, 'hidden', true);
  }

  ngOnDestroy(): void {
    // hide tooltip
  }

}
