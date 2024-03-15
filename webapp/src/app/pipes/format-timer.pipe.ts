import {Pipe, PipeTransform} from '@angular/core';

const SEPARATOR = ":";
const SECONDS_PER_MINUTE = 60;
const SECONDS_PER_HOUR = 3600;

@Pipe({
  name: 'formatTimer',
  standalone: true
})
export class FormatTimerPipe implements PipeTransform {
  transform(value: unknown, ..._args: unknown[]): unknown {
    const duration = parseInt(value as string, 10);
    if (isNaN(duration)) {
      return `--${SEPARATOR}--`;
    }
    const sign = duration < 0 ? "-" : "";
    const totalMinutes = Math.abs(duration);
    const nbHours = Math.floor(totalMinutes / SECONDS_PER_HOUR);
    const nbMinutes = Math.floor((totalMinutes % SECONDS_PER_HOUR) / SECONDS_PER_MINUTE);
    const nbSeconds = totalMinutes % SECONDS_PER_MINUTE;
    const hours = nbHours > 0 ? nbHours.toString() + SEPARATOR : "";
    const minutes = nbMinutes.toString().padStart(2, "0") + SEPARATOR;
    const seconds = nbSeconds.toString().padStart(2, "0");
    return `${sign}${hours}${minutes}${seconds}`;
  }

}
