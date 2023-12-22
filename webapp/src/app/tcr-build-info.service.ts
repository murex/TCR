import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {catchError, Observable, of, tap} from "rxjs";
import {TcrBuildInfo} from "./tcr-build-info";

@Injectable({
  providedIn: 'root'
})
export class TcrBuildInfoService {
  private apiUrl = `/api`// URL to web api

  constructor(
    private http: HttpClient) {
  }

  getBuildInfo(): Observable<TcrBuildInfo> {
    const url = `${this.apiUrl}/build-info`;
    const httpOptions = {
      headers: new HttpHeaders({
        'Accept': 'application/json',
      })
    };

    // TODO replace with HTTP call
    // const stub: TcrBuildInfo = {
    //   version: "0.12.0",
    //   os: "Windows",
    //   arch: "amd64",
    //   commit: "00000",
    //   date: "22/12/2023",
    //   author: "dmenanteau",
    // };
    // return of(stub);

    return this.http.get<TcrBuildInfo>(url, httpOptions)
      .pipe(
        tap(_ => this.log('fetched TCR build info')),
        catchError(this.handleError<TcrBuildInfo>('getBuildInfo'))
      );
  }

  private log(message: string) {
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
