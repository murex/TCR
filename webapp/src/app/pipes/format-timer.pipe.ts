/*
Copyright (c) 2024 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
    const duration: number = parseInt(value as string, 10);
    if (isNaN(duration)) {
      return `--${SEPARATOR}--`;
    }
    const sign: string = duration < 0 ? "-" : "";
    const totalMinutes: number = Math.abs(duration);
    const nbHours: number = Math.floor(totalMinutes / SECONDS_PER_HOUR);
    const nbMinutes: number = Math.floor((totalMinutes % SECONDS_PER_HOUR) / SECONDS_PER_MINUTE);
    const nbSeconds: number = totalMinutes % SECONDS_PER_MINUTE;
    const hours: string = nbHours > 0 ? nbHours.toString() + SEPARATOR : "";
    const minutes: string = nbMinutes.toString().padStart(2, "0") + SEPARATOR;
    const seconds: string = nbSeconds.toString().padStart(2, "0");
    return `${sign}${hours}${minutes}${seconds}`;
  }

}
