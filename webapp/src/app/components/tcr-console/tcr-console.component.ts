import {Component, ViewChild} from '@angular/core';
import {WebsocketService} from "../../services/websocket.service";
import {NgTerminal, NgTerminalModule} from "ng-terminal";
import {TcrMessage} from "../../interfaces/tcr-message";
import {bgDarkGray, cyan, green, lightCyan, red, yellow} from "ansicolor";
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
    this.ws.webSocket$.subscribe((m: TcrMessage) => this.writeMessage(m));
  }

  private writeMessage(message: TcrMessage): void {
    switch (message.type) {
      case "simple":
        this.write(message.text);
        break;
      case "info":
        this.write(cyan(message.text))
        break;
      case "title":
        this.write(lightCyan("─".repeat(80)));
        this.write(lightCyan(message.text));
        break;
      case "role":
        // ignore: handled by roles service
        break;
      case "timer":
        this.write("⏳ " + green(message.text));
        break;
      case "success":
        this.write("🟢 " + green(message.text));
        break;
      case "warning":
        this.write("🔶 " + yellow(message.text));
        break;
      case "error":
        this.write("🟥 " + red(message.text));
        break;
      default:
        this.write(bgDarkGray("[" + message.type + "]") + " " + message.text);
    }
  }

  private write(input: string) {
    // ng-console handles EOL in Windows style, e.g. it needs CRLF to properly
    // go back to beginning of next line in the console
    this.child.write(input.replace(/\n/g, "\r\n") + "\r\n");
  };

}
