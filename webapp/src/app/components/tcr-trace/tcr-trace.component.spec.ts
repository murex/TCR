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

import {
  MockNgTerminalComponent,
  createMockNgTerminal,
  setupTerminalTestEnvironment,
} from "../../../testing/terminal-test-utils";

describe("TcrTraceComponent", () => {
  let component: TcrTraceComponent;
  let fixture: ComponentFixture<TcrTraceComponent>;
  let mockTerminal: ReturnType<typeof createMockNgTerminal>;

  beforeAll(() => {
    // Set up safe terminal test environment
    setupTerminalTestEnvironment();
  });

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrTraceComponent],
      declarations: [MockNgTerminalComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TcrTraceComponent);
    component = fixture.componentInstance;

    // Create a mock terminal that won't cause afterAll errors
    mockTerminal = createMockNgTerminal();

    // Ensure the mock has a proper write method
    if (!mockTerminal.write) {
      mockTerminal.write = jasmine.createSpy("write");
    }

    component.ngTerminal = mockTerminal;

    // Override the component's private xterm property to prevent real terminal setup
    (component as unknown as { xterm: unknown }).xterm =
      mockTerminal.underlying;

    // Mock the setupTerminal method to prevent real terminal initialization
    spyOn(
      component as unknown as { setupTerminal: () => void },
      "setupTerminal",
    ).and.callFake(() => {
      (component as unknown as { xterm: unknown }).xterm =
        mockTerminal.underlying;
    });

    fixture.detectChanges();
  });

  afterEach(() => {
    // Clean up subscriptions first
    if (component && component.ngOnDestroy) {
      component.ngOnDestroy();
    }

    // Clean up mock terminal - prevent any async operations
    if (mockTerminal && mockTerminal.underlying) {
      // Stop any viewport operations
      const underlying = mockTerminal.underlying as {
        _core?: {
          viewport?: {
            _refreshAnimationFrame?: number | null;
            dimensions?: unknown;
            _coreBrowserService?: unknown;
          };
        };
        dispose?: () => void;
      };
      const viewport = underlying._core?.viewport;
      if (viewport) {
        if (viewport._refreshAnimationFrame) {
          cancelAnimationFrame(viewport._refreshAnimationFrame);
          viewport._refreshAnimationFrame = null;
        }
        viewport.dimensions = undefined;
        viewport._coreBrowserService = null;
      }

      // Dispose the terminal
      if (typeof underlying.dispose === "function") {
        try {
          underlying.dispose();
        } catch (_e) {
          // Ignore disposal errors in tests
        }
      }
    }

    // Clean up fixture
    if (fixture) {
      fixture.destroy();
    }
  });

  afterAll(() => {
    // Additional cleanup to prevent afterAll errors
    const globalWithCanvas = globalThis as {
      HTMLCanvasElement?: {
        prototype: { getContext?: unknown };
      };
    };
    if (typeof globalWithCanvas.HTMLCanvasElement !== "undefined") {
      delete globalWithCanvas.HTMLCanvasElement.prototype.getContext;
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

      // Use the mock terminal's reset method
      if (mockTerminal.underlying) {
        mockTerminal.underlying.reset = jasmine
          .createSpy("reset")
          .and.callFake(() => {
            cleared = true;
          });
      }

      const clearTrace = new Subject<void>();
      component.clearTrace = clearTrace.asObservable();

      component.ngAfterViewInit();
      clearTrace.next();

      expect(cleared).toBeTruthy();

      // Clean up the subject
      clearTrace.complete();
    });

    it("should print text upon reception of text observable", () => {
      let written = "";

      // Set up the write spy on the component's ngTerminal
      component.ngTerminal.write = jasmine
        .createSpy("write")
        .and.callFake((input: string) => {
          written = input;
        });

      const text = new Subject<string>();

      const input = "Hello World";
      component.text = text.asObservable();

      component.ngAfterViewInit();
      text.next(input);

      expect(written).toEqual(input + "\r\n");

      // Clean up the subject
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
      let written = "";

      // Set up the write spy on the component's ngTerminal
      component.ngTerminal.write = jasmine
        .createSpy("write")
        .and.callFake((input: string) => {
          written = input;
        });

      const input = "Hello World";
      component.print(input);
      expect(written).toEqual(input + "\r\n");
    });
  });

  describe("clear function", () => {
    it("should clear the terminal contents", () => {
      let cleared = false;
      if (mockTerminal.underlying) {
        mockTerminal.underlying.reset = jasmine
          .createSpy("reset")
          .and.callFake(() => {
            cleared = true;
          });
      }
      component.clear();
      expect(cleared).toBeTruthy();
    });
  });
});
