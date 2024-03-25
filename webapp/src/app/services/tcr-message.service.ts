import {Injectable} from '@angular/core';
import {Observable, retry} from "rxjs";
import {TcrMessage} from "../interfaces/tcr-message";
import {WebsocketService} from "./websocket.service";
import {takeUntilDestroyed} from "@angular/core/rxjs-interop";

@Injectable({
  providedIn: 'root'
})
export class TcrMessageService {
  public webSocket$: Observable<TcrMessage>;

  constructor(private ws: WebsocketService) {
    this.webSocket$ = this.ws.webSocket$.pipe(
      retry({delay: 5_000}),
      takeUntilDestroyed(),
    )
  }
}
