import { Component, OnInit, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import { MenuService } from '../../services/menu.service';

@Component({
  selector: 'app-stack-deploy',
  templateUrl: './docker-stack-deploy.component.html',
  styleUrls: ['./docker-stack-deploy.component.css']
})
export class DockerStackDeployComponent implements OnInit {
  @ViewChild ('f') form: NgForm;

  constructor(private menuService : MenuService) { }

  ngOnInit() {
    this.menuService.setItemMenu('stack', 'Deploy')
  }

  onDeploy() {
    console.log(this.form)
  }

}
