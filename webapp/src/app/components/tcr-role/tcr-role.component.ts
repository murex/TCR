import {Component, Input, OnInit} from '@angular/core';
import {TcrRole} from "../../interfaces/tcr-role";
import {TcrRolesService} from "../../services/trc-roles.service";
import {TcrMessage} from "../../interfaces/tcr-message";
import {NgClass, NgIf} from "@angular/common";
import {catchError, retry, throwError} from "rxjs";

@Component({
  selector: 'app-tcr-role',
  standalone: true,
  imports: [
    NgIf,
    NgClass
  ],
  templateUrl: './tcr-role.component.html',
  styleUrl: './tcr-role.component.css'
})
export class TcrRoleComponent implements OnInit {
  @Input() name = "";
  @Input() role?: TcrRole;

  constructor(
    private rolesService: TcrRolesService) {
  }

  ngOnInit(): void {
    this.getRole();
    this.rolesService.webSocket$
      .pipe(
        catchError((error) => {
          return throwError(() => new Error(error));
        }),
        retry({delay: 5_000}))
      .subscribe((m: TcrMessage) => this.refresh(m));
  }

  private refresh(message: TcrMessage): void {
    const name = message.text.split(":")[0];
    if (name === this.name) {
      this.getRole();
    }
  }

  private getRole(): void {
    this.rolesService.getRole(this.name).subscribe(r => {
        this.role = r
      }
    );
  }

  toggleRole(role: TcrRole) {
    this.rolesService.activateRole(role.name, !role.active)
      .subscribe(r => {
        console.log(r.name + ' set to ' + r.active);
      });
  }
}
