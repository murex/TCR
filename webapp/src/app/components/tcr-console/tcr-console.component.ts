import {Component} from '@angular/core';
import {WebsocketService} from "../../services/websocket.service";
import {catchError, retry, throwError} from "rxjs";
import {takeUntilDestroyed} from "@angular/core/rxjs-interop";
import {NgIf} from "@angular/common";

@Component({
  selector: 'app-tcr-console',
  standalone: true,
  imports: [
    NgIf
  ],
  templateUrl: './tcr-console.component.html',
  styleUrl: './tcr-console.component.css'
})
export class TcrConsoleComponent {
  title = "TCR Console";
  public trace: String = "";

  constructor(private ws: WebsocketService) {
    this.ws.webSocket$
      .pipe(
        catchError((error) => {
          this.trace = "";
          return throwError(() => new Error(error));
        }),
        retry({delay: 5_000}),
        takeUntilDestroyed())
      .subscribe((value) => {
        this.trace += value + "\n";
      });
  }
}
