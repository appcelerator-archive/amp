import { Directive, EventEmitter, ElementRef, Renderer2, OnInit, Input } from '@angular/core';
import { MenuService } from '../services/menu.service'

@Directive({
  selector: '[tooltip]',
  host: {
		'(mouseenter)': 'onMouseEnter($event)',
		'(mouseleave)': 'onMouseLeave()'
	}
})

export class TooltipDirective implements OnInit {
  @Input() tooltipLabel: string;
  constructor(
    private menuService : MenuService,
    private eRef : ElementRef,
    private renderer: Renderer2) {
  }

  ngOnInit() {
  }

  onMouseEnter($event) {
    this.menuService.tooltipLabel = this.tooltipLabel
  }

  onMouseLeave() {
    this.menuService.tooltipLabel = ""
  }

  ngOnDestroy(): void {
    this.menuService.tooltipLabel = ""
  }

}
