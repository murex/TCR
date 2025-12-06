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
import { configureComponentTestingModule } from "../../../test-helpers/angular-test-helpers";
import { injectService } from "../../../test-helpers/angular-test-helpers";

import { TcrTimerComponent } from "./tcr-timer.component";
import { Observable, of } from "rxjs";
import { TcrMessage, TcrMessageType } from "../../interfaces/tcr-message";
import { TcrTimerService } from "../../services/tcr-timer.service";
import { TcrTimer, TcrTimerState } from "../../interfaces/tcr-timer";
import { By } from "@angular/platform-browser";
import { FaIconLibrary } from "@fortawesome/angular-fontawesome";
import { registerFontAwesomeIcons } from "../../shared/font-awesome-icons";

class FakeTcrTimerService {
  message$ = new Observable<TcrMessage>();

  getTimer(): Observable<TcrTimer> {
    return of({
      state: TcrTimerState.OFF,
      timeout: "0",
      elapsed: "0",
      remaining: "0",
    });
  }
}

describe("TcrTimerComponent", () => {
  let component: TcrTimerComponent;
  let fixture: ComponentFixture<TcrTimerComponent>;
  let serviceFake: TcrTimerService;

  beforeEach(async () => {
    await configureComponentTestingModule(
      TcrTimerComponent,
      [],
      [
        { provide: TcrTimerService, useClass: FakeTcrTimerService },
        FaIconLibrary,
      ],
    );

    // Register FontAwesome icons
    const library = injectService(FaIconLibrary);
    registerFontAwesomeIcons(library);
  });

  beforeEach(() => {
    serviceFake = injectService(TcrTimerService);
    fixture = TestBed.createComponent(TcrTimerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  afterEach(() => {
    // Clean up the interval created in ngAfterViewInit
    if (component.ngOnDestroy) {
      component.ngOnDestroy();
    }
    fixture.destroy();
  });

  describe("component instance", () => {
    it("should be created", () => {
      expect(component).toBeTruthy();
    });
  });

  describe("component initialization", () => {
    const clockIcon = "clock";
    const warningIcon = "circle-exclamation";

    [
      {
        state: TcrTimerState.OFF,
        timeout: "100",
        elapsed: "0",
        remaining: "0",
        expectedColor: "rgb(128, 128, 128)",
        expectedIcon: clockIcon,
        expectedText: "00:00",
      },
      {
        state: TcrTimerState.PENDING,
        timeout: "100",
        elapsed: "0",
        remaining: "100",
        expectedColor: "rgb(255, 255, 255)",
        expectedIcon: clockIcon,
        expectedText: "01:40",
      },
      {
        state: TcrTimerState.RUNNING,
        timeout: "100",
        elapsed: "20",
        remaining: "80",
        expectedColor: "rgb(255, 204, 204)",
        expectedIcon: clockIcon,
        expectedText: "01:20",
      },
      {
        state: TcrTimerState.STOPPED,
        timeout: "100",
        elapsed: "60",
        remaining: "0",
        expectedColor: "rgb(128, 128, 128)",
        expectedIcon: clockIcon,
        expectedText: "00:00",
      },
      {
        state: TcrTimerState.TIMEOUT,
        timeout: "100",
        elapsed: "120",
        remaining: "-20",
        expectedColor: "rgb(255, 0, 0)",
        expectedIcon: warningIcon,
        expectedText: "-00:20",
      },
    ].forEach((testCase) => {
      it(`should work with timer in ${testCase.state} state`, () => {
        const timer: TcrTimer = {
          state: testCase.state,
          timeout: testCase.timeout,
          elapsed: testCase.elapsed,
          remaining: testCase.remaining,
        };

        // Clean up existing component first
        if (component.ngOnDestroy) {
          component.ngOnDestroy();
        }
        fixture.destroy();

        // Have the service fake's getTimer method return the timer data
        serviceFake.getTimer = () => of(timer);
        fixture = TestBed.createComponent(TcrTimerComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();

        // Verify that the component's timer attribute is set correctly
        expect(component.timer).toEqual(timer);

        // Verify that the component is rendered with the expected color
        const componentElement = fixture.debugElement.query(
          By.css(`[data-testid="timer-component"]`),
        );
        expect(componentElement).toBeTruthy();
        expect(componentElement.nativeElement.style.color).toEqual(
          testCase.expectedColor,
        );

        // Verify that the right icon is rendered
        const iconElement = fixture.debugElement.query(
          By.css(`[data-testid="timer-icon"]`),
        );
        expect(iconElement).toBeTruthy();
        // Check for fa-icon element with correct icon
        expect(iconElement.nativeElement.tagName.toLowerCase()).toBe("fa-icon");
        // The icon is rendered as an SVG, check the data-icon attribute
        const svgElement = iconElement.nativeElement.querySelector("svg");
        if (svgElement) {
          const dataIcon = svgElement.getAttribute("data-icon");
          expect(dataIcon).toBe(testCase.expectedIcon);
        } else {
          // Fallback: check if the icon name is in the class
          expect(iconElement.nativeElement.innerHTML).toContain(
            testCase.expectedIcon,
          );
        }

        // Verify that the timer text is rendered
        const textElement = fixture.debugElement.query(
          By.css(`[data-testid="timer-label"]`),
        );
        expect(textElement).toBeTruthy();
        expect(textElement.nativeElement.textContent).toEqual(
          testCase.expectedText,
        );
      });
    });
  });

  describe("component updateColor", () => {
    [
      {
        state: TcrTimerState.OFF,
        timeout: 100,
        elapsed: 0,
        remaining: 0,
        expectedColor: "rgb(128,128,128)",
      },
      {
        state: TcrTimerState.PENDING,
        timeout: 100,
        elapsed: 0,
        remaining: 100,
        expectedColor: "rgb(255,255,255)",
      },
      {
        state: TcrTimerState.RUNNING,
        timeout: 100,
        elapsed: 20,
        remaining: 80,
        expectedColor: "rgb(255,204,204)",
      },
      {
        state: TcrTimerState.STOPPED,
        timeout: 100,
        elapsed: 60,
        remaining: 0,
        expectedColor: "rgb(128,128,128)",
      },
      {
        state: TcrTimerState.TIMEOUT,
        timeout: 100,
        elapsed: 120,
        remaining: -20,
        expectedColor: "rgb(255,0,0)",
      },
    ].forEach((testCase) => {
      const input = `${testCase.state}/${testCase.timeout}/${testCase.elapsed}/${testCase.remaining}`;
      it(`should translate ${input} into ${testCase.expectedColor}`, () => {
        // Setup the timer component with the timer data
        component.timer = {
          state: testCase.state,
          timeout: `${testCase.timeout}`,
          elapsed: `${testCase.elapsed}`,
          remaining: `${testCase.remaining}`,
        };
        component.timeout = testCase.timeout;
        component.remaining = testCase.remaining;

        component.updateColor();
        expect(component.fgColor).toEqual(testCase.expectedColor);
      });
    });
  });

  describe("component periodicUpdate", () => {
    [
      {
        state: TcrTimerState.OFF,
        timeout: 100,
        elapsed: 0,
        remaining: 0,
        expectedRemaining: 0,
      },
      {
        state: TcrTimerState.PENDING,
        timeout: 100,
        elapsed: 0,
        remaining: 100,
        expectedRemaining: 100,
      },
      {
        state: TcrTimerState.RUNNING,
        timeout: 100,
        elapsed: 20,
        remaining: 80,
        expectedRemaining: 79,
      },
      {
        state: TcrTimerState.STOPPED,
        timeout: 100,
        elapsed: 60,
        remaining: 0,
        expectedRemaining: 0,
      },
      {
        state: TcrTimerState.TIMEOUT,
        timeout: 100,
        elapsed: 120,
        remaining: -20,
        expectedRemaining: -21,
      },
    ].forEach((testCase) => {
      it(`should change remaining time from ${testCase.remaining} to ${testCase.expectedRemaining} when ${testCase.state}`, () => {
        // Setup the timer component with the timer data
        component.timer = {
          state: testCase.state,
          timeout: `${testCase.timeout}`,
          elapsed: `${testCase.elapsed}`,
          remaining: `${testCase.remaining}`,
        };
        component.timeout = testCase.timeout;
        component.remaining = testCase.remaining;

        component.periodicUpdate();
        expect(component.remaining).toEqual(testCase.expectedRemaining);
      });
    });

    [
      {
        description: "ticking too fast",
        state: TcrTimerState.RUNNING,
        timeout: 100,
        initialRemaining: 80,
        serverRemaining: 71,
      },
      {
        description: "ticking too slow",
        state: TcrTimerState.RUNNING,
        timeout: 100,
        initialRemaining: 80,
        serverRemaining: 69,
      },
    ].forEach((testCase) => {
      it(`should re-sync with the server when ${testCase.description}`, () => {
        // Have the service fake's getTimer method return the server timer data
        serviceFake.getTimer = () =>
          of({
            state: TcrTimerState.RUNNING,
            timeout: `${testCase.timeout}`,
            elapsed: `${testCase.timeout - testCase.serverRemaining}`,
            remaining: `${testCase.serverRemaining}`,
          });

        // Setup the timer component starting state
        component.timer = {
          state: TcrTimerState.RUNNING,
          timeout: `${testCase.timeout}`,
          elapsed: `${testCase.timeout - testCase.initialRemaining}`,
          remaining: `${testCase.initialRemaining}`,
        };
        component.timeout = testCase.timeout;
        component.remaining = testCase.initialRemaining;

        // Simulate 10 seconds passing
        for (let tick = 1; tick <= 10; tick++) {
          component.periodicUpdate();
          expect(component.remaining).toEqual(testCase.initialRemaining - tick);
        }

        // Verify that the timer has re-synced with the server
        component.periodicUpdate();
        expect(component.remaining).toEqual(testCase.serverRemaining);
      });
    });
  });

  describe("component refresh", () => {
    [
      {
        expectation: "should fetch timer data on actual messages",
        timerBefore: {
          state: TcrTimerState.OFF,
          timeout: "0",
          elapsed: "0",
          remaining: "0",
        } as TcrTimer,
        message: {
          type: TcrMessageType.TIMER,
          text: "start:100:0:100",
        } as TcrMessage,
        timerAfter: {
          state: TcrTimerState.RUNNING,
          timeout: "100",
          elapsed: "5",
          remaining: "95",
        } as TcrTimer,
      },
      {
        expectation: "should not fetch timer data on empty messages",
        timerBefore: {
          state: TcrTimerState.OFF,
          timeout: "0",
          elapsed: "0",
          remaining: "0",
        } as TcrTimer,
        message: undefined,
        timerAfter: {
          state: TcrTimerState.RUNNING,
          timeout: "100",
          elapsed: "5",
          remaining: "95",
        } as TcrTimer,
      },
    ].forEach((testCase) => {
      it(`${testCase.expectation}`, () => {
        // Have the service fake's getTimer method return the starting timer data
        serviceFake.getTimer = () => of(testCase.timerBefore);

        // Initialize without triggering change detection
        component.ngOnInit();
        expect(component.timer).toEqual(testCase.timerBefore);

        // Update the service fake to return the expected new timer data
        serviceFake.getTimer = () => of(testCase.timerAfter);

        // Trigger the refresh method with a message
        component.refresh(testCase.message!);

        // Verify that the component's timer attribute was updated
        expect(component.timer).toEqual(
          testCase.message ? testCase.timerAfter : testCase.timerBefore,
        );
      });
    });
  });
});
