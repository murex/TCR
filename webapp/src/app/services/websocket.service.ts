import {Injectable} from '@angular/core';
import {webSocket, WebSocketSubject} from "rxjs/webSocket";
import {TcrMessage} from "../interfaces/tcr-message";
import {catchError, Observable, retry, throwError} from "rxjs";
import {takeUntilDestroyed} from "@angular/core/rxjs-interop";

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {
  private readonly url = this.webSocketURL("/ws");
  private webSocketSubject: WebSocketSubject<TcrMessage>;
  public webSocket$: Observable<TcrMessage>;

  constructor() {
    this.webSocketSubject = webSocket<TcrMessage>(this.url);
    this.webSocket$ = this.webSocketSubject.asObservable().pipe(
      catchError((error) => {
        return throwError(() => new Error(error));
      }),
      retry({delay: 5_000}),
      takeUntilDestroyed()
    )
  }

  private webSocketURL(path: string): string {
    return "ws://" + window.location.host + path;
  }
}
