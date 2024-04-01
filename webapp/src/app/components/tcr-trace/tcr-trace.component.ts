import {AfterViewInit, Component, Input, ViewChild} from '@angular/core';
import {NgTerminal, NgTerminalModule} from "ng-terminal";
import {Observable} from "rxjs";

@Component({
  selector: 'app-tcr-trace',
  standalone: true,
  imports: [NgTerminalModule],
  templateUrl: './tcr-trace.component.html',
  styleUrl: './tcr-trace.component.css',
})
export class TcrTraceComponent implements AfterViewInit {
  @Input() text?: Observable<string>;
  @Input() clearTrace?: Observable<void>;
  @ViewChild('term', {static: false}) child!: NgTerminal;

  constructor() {
  }

  ngAfterViewInit(): void {
    this.text?.subscribe((text) => this.print(text));
    this.clearTrace?.subscribe(() => this.clear());
  }

  print(input: string): void {
    // ng-console handles EOL in Windows style, e.g. it needs CRLF to properly
    // go back to beginning of next line in the console
    this.child?.write(toCRLF(input));
  }

  clear() {
    this.child?.underlying?.reset();
  }
}

export function toCRLF(input: string) {
  return input ? (input.replace(/\n/g, "\r\n") + "\r\n") : "";
}
