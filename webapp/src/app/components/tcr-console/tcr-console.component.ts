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

import {Component} from '@angular/core';
import {TcrMessage, TcrMessageType} from "../../interfaces/tcr-message";
import {TcrRolesComponent} from "../tcr-roles/tcr-roles.component";
import {TcrMessageService} from "../../services/tcr-message.service";
import {
  bgDarkGray,
  cyan,
  green,
  lightCyan,
  lightYellow,
  red,
  yellow
} from "ansicolor";
import {TcrTraceComponent} from "../tcr-trace/tcr-trace.component";
import {Observable, Subject} from "rxjs";
import {TcrControlsComponent} from "../tcr-controls/tcr-controls.component";

@Component({
  selector: 'app-tcr-console',
  imports: [
    TcrRolesComponent,
    TcrTraceComponent,
    TcrControlsComponent
  ],
  templateUrl: './tcr-console.component.html',
  styleUrl: './tcr-console.component.css',
})
export class TcrConsoleComponent {
  title: string = "TCR Console";
  message$: Observable<TcrMessage> = this.messageService.message$;
  text: Subject<string> = new Subject<string>();
  clearTrace: Subject<void> = new Subject<void>();

  constructor(private messageService: TcrMessageService) {
    this.messageService.message$.subscribe(msg => this.printMessage(msg));
  }

  printMessage(message: TcrMessage): void {
    // clear the console every time a role is starting
    if (isRoleStartMessage(message))
      this.clear();

    switch (message.type) {
      case TcrMessageType.SIMPLE:
        this.printSimple(message.text);
        break;
      case TcrMessageType.INFO:
        this.printInfo(message.text);
        break;
      case TcrMessageType.TITLE:
        this.printTitle(message.text);
        break;
      case TcrMessageType.ROLE:
        this.printRole(message.text);
        break;
      case TcrMessageType.TIMER:
        // ignore: handled by timer service
        break;
      case TcrMessageType.SUCCESS:
        this.printSuccess(message.text);
        break;
      case TcrMessageType.WARNING:
        this.printWarning(message.text);
        break;
      case TcrMessageType.ERROR:
        this.printError(message.text);
        break;
      default:
        this.printUnhandled(message.type, message.text);
    }
  }

  printSimple(text: string): void {
    this.print(text);
  }

  printInfo(text: string): void {
    this.print(cyan(text));
  }

  printTitle(text: string): void {
    const lineSep = lightCyan("─".repeat(80));
    this.print(lineSep + "\n" + lightCyan(text));
  }

  printRole(text: string): void {
    const sepLine = yellow("─".repeat(80));
    this.print(sepLine + "\n" + lightYellow(formatRoleMessage(text)) + "\n" + sepLine);
  }

  printSuccess(text: string): void {
    this.print("🟢 " + green(text));
  }

  printWarning(text: string): void {
    this.print("🔶 " + yellow(text));
  }

  printError(text: string): void {
    this.print("🟥 " + red(text));
  }

  printUnhandled(type: string, text: string): void {
    this.print(bgDarkGray("[" + type + "]") + " " + text);
  }

  print(input: string): void {
    this.text.next(input);
  }

  clear(): void {
    this.clearTrace.next();
  }
}

export function getRoleAction(message: string): string {
  return message ? message.split(":")[1] : "";
}

export function getRoleName(message: string): string {
  return message ? message.split(":")[0] : "";
}

export function formatRoleMessage(message: string): string {
  const action = getRoleAction(message);
  return action ? capitalize(action) + "ing " + getRoleName(message) + " role" : "";
}

function capitalize(text: string): string {
  return text.charAt(0).toUpperCase() + text.slice(1);
}

const ROLE_START = "start";

export function isRoleStartMessage(msg: TcrMessage): boolean {
  return msg.type === TcrMessageType.ROLE && getRoleAction(msg.text) === ROLE_START;
}
