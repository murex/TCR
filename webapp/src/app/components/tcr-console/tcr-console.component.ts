import {Component, ViewChild} from '@angular/core';
import {WebsocketService} from "../../services/websocket.service";
import {catchError, retry, throwError} from "rxjs";
import {takeUntilDestroyed} from "@angular/core/rxjs-interop";
import {NgTerminal, NgTerminalModule} from "ng-terminal";

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
      .subscribe((value) => {
        this.write(value);
      });
  }

  private write(input: string) {
    this.child.write(input + "\r\n");
  };

}
