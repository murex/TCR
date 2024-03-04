import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'onOff',
  standalone: true
})
export class OnOffPipe implements PipeTransform {

  transform(value: unknown, ..._args: unknown[]): unknown {
    return value ? "✅" : "❌";
  }

}
