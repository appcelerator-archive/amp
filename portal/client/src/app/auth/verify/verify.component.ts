import { Component, OnInit, OnDestroy } from '@angular/core';
import { MenuService } from '../../services/menu.service';
import { HttpService } from '../../services/http.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-verify',
  templateUrl: './verify.component.html',
  styleUrls: ['./verify.component.css']
})

export class VerifyComponent implements OnInit, OnDestroy {
  routeSub : any
  successStatus = "ongoing"
  message = ""
  messageError = ""

  constructor(
    public menuService : MenuService,
    private route: ActivatedRoute,
    private httpService : HttpService) { }

  ngOnInit() {
    this.routeSub = this.route.params.subscribe(params => {
      this.successStatus = ""
      this.message = ""
      this.messageError = ""
      let token = params['token'];
      this.httpService.verify(token).subscribe(
        () => {
          this.successStatus = "ok"
          this.message = "You have successfully validated your email, please click to login"
        },
        (err) => {
          this.successStatus = "ko"
          let error = err.json()
          this.messageError = "Email validation error: "+error.error
        }
      )
    })
  }

  ngOnDestroy() {
    this.routeSub.unsubscribe();
  }

  routeToLogin() {
    this.menuService.navigate(["/auth", "signin"])
  }

}
