import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {catchError, Observable, of} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class TcrControlsService {
  private apiUrl: string = '/api'; // URL to web api

  constructor(private http: HttpClient) {
  }

  abortCommand(): Observable<unknown> {
    return this.sendControl(`abort-command`);
  }

  private sendControl(command: string) {
    const url: string = `${this.apiUrl}/controls/${command}`;
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      })
    };

    return this.http.post(url, httpOptions).pipe(
      catchError(this.handleError<unknown>(command)),
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
