import {Component, effect, Input, OnInit, Signal} from '@angular/core';
import {TcrRole} from "../../interfaces/tcr-role";
import {TcrRolesService} from "../../services/trc-roles.service";
import {TcrMessage} from "../../interfaces/tcr-message";
import {NgClass, NgIf} from "@angular/common";
import {toSignal} from "@angular/core/rxjs-interop";

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
  @Input({required: true}) name = "";
  role?: TcrRole;
  roleMessage: Signal<TcrMessage | undefined>;

  constructor(private rolesService: TcrRolesService) {
    this.roleMessage = toSignal(this.rolesService.message$);

    effect(() => {
      // When receiving a role message from the server
      // trigger a refresh query to ensure that we keep in sync
      this.refresh(this.roleMessage()!);
    });
  }

  ngOnInit(): void {
    this.getRole();
  }

  refresh(message: TcrMessage): void {
    if (message) {
      const name = message.text.split(":")[0];
      if (name === this.name) {
        this.getRole();
      }
    }
  }

  private getRole(): void {
    this.rolesService.getRole(this.name)
      .subscribe(r => {
        this.role = r;
      });
  }

  toggleRole(role: TcrRole) {
    this.rolesService.activateRole(role.name, !role.active)
      .subscribe(r => {
        console.log(r.name + ' set to ' + r.active);
        this.role = r;
      });
  }
}
