import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {catchError, filter, Observable, of} from "rxjs";
import {TcrRole} from "../interfaces/tcr-role";
import {WebsocketService} from "./websocket.service";
import {TcrMessage} from "../interfaces/tcr-message";

@Injectable({
  providedIn: 'root'
})
export class TcrRolesService {
  private apiUrl = `/api`; // URL to web api
  public webSocket$: Observable<TcrMessage>;

  constructor(
    private http: HttpClient,
    private ws: WebsocketService) {
    this.webSocket$ = this.ws.webSocket$.pipe(
      filter(message => message.type === "role")
    )
  }

  getRoles(): Observable<TcrRole[]> {
    const url = `${this.apiUrl}/roles`;
    const httpOptions = {
      headers: new HttpHeaders({
        'Accept': 'application/json',
      })
    };

    return this.http.get<TcrRole[]>(url, httpOptions)
      .pipe(
        catchError(this.handleError<TcrRole[]>('getRoles', []))
      );
  }

  getRole(name: string): Observable<TcrRole> {
    const url = `${this.apiUrl}/roles/${name}`;
    const httpOptions = {
      headers: new HttpHeaders({
        'Accept': 'application/json',
      })
    };

    return this.http.get<TcrRole>(url, httpOptions)
      .pipe(
        catchError(this.handleError<TcrRole>('getRole'))
      );
  }

  activateRole(name: string, state: boolean): Observable<TcrRole> {
    const url = `${this.apiUrl}/roles/${name}/${state ? "start" : "stop"}`;
    const httpOptions = {
      headers: new HttpHeaders({
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      })
    };

    return this.http.post<TcrRole>(url, httpOptions).pipe(
      catchError(this.handleError<TcrRole>('toggleRole'))
    );
  }

  /**
   * Handle HTTP operation that failed.
   * Let the app continue.
   *
   * @param operation - name of the operation that failed
   * @param result - optional value to return as the observable result
   */
  private handleError<T>(operation = 'operation', result?: T) {
    return (error: unknown): Observable<T> => {
      console.error(`${operation} - ` + error);
      // Let the app keep running by returning an empty result.
      return of(result as T);
    };
  }
}
