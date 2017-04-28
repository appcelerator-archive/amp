import { Component, OnInit } from '@angular/core';
import { OrganizationsService } from '../services/organizations.service'

@Component({
  selector: 'app-organizations',
  templateUrl: './organizations.component.html',
  styleUrls: ['./organizations.component.css']
})
export class OrganizationsComponent implements OnInit {

  constructor(public organizationsService : OrganizationsService) { }

  ngOnInit() {
  }

}
