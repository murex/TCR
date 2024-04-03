import {AfterViewInit, Component, effect, OnInit, Signal} from '@angular/core';
import {TcrMessage} from "../../interfaces/tcr-message";
import {TcrTimerService} from "../../services/tcr-timer.service";
import {toSignal} from "@angular/core/rxjs-interop";
import {TcrTimer, TcrTimerState} from "../../interfaces/tcr-timer";
import {FormatTimerPipe} from "../../pipes/format-timer.pipe";
import {NgClass, NgIf, NgStyle} from "@angular/common";

@Component({
  selector: 'app-tcr-timer',
  standalone: true,
  imports: [
    FormatTimerPipe,
    NgClass,
    NgIf,
    NgStyle
  ],
  templateUrl: './tcr-timer.component.html',
  styleUrl: './tcr-timer.component.css'
})
export class TcrTimerComponent implements OnInit, AfterViewInit {
  timer?: TcrTimer;
  progressRatio: number | undefined;
  remaining: number | undefined;
  timeout: number | undefined;
  fgColor: string | undefined;
  timerMessage: Signal<TcrMessage | undefined>;
  private syncCounter: number = 0;
  private SYNC_INTERVAL: number = 10;

  constructor(private timerService: TcrTimerService) {
    this.timerMessage = toSignal(this.timerService.message$);

    effect(() => {
      // When receiving a timer message from the server
      // trigger a refresh query to ensure that we keep in sync
      this.refresh(this.timerMessage()!);
    });
  }

  ngOnInit(): void {
    this.getTimer();
  }

  ngAfterViewInit(): void {
    setInterval(() => this.periodicUpdate(), 1000);
  }

  // Timer periodic update. We re-sync with the server every SYNC_INTERVAL seconds
  periodicUpdate(): void {
    const activeStates = [TcrTimerState.RUNNING, TcrTimerState.TIMEOUT];
    if (this.syncCounter++ >= this.SYNC_INTERVAL) {
      this.getTimer();
      this.syncCounter = 0;
    } else if (activeStates.includes(this.timer?.state as TcrTimerState)) {
      this.remaining = this.remaining! - 1;
      this.updateColor();
    }
  }

  refresh(message: TcrMessage): void {
    if (message)
      this.getTimer();
  }

  public getTimer(): void {
    this.timerService.getTimer().subscribe(t => {
      this.timer = t;
      this.timeout = parseInt(t.timeout, 10);
      this.remaining = parseInt(t.remaining, 10);
      this.updateColor()
    });
  }

  updateColor(): void {
    let color = {red: 0, green: 0, blue: 0};
    if (this.timer) {
      switch (this.timer.state) {
        case TcrTimerState.OFF:
        case TcrTimerState.STOPPED:
          color = {red: 128, green: 128, blue: 128};
          break;
        case TcrTimerState.TIMEOUT:
          color = {red: 255, green: 0, blue: 0};
          break;
        default:
          this.progressRatio = (this.timeout! - this.remaining!) / this.timeout!;
          color = {
            red: 255,
            green: 255 * (1 - this.progressRatio),
            blue: 255 * (1 - this.progressRatio),
          };
          break;
      }
    }
    this.fgColor = `rgb(${color.red},${color.green},${color.blue})`;
  }
}
