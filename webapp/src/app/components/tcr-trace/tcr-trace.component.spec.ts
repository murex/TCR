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
import { Component, Input, EventEmitter, Output } from "@angular/core";

/**
 * Complete mock of NgTerminal to prevent any xterm initialization
 */
@Component({
  // eslint-disable-next-line @angular-eslint/component-selector
  selector: "ng-terminal",
  template: '<div class="mock-terminal"></div>',
  standalone: true,
})
class MockNgTerminalComponent {
  @Input() setXtermOptions: unknown;
  @Input() setRows: number | undefined;
  @Input() setCols: number | undefined;
  @Input() setMinWidth: string | undefined;
  @Input() setMinHeight: string | undefined;
  @Input() setDraggable: boolean | undefined;
  @Output() keyEventInput = new EventEmitter<unknown>();

  underlying: Record<string, unknown> = {
    reset: jasmine.createSpy("reset"),
    dispose: jasmine.createSpy("dispose"),
    loadAddon: jasmine.createSpy("loadAddon"),
    unicode: { activeVersion: "11" },
  };

  write = jasmine.createSpy("write");
  clear = jasmine.createSpy("clear");
}

describe("TcrTraceComponent", () => {
  let component: TcrTraceComponent;
  let fixture: ComponentFixture<TcrTraceComponent>;
  let mockTerminal: MockNgTerminalComponent;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrTraceComponent, MockNgTerminalComponent],
    })
      .overrideComponent(TcrTraceComponent, {
        set: {
          imports: [MockNgTerminalComponent],
        },
      })
      .compileComponents();

    fixture = TestBed.createComponent(TcrTraceComponent);
    component = fixture.componentInstance;

    // Create and assign mock terminal
    mockTerminal = new MockNgTerminalComponent();
    component.ngTerminal =
      mockTerminal as unknown as typeof component.ngTerminal;

    // Mock the private properties and methods
    (component as unknown as Record<string, unknown>)["xterm"] =
      mockTerminal.underlying;
    spyOn(
      component as unknown as Record<string, jasmine.Spy>,
      "setupTerminal",
    ).and.stub();
  });

  afterEach(() => {
    // Clean up subscriptions
    if ((component as unknown as Record<string, unknown>)["subscriptions"]) {
      (
        component as unknown as {
          subscriptions: Array<{ unsubscribe?: () => void }>;
        }
      ).subscriptions.forEach((sub: { unsubscribe?: () => void }) => {
        if (sub && sub.unsubscribe) {
          sub.unsubscribe();
        }
      });
    }

    // Clean up component
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
      const clearTrace = new Subject<void>();
      component.clearTrace = clearTrace.asObservable();

      // Initialize component
      component.ngAfterViewInit();

      // Trigger the clear
      clearTrace.next();

      // Verify reset was called
      expect(mockTerminal.underlying["reset"]).toHaveBeenCalled();

      // Clean up
      clearTrace.complete();
    });

    it("should print text upon reception of text observable", () => {
      const text = new Subject<string>();
      const input = "Hello World";
      component.text = text.asObservable();

      // Initialize component
      component.ngAfterViewInit();

      // Send the text
      text.next(input);

      // Verify write was called
      expect(mockTerminal.write).toHaveBeenCalledWith(input + "\r\n");

      // Clean up
      text.complete();
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
      const input = "Hello World";
      component.print(input);
      expect(mockTerminal.write).toHaveBeenCalledWith(input + "\r\n");
    });
  });

  describe("clear function", () => {
    it("should clear the terminal contents", () => {
      component.clear();
      expect(mockTerminal.underlying["reset"]).toHaveBeenCalled();
    });
  });

  describe("lifecycle hooks", () => {
    it("should clean up subscriptions on destroy", () => {
      const text = new Subject<string>();
      const clearTrace = new Subject<void>();
      component.text = text.asObservable();
      component.clearTrace = clearTrace.asObservable();

      component.ngAfterViewInit();

      // Verify subscriptions are created
      expect(
        (component as unknown as { subscriptions: unknown[] })["subscriptions"]
          .length,
      ).toBe(2);

      // Destroy component
      component.ngOnDestroy();

      // Verify dispose was called
      expect(mockTerminal.underlying["dispose"]).toHaveBeenCalled();

      // Clean up subjects
      text.complete();
      clearTrace.complete();
    });

    it("should handle undefined observables gracefully", () => {
      component.text = undefined;
      component.clearTrace = undefined;

      expect(() => component.ngAfterViewInit()).not.toThrow();
      expect(
        (component as unknown as { subscriptions: unknown[] })["subscriptions"]
          .length,
      ).toBe(0);
    });
  });
});
