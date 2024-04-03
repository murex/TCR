import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {catchError, Observable, of} from "rxjs";
import {TcrBuildInfo} from "../interfaces/tcr-build-info";

@Injectable({
  providedIn: 'root'
})
export class TcrBuildInfoService {
  private apiUrl: string = `/api` // URL to web api

  constructor(
    private http: HttpClient) {
  }

  getBuildInfo(): Observable<TcrBuildInfo> {
    const url: string = `${this.apiUrl}/build-info`;
    const httpOptions = {
      headers: new HttpHeaders({
        'Accept': 'application/json',
      })
    };

    return this.http.get<TcrBuildInfo>(url, httpOptions)
      .pipe(
        catchError(this.handleError<TcrBuildInfo>('getBuildInfo'))
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
