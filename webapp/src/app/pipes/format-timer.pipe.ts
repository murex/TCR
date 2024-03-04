import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'formatTimer',
  standalone: true
})
export class FormatTimerPipe implements PipeTransform {

  transform(value: unknown, ..._args: unknown[]): unknown {
    const duration = parseInt(value as string, 10);
    if (isNaN(duration)) {
      return "--:--";
    }
    const sign = duration < 0 ? "-" : "";
    const minutes = Math.floor(Math.abs(duration) / 60).toString().padStart(2, "0");
    const seconds = (Math.abs(duration) % 60).toString().padStart(2, "0");
    return `${sign}${minutes}:${seconds}`;
  }

}
