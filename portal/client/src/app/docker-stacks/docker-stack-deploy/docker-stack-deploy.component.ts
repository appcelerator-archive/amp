import { Component, OnInit, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import { MenuService } from '../../services/menu.service';
import { HttpService } from '../../services/http.service';
import { ActivatedRoute, Params } from '@angular/router';

@Component({
  selector: 'app-stack-deploy',
  templateUrl: './docker-stack-deploy.component.html',
  styleUrls: ['./docker-stack-deploy.component.css']
})
export class DockerStackDeployComponent implements OnInit {
  @ViewChild ('f') form: NgForm;
  message = ""

  constructor(
    private menuService : MenuService,
    private route : ActivatedRoute,
    private httpService : HttpService) { }
    private updateStackName = "new stack"
    fileText=""

  ngOnInit() {
    this.menuService.setItemMenu('stack', 'Deploy')
    this.route.params.subscribe(
      (params : Params) => {
        if (params['stackName']) {
          this.menuService.setItemMenu('stack', params['stackName'])
          this.updateStackName = params['stackName']
        }
      }
    );
  }

  onDeploy(form : NgForm) {
    this.message = "submitting..."
    let name = form.value.name
    if (this.updateStackName != 'new stack') {
      name = this.updateStackName
    }

    this.httpService.deployStack(name, form.value.filedata).subscribe(
      data => {
        this.message=""
        this.menuService.returnToPreviousPath()
      },
      error => {
        let data = error.json()
        this.message = data.error
      }
    )
  }

  fileSelected(event) {
    let files = event.srcElement.files;
    if (!files || !files[0]) {
      return
    }
    var reader = new FileReader();
    reader.onload = file => {
      var contents: any = file.target;
      this.fileText= contents.result;
    };
    reader.readAsText(files[0]);
  }

  returnBack() {
    this.menuService.returnToPreviousPath()
  }

}
