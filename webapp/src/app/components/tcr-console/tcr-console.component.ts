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
import {Subject} from "rxjs";


@Component({
  selector: 'app-tcr-console',
  standalone: true,
  imports: [
    TcrRolesComponent,
    TcrTraceComponent
  ],
  templateUrl: './tcr-console.component.html',
  styleUrl: './tcr-console.component.css'
})
export class TcrConsoleComponent {
  title = "TCR Console";
  message$ = this.messageService.message$;
  text: Subject<string> = new Subject<string>();
  clearTrace: Subject<void> = new Subject<void>();

  constructor(private messageService: TcrMessageService) {
    this.messageService.message$.subscribe(msg => this.printMessage(msg));
  }

  printMessage(message: TcrMessage): void {
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
        if (isRoleStartMessage(message.text))
          this.clear();
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

  printSimple(text: string) {
    this.print(text);
  }

  printInfo(text: string) {
    this.print(cyan(text));
  }

  printTitle(text: string) {
    const lineSep = lightCyan("─".repeat(80));
    this.print(lineSep + "\n" + lightCyan(text));
  }

  printRole(text: string) {
    const sepLine = yellow("─".repeat(80));
    this.print(sepLine + "\n" + lightYellow(formatRoleMessage(text)) + "\n" + sepLine);
  }

  printSuccess(text: string) {
    this.print("🟢- " + green(text));
  }

  printWarning(text: string) {
    this.print("🔶- " + yellow(text));
  }

  printError(text: string) {
    this.print("🟥- " + red(text));
  }

  printUnhandled(type: string, text: string) {
    this.print(bgDarkGray("[" + type + "]") + " " + text);
  }

  print(input: string): void {
    this.text.next(input);
  }

  clear() {
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

export function isRoleStartMessage(text: string) {
  return getRoleAction(text) === ROLE_START;
}
