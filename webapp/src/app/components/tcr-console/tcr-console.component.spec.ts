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
import {
  formatRoleMessage,
  getRoleAction,
  getRoleName,
  isRoleStartMessage,
  TcrConsoleComponent,
} from "./tcr-console.component";
import { TcrMessage, TcrMessageType } from "../../interfaces/tcr-message";
import { TcrMessageService } from "../../services/tcr-message.service";
import { TcrControlsService } from "../../services/tcr-controls.service";
import { TcrRolesComponent } from "../tcr-roles/tcr-roles.component";
import { TcrTraceComponent } from "../tcr-trace/tcr-trace.component";
import { TcrControlsComponent } from "../tcr-controls/tcr-controls.component";
import {
  bgDarkGray,
  cyan,
  green,
  lightCyan,
  lightYellow,
  red,
  yellow,
} from "ansicolor";
import { Component, Input } from "@angular/core";
import { Observable, Subject, of } from "rxjs";
import {
  MockNgTerminalComponent,
  setupTerminalTestEnvironment,
} from "../../../testing/terminal-test-utils";

// Mock components for testing
@Component({
  selector: "app-tcr-roles",
  template: '<div class="mock-roles"></div>',
  standalone: true,
})
class MockTcrRolesComponent {}

@Component({
  selector: "app-tcr-trace",
  template: '<div class="mock-trace"></div>',
  standalone: true,
})
class MockTcrTraceComponent {
  @Input() text?: Observable<string>;
  @Input() clearTrace?: Observable<void>;
}

@Component({
  selector: "app-tcr-controls",
  template: '<div class="mock-controls"></div>',
  standalone: true,
})
class MockTcrControlsComponent {}

class FakeTcrMessageService {
  message$ = new Subject<TcrMessage>();
}

class FakeTcrControlsService {
  abortCommand() {
    return of({});
  }
}

describe("TcrConsoleComponent", () => {
  let component: TcrConsoleComponent;
  let fixture: ComponentFixture<TcrConsoleComponent>;

  beforeAll(() => {
    // Set up safe terminal test environment
    setupTerminalTestEnvironment();
  });

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrConsoleComponent],
      declarations: [MockNgTerminalComponent],
      providers: [
        { provide: TcrMessageService, useClass: FakeTcrMessageService },
        { provide: TcrControlsService, useClass: FakeTcrControlsService },
      ],
    })
      .overrideComponent(TcrConsoleComponent, {
        remove: {
          imports: [TcrRolesComponent, TcrTraceComponent, TcrControlsComponent],
        },
        add: {
          imports: [
            MockTcrRolesComponent,
            MockTcrTraceComponent,
            MockTcrControlsComponent,
          ],
        },
      })
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TcrConsoleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  afterEach(() => {
    // Clean up any terminal instances to prevent dimension errors
    if (fixture) {
      fixture.destroy();
    }
  });

  describe("component instance", () => {
    it("should be created", () => {
      expect(component).toBeTruthy();
    });
  });

  describe("component DOM", () => {
    it(`should have a title`, () => {
      const element = fixture.nativeElement.querySelector("h1");
      expect(element).toBeTruthy();
      expect(element.textContent).toContain("TCR Console");
    });

    [
      {
        description: "a roles component",
        selector: "app-tcr-roles",
      },
      {
        description: "a trace component",
        selector: "app-tcr-trace",
      },
    ].forEach((testCase) => {
      it(`should contain ${testCase.description}`, () => {
        const element = fixture.nativeElement.querySelector(testCase.selector);
        expect(element).toBeTruthy();
      });
    });
  });

  describe("getRoleAction function", () => {
    it("should extract the action from the message", () => {
      const message = "driver:start";
      const result = getRoleAction(message);
      expect(result).toEqual("start");
    });

    it("should return an empty string if the message is empty", () => {
      const message = "";
      const result = getRoleAction(message!);
      expect(result).toEqual("");
    });

    it("should return an empty string if the message is undefined", () => {
      const message = undefined;
      const result = getRoleAction(message!);
      expect(result).toEqual("");
    });
  });

  describe("getRoleName function", () => {
    it("should extract the role name from the message", () => {
      const message = "driver:start";
      const result = getRoleName(message);
      expect(result).toEqual("driver");
    });

    it("should return an empty string if the message is empty", () => {
      const message = "";
      const result = getRoleName(message!);
      expect(result).toEqual("");
    });

    it("should return an empty string if the message is undefined", () => {
      const message = undefined;
      const result = getRoleName(message!);
      expect(result).toEqual("");
    });
  });

  describe("formatRoleMessage function", () => {
    [
      { message: "driver:start", expected: "Starting driver role" },
      { message: "driver:end", expected: "Ending driver role" },
      { message: "navigator:start", expected: "Starting navigator role" },
      { message: "navigator:end", expected: "Ending navigator role" },
      { message: "", expected: "" },
      { message: undefined, expected: "" },
    ].forEach((testCase) => {
      it(`should format the message '${testCase.message}'`, () => {
        const result = formatRoleMessage(testCase.message!);
        expect(result).toEqual(testCase.expected);
      });
    });
  });

  describe("isRoleStartMessage function", () => {
    [
      { type: TcrMessageType.ROLE, message: "driver:start", expected: true },
      { type: TcrMessageType.ROLE, message: "driver:end", expected: false },
      { type: TcrMessageType.ROLE, message: "navigator:start", expected: true },
      { type: TcrMessageType.ROLE, message: "navigator:end", expected: false },
      { type: TcrMessageType.ROLE, message: "other", expected: false },
      { type: TcrMessageType.INFO, message: "other", expected: false },
    ].forEach((testCase) => {
      const expectation = `should return ${testCase.expected} for '${testCase.type}:${testCase.message}' messages`;
      it(expectation, () => {
        const result = isRoleStartMessage({
          type: testCase.type,
          text: testCase.message,
        } as TcrMessage);
        expect(result).toEqual(testCase.expected);
      });
    });
  });

  describe("printSimple function", () => {
    it("should print a simple text", () => {
      const text = "some simple text";

      let actual: string | undefined;
      component.text.subscribe((value) => {
        actual = value;
      });

      component.printSimple(text);
      expect(actual).toEqual(text);
    });
  });

  describe("printInfo function", () => {
    it("should format and print an info text", () => {
      const text = "some info text";

      let actual: string | undefined;
      component.text.subscribe((value) => {
        actual = value;
      });

      component.printInfo(text);
      expect(actual).toEqual(cyan(text));
    });
  });

  describe("printTitle function", () => {
    it("should format and print a title text", () => {
      const text = "some title text";

      let actual: string | undefined;
      component.text.subscribe((value) => {
        actual = value;
      });

      component.printTitle(text);
      const lineSep = lightCyan("─".repeat(80));
      expect(actual).toEqual(lineSep + "\n" + lightCyan(text));
    });
  });

  describe("printRole function", () => {
    it("should format and print a role text", () => {
      const text = "driver:start";

      let actual: string | undefined;
      component.text.subscribe((value) => {
        actual = value;
      });

      component.printRole(text);
      const lineSep = yellow("─".repeat(80));
      expect(actual).toEqual(
        lineSep + "\n" + lightYellow("Starting driver role") + "\n" + lineSep,
      );
    });
  });

  describe("printSuccess function", () => {
    it("should format and print a success text", () => {
      const text = "some success text";

      let actual: string | undefined;
      component.text.subscribe((value) => {
        actual = value;
      });

      component.printSuccess(text);
      expect(actual).toEqual("🟢 " + green(text));
    });
  });

  describe("printWarning function", () => {
    it("should format and print a warning text", () => {
      const text = "some warning text";

      let actual: string | undefined;
      component.text.subscribe((value) => {
        actual = value;
      });

      component.printWarning(text);
      expect(actual).toEqual("🔶 " + yellow(text));
    });
  });

  describe("printError function", () => {
    it("should format and print an error text", () => {
      const text = "some error text";

      let actual: string | undefined;
      component.text.subscribe((value) => {
        actual = value;
      });

      component.printError(text);
      expect(actual).toEqual("🟥 " + red(text));
    });
  });

  describe("printUnhandled function", () => {
    it("should format and print an unhandled message", () => {
      const text = "some unhandled text";
      const type = "xxx";

      let actual: string | undefined;
      component.text.subscribe((value) => {
        actual = value;
      });

      component.printUnhandled(type, text);
      expect(actual).toEqual(bgDarkGray("[xxx]") + " " + text);
    });
  });

  describe("clear function", () => {
    it("should update clearTrace subject", () => {
      let actual: boolean = false;
      component.clearTrace.subscribe(() => {
        actual = true;
      });

      component.clear();
      expect(actual).toBeTruthy();
    });
  });

  // test cases for printMessage function
  describe("printMessage function", () => {
    [
      {
        type: "simple",
        text: "some simple text",
        expectedFunction: "printSimple",
      },
      {
        type: "info",
        text: "some info text",
        expectedFunction: "printInfo",
      },
      {
        type: "title",
        text: "some title text",
        expectedFunction: "printTitle",
      },
      {
        type: "role",
        text: "driver:start",
        expectedFunction: "printRole",
      },
      {
        type: "success",
        text: "some success text",
        expectedFunction: "printSuccess",
      },
      {
        type: "warning",
        text: "some warning text",
        expectedFunction: "printWarning",
      },
      {
        type: "error",
        text: "some error text",
        expectedFunction: "printError",
      },
      {
        type: "unhandled",
        text: "some unhandled text",
        expectedFunction: "printUnhandled",
      },
    ].forEach((testCase) => {
      it(`should format and print ${testCase.type} messages`, () => {
        const printFunction = spyOn(
          component as any, // eslint-disable-line @typescript-eslint/no-explicit-any
          testCase.expectedFunction,
        ).and.callThrough();

        component.printMessage({
          type: testCase.type,
          text: testCase.text,
        } as TcrMessage);

        expect(printFunction).toHaveBeenCalled();
      });
    });

    it("should ignore timer messages", () => {
      const printFunction = spyOn(component, "print").and.callThrough();

      component.printMessage({
        type: "timer",
        text: "some timer text",
      } as TcrMessage);

      expect(printFunction).not.toHaveBeenCalled();
    });
  });
});
