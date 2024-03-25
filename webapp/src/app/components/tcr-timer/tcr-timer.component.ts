import {AfterViewInit, Component, effect, OnInit, Signal} from '@angular/core';
import {TcrMessage} from "../../interfaces/tcr-message";
import {TcrTimerService} from "../../services/tcr-timer.service";
import {toSignal} from "@angular/core/rxjs-interop";
import {TcrTimer} from "../../interfaces/tcr-timer";
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
  private syncCounter = 0;
  private SYNC_INTERVAL = 10;

  constructor(private timerService: TcrTimerService) {
    this.timerMessage = toSignal(this.timerService.message$);

    effect(() => {
      // When receiving a timer message from the server
      // trigger a refresh query to ensure that we keep in sync
      this.refresh(this.timerMessage()!);
    });
  }

  ngAfterViewInit(): void {
    // Timer periodic update. We re-sync with the server every SYNC_INTERVAL seconds
    setInterval(() => {
      if (this.syncCounter++ >= this.SYNC_INTERVAL) {
        this.getTimer();
        this.syncCounter = 0;
      } else if (this.timer?.state != 'off') {
        this.remaining = this.remaining! - 1;
        this.updateColor();
      }
    }, 1000);
  }

  ngOnInit(): void {
    this.getTimer();
  }

  private refresh(message: TcrMessage): void {
    if (message) {
      this.getTimer();
    }
  }

  public getTimer(): void {
    this.timerService.getTimer().subscribe(t => {
        this.timer = t;
        this.timeout = parseInt(t.timeout, 10);
        this.remaining = parseInt(t.remaining, 10);
        this.updateColor()
      }
    );
  }

  private updateColor(): void {
    if (this.timer) {
      if (this.timer.state === 'off') {
        this.fgColor = `gray`
      } else if (this.remaining! < 0) {
        this.fgColor = `rgb(255, 0, 0)`
      } else {
        this.progressRatio = (this.timeout! - this.remaining!) / this.timeout!;
        const red = 255 * this.progressRatio
        const green = 255 * (1 - this.progressRatio)
        const blue = 0
        this.fgColor = `rgb(${red},${green},${blue})`
      }
    }
  }
}
