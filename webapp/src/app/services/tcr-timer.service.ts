import {Injectable} from '@angular/core';
import {catchError, filter, Observable, of} from "rxjs";
import {TcrMessage, TcrMessageType} from "../interfaces/tcr-message";
import {WebsocketService} from "./websocket.service";
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {TcrTimer} from "../interfaces/tcr-timer";

@Injectable({
  providedIn: 'root'
})
export class TcrTimerService {
  private apiUrl = `/api`; // URL to web api
  public webSocket$: Observable<TcrMessage>;

  constructor(
    private http: HttpClient,
    private ws: WebsocketService) {
    this.webSocket$ = this.ws.webSocket$.pipe(
      filter(message => message.type === TcrMessageType.TIMER)
    )
  }

  getTimer(): Observable<TcrTimer> {
    const url = `${this.apiUrl}/timer`;
    const httpOptions = {
      headers: new HttpHeaders({
        'Accept': 'application/json',
      })
    };

    return this.http.get<TcrTimer>(url, httpOptions)
      .pipe(
        catchError(this.handleError<TcrTimer>('getTimer'))
      );
  }

  /**
   * Handle HTTP operation that failed.
   * Let the app continue.
   *
   * @param operation - name of the operation that failed
   * @param result - optional value to return as the observable result
   */
  private handleError<T>(operation: string, result?: T) {
    return (error: unknown): Observable<T> => {
      console.error(`${operation} - ` + error);
      // Let the app keep running by returning an empty result.
      return of(result as T);
    };
  }
}
