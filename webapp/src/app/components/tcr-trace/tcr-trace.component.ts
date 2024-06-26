import {AfterViewInit, Component, Input, ViewChild} from '@angular/core';
import {NgTerminal, NgTerminalModule} from "ng-terminal";
import {Observable} from "rxjs";
import {Terminal} from '@xterm/xterm';
import {WebLinksAddon} from '@xterm/addon-web-links';
import {Unicode11Addon} from '@xterm/addon-unicode11';

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
  @ViewChild('term', {static: false}) ngTerminal!: NgTerminal;
  private xterm?: Terminal;

  constructor() {
  }

  ngAfterViewInit(): void {
    this.setupTerminal();

    this.text?.subscribe((text) => this.print(text));
    this.clearTrace?.subscribe(() => this.clear());
  }

  private setupTerminal(): void {
    this.xterm = this.ngTerminal.underlying!;
    this.xterm.loadAddon(new WebLinksAddon());
    this.xterm.loadAddon(new Unicode11Addon());
    this.xterm.unicode.activeVersion = '11';
    this.ngTerminal.setXtermOptions({
      fontFamily: '"Cascadia Code", Menlo, monospace',
      theme: {
        background: '#333333',
        foreground: '#CCCCCC',
        cursor: '#CCCCCC',
      },
      cursorBlink: true
    });
    this.ngTerminal.setRows(20);
    this.ngTerminal.setCols(120);
    this.ngTerminal.setMinWidth(1000);
    this.ngTerminal.setMinHeight(400);
    this.ngTerminal.setDraggable(true);
  }

  print(input: string): void {
    // ng-console handles EOL in Windows style, e.g. it needs CRLF to properly
    // go back to beginning of next line in the console
    this.ngTerminal?.write(toCRLF(input));
  }

  clear(): void {
    this.xterm?.reset();
  }
}

export function toCRLF(input: string): string {
  return input ? (input.replace(/\n/g, "\r\n") + "\r\n") : "";
}
