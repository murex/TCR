import {Injectable} from '@angular/core';
import {catchError, filter, Observable, of, tap} from "rxjs";
import {TcrMessage} from "../interfaces/tcr-message";
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
      filter(message => message.type === "timer")
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
        tap(t => this.log(`fetched TCR timer ${t.state}`)),
        catchError(this.handleError<TcrTimer>('getTimer'))
      );
  }

  private log(_message: string) {
    // TODO - add messageService component
    // this.messageService.add(`AlbumService: ${message}`);
  }

  /**
   * Handle HTTP operation that failed.
   * Let the app continue.
   *
   * @param operation - name of the operation that failed
   * @param result - optional value to return as the observable result
   */
  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {

      // TODO: send the error to remote logging infrastructure
      console.error(error); // log to console instead

      // TODO: better job of transforming error for user consumption
      this.log(`${operation} failed: ${error.message}`);

      // Let the app keep running by returning an empty result.
      return of(result as T);
    };
  }
}
