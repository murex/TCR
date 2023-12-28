import {Injectable} from '@angular/core';
import {webSocket} from "rxjs/webSocket";

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {
  private readonly url = this.webSocketURL("/ws");
  private webSocketSubject = webSocket<string>(this.url);
  public webSocket$ = this.webSocketSubject.asObservable();

  private webSocketURL(path: string): string {
    return "ws://" + window.location.host + path;
  }
}
