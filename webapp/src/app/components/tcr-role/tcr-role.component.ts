/*
Copyright (c) 2024 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
  @Input({required: true}) name: string = "";
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

  toggleRole(role: TcrRole): void {
    this.rolesService.activateRole(role.name, !role.active)
      .subscribe(r => {
        console.log(r.name + ' set to ' + r.active);
        this.role = r;
      });
  }
}
