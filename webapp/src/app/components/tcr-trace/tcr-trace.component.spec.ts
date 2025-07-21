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

import { ComponentFixture, TestBed } from "@angular/core/testing";
import { TcrTraceComponent, toCRLF } from "./tcr-trace.component";
import { Subject } from "rxjs";
import { Component } from "@angular/core";
import { NgTerminal } from "ng-terminal";

@Component({
  selector: "ng-terminal", // eslint-disable-line @angular-eslint/component-selector
  template: "",
  standalone: false, // eslint-disable-line @angular-eslint/prefer-standalone
})
// eslint-disable-next-line @angular-eslint/component-class-suffix
class StubNgTerminal {
  underlying: any; // eslint-disable-line @typescript-eslint/no-explicit-any

  constructor() {
    this.underlying = {
      reset: jasmine.createSpy("reset"),
      dispose: jasmine.createSpy("dispose"),
      loadAddon: jasmine.createSpy("loadAddon"),
      unicode: { activeVersion: "11" },
      // Mock viewport with dimensions to prevent the error
      _core: {
        viewport: {
          dimensions: { cols: 80, rows: 24 },
          syncScrollArea: jasmine.createSpy("syncScrollArea"),
          _innerRefresh: jasmine.createSpy("_innerRefresh"),
        },
      },
    };
  }

  write(_data: string): void {}

  setXtermOptions(_options: unknown): void {}

  setRows(_rows: number): void {}

  setCols(_cols: number): void {}

  setMinWidth(_width: number): void {}

  setMinHeight(_height: number): void {}

  setDraggable(_draggable: boolean): void {}
}

describe("TcrTraceComponent", () => {
  let component: TcrTraceComponent;
  let fixture: ComponentFixture<TcrTraceComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrTraceComponent],
      declarations: [StubNgTerminal],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TcrTraceComponent);
    component = fixture.componentInstance;
    component.ngTerminal = TestBed.createComponent(StubNgTerminal)
      .componentInstance as NgTerminal;
    fixture.detectChanges();
  });

  afterEach(() => {
    // Clean up terminal instances to prevent memory leaks and dimension errors
    if (component.ngTerminal?.underlying?.dispose) {
      component.ngTerminal.underlying.dispose();
    }
    if (fixture) {
      fixture.destroy();
    }
  });

  describe("component instance", () => {
    it("should be created", () => {
      expect(component).toBeTruthy();
    });

    it("should have ng-terminal child component", () => {
      expect(component.ngTerminal).toBeTruthy();
    });

    it("should have ng-terminal child.underlying component", () => {
      expect(component.ngTerminal.underlying).toBeTruthy();
    });

    it("should clear the terminal upon reception of clearTrace observable", () => {
      let cleared = false;
      component.ngTerminal.underlying!.reset = () => {
        cleared = true;
      };

      const clearTrace = new Subject<void>();
      component.clearTrace = clearTrace.asObservable();

      component.ngAfterViewInit();
      clearTrace.next();

      expect(cleared).toBeTruthy();
    });

    it("should print text upon reception of text observable", () => {
      let written = "";
      component.ngTerminal.write = (input: string) => {
        written = input;
      };
      const text = new Subject<string>();

      const input = "Hello World";
      component.text = text.asObservable();

      component.ngAfterViewInit();
      text.next(input);

      expect(written).toEqual(input + "\r\n");
    });
  });

  describe("toCRLF function", () => {
    it("should replace all LF with CRLF in the input string", () => {
      const input = "Hello\nWorld\n";
      const result = toCRLF(input);
      expect(result).toEqual("Hello\r\nWorld\r\n\r\n");
    });

    it("should append CRLF to the input string if it does not end with LF", () => {
      const input = "Hello World";
      const result = toCRLF(input);
      expect(result).toEqual("Hello World\r\n");
    });

    it("should return an empty string if the input string is empty", () => {
      const input = "";
      const result = toCRLF(input);
      expect(result).toEqual("");
    });

    it("should return an empty string if the input string is undefined", () => {
      const input = undefined;
      const result = toCRLF(input!);
      expect(result).toEqual("");
    });
  });

  describe("print function", () => {
    it("should send text to the terminal", () => {
      let written = "";
      component.ngTerminal.write = (input: string) => {
        written = input;
      };
      const input = "Hello World";
      component.print(input);
      expect(written).toEqual(input + "\r\n");
    });
  });

  describe("clear function", () => {
    it("should clear the terminal contents", () => {
      let cleared = false;
      component.ngTerminal.underlying!.reset = () => {
        cleared = true;
      };
      component.clear();
      expect(cleared).toBeTruthy();
    });
  });
});
