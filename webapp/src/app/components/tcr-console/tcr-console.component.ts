import {Component, effect, signal, Signal} from '@angular/core';
import {TcrMessage, TcrMessageType} from "../../interfaces/tcr-message";
import {TcrRolesComponent} from "../tcr-roles/tcr-roles.component";
import {TcrMessageService} from "../../services/tcr-message.service";
import {toSignal} from "@angular/core/rxjs-interop";
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
  tcrMessage: Signal<TcrMessage | undefined>;
  text = signal("");
  clearTrace: Subject<void> = new Subject<void>();

  constructor(private messageService: TcrMessageService) {
    this.tcrMessage = toSignal(this.messageService.message$);

    effect(() => {
      // When receiving a message from the server
      // print it in the terminal
      this.printMessage(this.tcrMessage()!);
    }, {allowSignalWrites: true});
  }

  private printMessage(message: TcrMessage): void {
    if (message === undefined) {
      return;
    }
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
        this.printUnhandled(message);
    }
  }

  private printSimple(text: string) {
    this.print(text);
  }

  private printInfo(text: string) {
    this.print(cyan(text));
  }

  private printTitle(text: string) {
    const lineSep = lightCyan("─".repeat(80));
    this.print(lineSep + "\n" + lightCyan(text));
  }

  private printRole(text: string) {
    const sepLine = yellow("─".repeat(80));
    this.print(sepLine + "\n" + lightYellow(formatRoleMessage(text)) + "\n" + sepLine);
  }

  private printSuccess(text: string) {
    this.print("🟢- " + green(text));
  }

  private printWarning(text: string) {
    this.print("🔶- " + yellow(text));
  }

  private printError(text: string) {
    this.print("🟥- " + red(text));
  }

  private printUnhandled(message: TcrMessage) {
    this.print(bgDarkGray("[" + message.type + "]") + " " + message.text);
  }

  print(input: string): void {
    this.text.set(input);
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
