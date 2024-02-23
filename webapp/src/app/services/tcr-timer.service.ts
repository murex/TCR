import {Injectable} from '@angular/core';
import {filter, Observable} from "rxjs";
import {TcrMessage} from "../interfaces/tcr-message";
import {WebsocketService} from "./websocket.service";

@Injectable({
  providedIn: 'root'
})
export class TcrTimerService {
  public webSocket$: Observable<TcrMessage>;

  constructor(
    private ws: WebsocketService) {
    this.webSocket$ = this.ws.webSocket$.pipe(
      filter(message => message.type === "timer")
    )
  }

}
