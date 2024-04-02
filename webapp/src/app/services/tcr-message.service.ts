import {Injectable} from '@angular/core';
import {Observable, retry} from "rxjs";
import {TcrMessage} from "../interfaces/tcr-message";
import {WebsocketService} from "./websocket.service";
import {takeUntilDestroyed} from "@angular/core/rxjs-interop";

@Injectable({
  providedIn: 'root'
})
export class TcrMessageService {
  public message$: Observable<TcrMessage>;

  constructor(public ws: WebsocketService) {
    this.message$ = this.ws.webSocket$.pipe(
      retry({delay: 5_000}),
      takeUntilDestroyed(),
    )
  }
}
