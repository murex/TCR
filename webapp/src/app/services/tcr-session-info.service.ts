import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {catchError, Observable, of} from "rxjs";
import {TcrSessionInfo} from "../interfaces/tcr-session-info";

@Injectable({
  providedIn: 'root'
})
export class TcrSessionInfoService {
  private apiUrl = `/api`// URL to web api

  constructor(
    private http: HttpClient) {
  }

  getSessionInfo(): Observable<TcrSessionInfo> {
    const url = `${this.apiUrl}/session-info`;
    const httpOptions = {
      headers: new HttpHeaders({
        'Accept': 'application/json',
      })
    };

    return this.http.get<TcrSessionInfo>(url, httpOptions)
      .pipe(
        catchError(this.handleError<TcrSessionInfo>('getSessionInfo'))
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
