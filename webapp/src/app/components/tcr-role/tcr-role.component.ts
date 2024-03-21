import {Component, effect, Input, OnInit, Signal} from '@angular/core';
import {TcrRole} from "../../interfaces/tcr-role";
import {TcrRolesService} from "../../services/trc-roles.service";
import {TcrMessage} from "../../interfaces/tcr-message";
import {NgClass, NgIf, NgStyle} from "@angular/common";
import {toSignal} from "@angular/core/rxjs-interop";
import {TcrTimerService} from "../../services/tcr-timer.service";
import {TcrTimer} from "../../interfaces/tcr-timer";

@Component({
  selector: 'app-tcr-role',
  standalone: true,
  imports: [
    NgIf,
    NgClass,
    NgStyle
  ],
  templateUrl: './tcr-role.component.html',
  styleUrl: './tcr-role.component.css'
})
export class TcrRoleComponent implements OnInit {
  @Input() name = "";
  @Input() role?: TcrRole;
  @Input() timer?: TcrTimer;
  @Input() remaining: number | undefined;
  @Input() timeout: number | undefined;
  @Input() fgColor: string | undefined;
  roleMessage: Signal<TcrMessage | undefined>;
  timerMessage: Signal<TcrMessage | undefined>;

  constructor(
    private rolesService: TcrRolesService,
    private timerService: TcrTimerService) {
    this.roleMessage = toSignal(this.rolesService.webSocket$);
    this.timerMessage = toSignal(this.timerService.webSocket$);


    effect(() => {
      // When receiving a role message from the server
      // trigger a refresh query to ensure that we keep in sync
      this.refreshRole(this.roleMessage()!);
      this.refreshColorOfTime(this.timerMessage()!);
    });
  }

  ngOnInit(): void {
    this.getRole();
    this.getTimer();
  }

  private refreshColorOfTime(message: TcrMessage): void {
    if(message) {
      this.getTimer();
    }
  }

  private refreshRole(message: TcrMessage): void {
    if (message) {
      const name = message.text.split(":")[0];
      if (name === this.name) {
        this.getRole();
      }
    }
  }

  private getRole(): void {
    this.rolesService.getRole(this.name).subscribe(r => {
        this.role = r;
      }
    );
  }

  toggleRole(role: TcrRole) {
    this.rolesService.activateRole(role.name, !role.active)
      .subscribe(r => {
        console.log(r.name + ' set to ' + r.active);
      });
  }

  public getTimer(): void {
    this.timerService.getTimer().subscribe(t => {
        this.timer = t;
        this.timeout = parseInt(t.timeout, 10);
        this.remaining = parseInt(t.remaining, 10);
        this.updateColor()
      }
    );
  }

  private updateColor(): void {
    if (this.timer) {
      if (this.remaining! < 0) {
        this.fgColor = `rgb(255, 0, 0)`
      }
    }
  }
}
