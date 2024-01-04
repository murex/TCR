import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {catchError, filter, Observable, of, tap} from "rxjs";
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
        tap(_ => this.log('fetched TCR roles')),
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
        tap(r => this.log(`fetched TCR role ${r.name}`)),
        catchError(this.handleError<TcrRole>('getRole'))
      );
  }

  activateRole(name: string, state: boolean): Observable<TcrRole> {
    this.log("name: ${name} state: ${state}");

    const url = `${this.apiUrl}/roles/${name}/${state ? "start" : "stop"}`;
    const httpOptions = {
      headers: new HttpHeaders({
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      })
    };

    return this.http.post<TcrRole>(url, httpOptions).pipe(
      tap((role: TcrRole) =>
        this.log(`set role name=${role.name} to active=${role.active}`)),
      catchError(this.handleError<TcrRole>('toggleRole'))
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
