import {Injectable} from '@angular/core';
import {webSocket} from "rxjs/webSocket";
import {TcrMessage} from "../interfaces/tcr-message";

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {
  private readonly url = this.webSocketURL("/ws");
  private webSocketSubject = webSocket<TcrMessage>(this.url);
  public webSocket$ = this.webSocketSubject.asObservable();

  private webSocketURL(path: string): string {
    return "ws://" + window.location.host + path;
  }
}
