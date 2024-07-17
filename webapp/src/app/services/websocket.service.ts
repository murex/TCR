/*
Copyright (c) 2024 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
