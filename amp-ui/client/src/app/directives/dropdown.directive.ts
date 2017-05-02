import { Directive, ElementRef, Output, HostListener, HostBinding } from '@angular/core';

@Directive({
  selector: '[appDropdown]'
})
export class DropdownDirective {
  @HostBinding('class.open') isOpen = false;

  @HostListener('click') toggleOpen() {
    this.isOpen = !this.isOpen;
  }

  @HostListener('document.click') close() {
    this.isOpen = false;
  }

  @Output()

  @HostListener('document:click', ['$event', '$event.target'])
  public onClick(event: MouseEvent, targetElement: HTMLElement): void {
      if (!targetElement) {
          return;
      }

      const clickedInside = this._elementRef.nativeElement.contains(targetElement);
      if (!clickedInside) {
          this.isOpen = false
      }
  }
  constructor(private _elementRef: ElementRef) {}
}
