import {Pipe, PipeTransform} from '@angular/core';

const NOT_SET: string = "[not set]";

@Pipe({
  name: 'showEmpty',
  standalone: true
})
export class ShowEmptyPipe implements PipeTransform {
  transform(value: unknown, ..._args: unknown[]): unknown {
    return value || NOT_SET;
  }

}
