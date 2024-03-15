import {Component, ViewChild} from '@angular/core';
import {WebsocketService} from "../../services/websocket.service";
import {NgTerminal, NgTerminalModule} from "ng-terminal";
import {TcrMessage} from "../../interfaces/tcr-message";
import {bgDarkGray, cyan, green, lightCyan, lightYellow, red, yellow} from "ansicolor";
import {TcrRolesComponent} from "../tcr-roles/tcr-roles.component";

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
export class TcrConsoleComponent {
  title = "TCR Console";
  @ViewChild('term', {static: false}) child!: NgTerminal;

  constructor(private ws: WebsocketService) {
    this.ws.webSocket$.subscribe((m: TcrMessage) => this.printMessage(m));
  }

  private printMessage(message: TcrMessage): void {
    switch (message.type) {
      case "simple":
        this.print(message.text);
        break;
      case "info":
        this.print(cyan(message.text))
        break;
      case "title":
        this.print(lightCyan("─".repeat(80)));
        this.print(lightCyan(message.text));
        break;
      case "role":
        if (getRoleAction(message.text) === "start") {
          this.clear();
        }
        this.print(yellow("─".repeat(80)));
        this.print(lightYellow(formatRoleMessage(message.text)));
        this.print(yellow("─".repeat(80)));
        break;
      case "timer":
        // ignore: handled by timer service
        // this.print("⏳ " + green(message.text));
        break;
      case "success":
        this.print("🟢 " + green(message.text));
        break;
      case "warning":
        this.print("🔶 " + yellow(message.text));
        break;
      case "error":
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
