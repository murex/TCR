import {Component, ViewChild} from '@angular/core';
import {WebsocketService} from "../../services/websocket.service";
import {catchError, retry, throwError} from "rxjs";
import {takeUntilDestroyed} from "@angular/core/rxjs-interop";
import {NgTerminal, NgTerminalModule} from "ng-terminal";
import {TcrMessage} from "../../interfaces/tcr-message";
import {bgWhite, cyan, green, lightCyan, red, yellow} from "ansicolor";

@Component({
  selector: 'app-tcr-console',
  standalone: true,
  imports: [
    NgTerminalModule
  ],
  templateUrl: './tcr-console.component.html',
  styleUrl: './tcr-console.component.css'
})
export class TcrConsoleComponent {
  title = "TCR Console";
  @ViewChild('term', {static: false}) child!: NgTerminal;

  constructor(private ws: WebsocketService) {
    this.ws.webSocket$
      .pipe(
        catchError((error) => {
          return throwError(() => new Error(error));
        }),
        retry({delay: 5_000}),
        takeUntilDestroyed())
      .subscribe((m: TcrMessage) => this.writeMessage(m));
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
        this.write(bgWhite("[" + message.type + "]") + " " + message.text);
    }
  }

  private write(input: string) {
    // ng-console handles EOL in Windows style, e.g. it needs CRLF to properly
    // go back to beginning of next line in the console
    this.child.write(input.replace(/\n/g, "\r\n") + "\r\n");
  };

}
