import {Component} from '@angular/core';
import {NgForOf} from "@angular/common";
import {TcrRoleComponent} from "../tcr-role/tcr-role.component";

@Component({
  selector: 'app-tcr-roles',
  standalone: true,
  imports: [
    NgForOf,
    TcrRoleComponent
  ],
  templateUrl: './tcr-roles.component.html',
  styleUrl: './tcr-roles.component.css'
})
export class TcrRolesComponent {
  roles: string[] = ['driver', 'navigator'];

  constructor() {
  }
}
