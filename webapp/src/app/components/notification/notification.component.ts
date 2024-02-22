import {Component, OnInit, signal} from '@angular/core';
import {catchError, retry, throwError} from "rxjs";
import {TcrMessage} from "../../interfaces/tcr-message";
import {TcrTimerService} from "../../services/tcr-timer.service";

@Component({
  selector: 'app-notification',
  standalone: true,
  imports: [],
  templateUrl: './notification.component.html',
  styleUrl: './notification.component.css'
})
export class NotificationComponent implements OnInit {
  message = signal("timeout");

  constructor(
    private tcrTimerService: TcrTimerService) {
  }

  ngOnInit(): void {
    this.tcrTimerService.webSocket$
      .pipe(
        catchError((error) => {
          return throwError(() => new Error(error));
        }),
        retry({delay: 5_000}))
      .subscribe((m: TcrMessage) => this.message.set(m.text));
  }
  
}
