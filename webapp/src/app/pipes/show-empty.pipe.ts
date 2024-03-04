import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'showEmpty',
  standalone: true
})
export class ShowEmptyPipe implements PipeTransform {

  transform(value: unknown, ..._args: unknown[]): unknown {
    return value || "[not set]";
  }

}
