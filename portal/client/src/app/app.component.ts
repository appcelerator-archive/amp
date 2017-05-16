import { Component, OnInit, Renderer } from '@angular/core';
import { MenuService } from './services/menu.service';
import { Router } from '@angular/router';
import 'rxjs/add/operator/pairwise';


@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})

export class AppComponent {

  constructor(
    public menuService : MenuService,
    private renderer: Renderer,
    private router : Router) { }

  ngOnInit() {
    this.menuService.waitingCursor(false)
    let that = this
    this.renderer.listen('window', 'resize', (evt) => {
      this.menuService.resize(evt)
    })
    this.router.events
      .subscribe((e: any) => {
        this.menuService.pushPath(e.urlAfterRedirects);
      }
    );
  }
}
