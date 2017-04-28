import { Component, OnInit, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';

@Component({
  selector: 'app-stack-deploy',
  templateUrl: './docker-stack-deploy.component.html',
  styleUrls: ['./docker-stack-deploy.component.css']
})
export class DockerStackDeployComponent implements OnInit {
  @ViewChild ('f') form: NgForm;

  constructor() { }

  ngOnInit() {
  }

  onDeploy() {
    console.log(this.form)
  }

}
