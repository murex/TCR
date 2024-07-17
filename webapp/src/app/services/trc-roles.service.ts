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

import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {catchError, filter, Observable, of, retry} from "rxjs";
import {TcrRole} from "../interfaces/tcr-role";
import {WebsocketService} from "./websocket.service";
import {TcrMessage, TcrMessageType} from "../interfaces/tcr-message";
import {takeUntilDestroyed} from "@angular/core/rxjs-interop";

@Injectable({
  providedIn: 'root'
})
export class TcrRolesService {
  private apiUrl: string = `/api`; // URL to web api
  public message$: Observable<TcrMessage>;

  constructor(private http: HttpClient, private ws: WebsocketService) {
    this.message$ = this.ws.webSocket$.pipe(
      filter(message => message.type === TcrMessageType.ROLE),
      retry({delay: 5_000}),
      takeUntilDestroyed(),
    )
  }

  getRole(name: string): Observable<TcrRole> {
    const url: string = `${this.apiUrl}/roles/${name}`;
    const httpOptions = {
      headers: new HttpHeaders({
        'Accept': 'application/json',
      })
    };

    return this.http.get<TcrRole>(url, httpOptions).pipe(
      catchError(this.handleError<TcrRole>('getRole'))
    );
  }

  activateRole(name: string, state: boolean): Observable<TcrRole> {
    const url: string = `${this.apiUrl}/roles/${name}/${state ? "start" : "stop"}`;
    const httpOptions = {
      headers: new HttpHeaders({
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      })
    };

    return this.http.post<TcrRole>(url, httpOptions).pipe(
      catchError(this.handleError<TcrRole>('activateRole')),
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
