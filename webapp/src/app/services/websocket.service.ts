import {Inject, Injectable, InjectionToken, OnDestroy} from '@angular/core';
import {webSocket, WebSocketSubject} from "rxjs/webSocket";
import {TcrMessage} from "../interfaces/tcr-message";
import {catchError, retry, throwError} from "rxjs";
import {takeUntilDestroyed} from "@angular/core/rxjs-interop";

export const WEBSOCKET_CTOR = new InjectionToken<typeof webSocket>(
  'rxjs/webSocket.webSocket',
  {
    providedIn: 'root',
    factory: () => webSocket,
  }
);

@Injectable({
  providedIn: 'root'
})
export class WebsocketService implements OnDestroy {
  private readonly url: string = "ws://" + window.location.host + "/ws";
  public webSocket$: WebSocketSubject<TcrMessage>;

  constructor(@Inject(WEBSOCKET_CTOR) private webSocketSubject: typeof webSocket) {
    this.webSocket$ = this.webSocketSubject<TcrMessage>(this.url);
    this.webSocket$.asObservable().pipe(
      catchError((error) => {
        return throwError(() => new Error(error));
      }),
      retry({delay: 5_000}),
      takeUntilDestroyed(),
    );
  }

  ngOnDestroy(): void {
    console.info(`closed websocket connection to ${this.url}`);
    this.webSocket$.complete();
  }
}
