import {Component, OnInit, Signal} from '@angular/core';
import {TcrMessage} from "../../interfaces/tcr-message";
import {TcrTimerService} from "../../services/tcr-timer.service";
import {toSignal} from "@angular/core/rxjs-interop";

@Component({
  selector: 'app-notification',
  standalone: true,
  imports: [],
  templateUrl: './notification.component.html',
  styleUrl: './notification.component.css'
})
export class NotificationComponent implements OnInit {
  message: Signal<TcrMessage | undefined>;

  constructor(
    private tcrTimerService: TcrTimerService) {
    this.message = toSignal(this.tcrTimerService.webSocket$);
  }

  ngOnInit(): void {
  }
}
