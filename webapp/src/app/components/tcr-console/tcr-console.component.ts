import {Component, effect, OnInit, Signal, ViewChild} from '@angular/core';
import {NgTerminal, NgTerminalModule} from "ng-terminal";
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

@Component({
  selector: 'app-tcr-console',
  standalone: true,
  imports: [
    NgTerminalModule,
    TcrRolesComponent
  ],
  templateUrl: './tcr-console.component.html',
  styleUrl: './tcr-console.component.css'
})
export class TcrConsoleComponent implements OnInit {
  title = "TCR Console";
  tcrMessage: Signal<TcrMessage | undefined>;

  @ViewChild('term', {static: false}) child!: NgTerminal;

  constructor(private messageService: TcrMessageService) {
    this.tcrMessage = toSignal(this.messageService.message$);

    effect(() => {
      // When receiving a message from the server
      // print it in the terminal
      this.printMessage(this.tcrMessage()!);
    });
  }

  ngOnInit(): void {
    this.clear();
  }

  private printMessage(message: TcrMessage): void {
    if (message === undefined) {
      return;
    }
    switch (message.type) {
      case TcrMessageType.SIMPLE:
        this.print(message.text);
        break;
      case TcrMessageType.INFO:
        this.print(cyan(message.text))
        break;
      case TcrMessageType.TITLE:
        this.print(lightCyan("─".repeat(80)));
        this.print(lightCyan(message.text));
        break;
      case TcrMessageType.ROLE:
        if (getRoleAction(message.text) === "start") {
          this.clear();
        }
        this.print(yellow("─".repeat(80)));
        this.print(lightYellow(formatRoleMessage(message.text)));
        this.print(yellow("─".repeat(80)));
        break;
      case TcrMessageType.TIMER:
        // ignore: handled by timer service
        break;
      case TcrMessageType.SUCCESS:
        this.print("🟢 " + green(message.text));
        break;
      case TcrMessageType.WARNING:
        this.print("🔶 " + yellow(message.text));
        break;
      case TcrMessageType.ERROR:
        this.print("🟥 " + red(message.text));
        break;
      default:
        this.print(bgDarkGray("[" + message.type + "]") + " " + message.text);
    }
  }

  private print(input: string) {
    // ng-console handles EOL in Windows style, e.g. it needs CRLF to properly
    // go back to beginning of next line in the console
    this.child.write(input.replace(/\n/g, "\r\n") + "\r\n");
  }

  private clear() {
    if (this.child)
      this.child.underlying?.reset();
  }

}

function getRoleAction(message: string): string {
  return message ? message.split(":")[1] : "";
}

function getRoleName(message: string): string {
  return message ? message.split(":")[0] : "";
}

function formatRoleMessage(message: string): string {
  return message
    ? capitalize(getRoleAction(message)) + "ing "
    + capitalize(getRoleName(message)) + " role"
    : "";
}

function capitalize(s: string): string {
  return s.charAt(0).toUpperCase() + s.slice(1);
}
