import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'showEmpty',
  standalone: true
})
export class ShowEmptyPipe implements PipeTransform {

  private readonly NOT_SET = "[not set]";

  transform(value: unknown, ..._args: unknown[]): unknown {
    return value || this.NOT_SET;
  }

}
