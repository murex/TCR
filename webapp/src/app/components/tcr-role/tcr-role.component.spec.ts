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

import {
  ComponentFixture,
  TestBed,
  fakeAsync,
  tick,
} from "@angular/core/testing";
import { configureComponentTestingModule } from "../../../test-helpers/angular-test-helpers";
import { injectService } from "../../../test-helpers/angular-test-helpers";
import { TcrRoleComponent } from "./tcr-role.component";
import { Observable, of } from "rxjs";
import { TcrMessage, TcrMessageType } from "../../interfaces/tcr-message";
import { TcrRolesService } from "../../services/trc-roles.service";
import { TcrRole } from "../../interfaces/tcr-role";
import { By } from "@angular/platform-browser";
import { ChangeDetectorRef } from "@angular/core";
import { FaIconLibrary } from "@fortawesome/angular-fontawesome";
import { registerFontAwesomeIcons } from "../../shared/font-awesome-icons";

class FakeTcrRolesService {
  message$ = new Observable<TcrMessage>();
  private roleStates: Map<string, boolean> = new Map();

  getRole(name: string): Observable<TcrRole> {
    const active = this.roleStates.get(name) || false;
    return of({ name: name, description: name + " role", active: active });
  }

  activateRole(name: string, state: boolean): Observable<TcrRole> {
    this.roleStates.set(name, state);
    return of({ name: name, description: name + " role", active: state });
  }

  setRoleState(name: string, state: boolean): void {
    this.roleStates.set(name, state);
  }
}

describe("TcrRoleComponent", () => {
  let component: TcrRoleComponent;
  let fixture: ComponentFixture<TcrRoleComponent>;
  let serviceFake: FakeTcrRolesService;
  let changeDetectorRef: ChangeDetectorRef;

  beforeEach(async () => {
    await configureComponentTestingModule(
      TcrRoleComponent,
      [],
      [
        { provide: TcrRolesService, useClass: FakeTcrRolesService },
        FaIconLibrary,
      ],
    );

    // Register FontAwesome icons
    const library = injectService(FaIconLibrary);
    registerFontAwesomeIcons(library);
  });

  beforeEach(() => {
    serviceFake = TestBed.inject(
      TcrRolesService,
    ) as unknown as FakeTcrRolesService;
    fixture = TestBed.createComponent(TcrRoleComponent);
    component = fixture.componentInstance;
    changeDetectorRef = fixture.changeDetectorRef;
  });

  describe("component instance", () => {
    it("should be created", () => {
      expect(component).toBeTruthy();
    });
  });

  describe("component initialization", () => {
    [
      {
        name: "driver",
        description: "driver role",
        active: true,
        componentClass: "driver-on",
        iconClass: "keyboard",
      },
      {
        name: "driver",
        description: "driver role",
        active: false,
        componentClass: "role-off",
        iconClass: "keyboard",
      },
      {
        name: "navigator",
        description: "navigator role",
        active: true,
        componentClass: "navigator-on",
        iconClass: "compass",
      },
      {
        name: "navigator",
        description: "navigator role",
        active: false,
        componentClass: "role-off",
        iconClass: "compass",
      },
    ].forEach((testCase) => {
      it(`should work with ${testCase.name} role ${testCase.active ? "on" : "off"}`, fakeAsync(() => {
        // Setup the service to return the expected role
        serviceFake.setRoleState(testCase.name, testCase.active);

        // Set the component input
        component.name = testCase.name;

        // Trigger ngOnInit
        fixture.detectChanges();
        tick();

        // Verify that the component's role attribute is set correctly
        expect(component.role).toBeTruthy();
        expect(component.role?.name).toEqual(testCase.name);
        expect(component.role?.description).toEqual(testCase.description);
        expect(component.role?.active).toEqual(testCase.active);

        // Verify that the component is rendered with the expected CSS class
        const componentElement = fixture.debugElement.query(
          By.css(`[data-testid="role-component"]`),
        );
        expect(componentElement).toBeTruthy();
        expect(componentElement.nativeElement.classList).toContain(
          testCase.componentClass,
        );

        // Verify that the right role icon is rendered
        const iconElement = fixture.debugElement.query(
          By.css(`[data-testid="role-icon"]`),
        );
        expect(iconElement).toBeTruthy();
        // Check for fa-icon element with correct icon
        expect(iconElement.nativeElement.tagName.toLowerCase()).toBe("fa-icon");
        // The icon is rendered as an SVG, check the data-icon attribute
        const svgElement = iconElement.nativeElement.querySelector("svg");
        if (svgElement) {
          const dataIcon = svgElement.getAttribute("data-icon");
          expect(dataIcon).toBe(testCase.iconClass);
        } else {
          // Fallback: check if the icon name is in the class
          expect(iconElement.nativeElement.innerHTML).toContain(
            testCase.iconClass,
          );
        }

        // Verify that the role label is rendered
        const labelElement = fixture.debugElement.query(
          By.css(`[data-testid="role-label"]`),
        );
        expect(labelElement).toBeTruthy();
        expect(labelElement.nativeElement.textContent).toEqual(
          testCase.description,
        );
      }));
    });
  });

  describe("component refresh", () => {
    [
      {
        expectation: "should activate on own role start messages",
        name: "driver",
        description: "driver role",
        activeBefore: false,
        message: "driver:start",
        activeAfter: true,
      },
      {
        expectation: "should deactivate on own role stop messages",
        name: "driver",
        description: "driver role",
        activeBefore: true,
        message: "driver:stop",
        activeAfter: false,
      },
      {
        expectation: "should ignore other role start messages",
        name: "driver",
        description: "driver role",
        activeBefore: false,
        message: "navigator:start",
        activeAfter: false,
      },
      {
        expectation: "should ignore other role stop messages",
        name: "driver",
        description: "driver role",
        activeBefore: true,
        message: "navigator:stop",
        activeAfter: true,
      },
    ].forEach((testCase) => {
      it(`${testCase.expectation}`, fakeAsync(() => {
        // Setup initial state
        serviceFake.setRoleState(testCase.name, testCase.activeBefore);
        component.name = testCase.name;

        // Initialize component
        fixture.detectChanges();
        tick();
        expect(component.role?.active).toEqual(testCase.activeBefore);

        // Update the service state for after the refresh
        serviceFake.setRoleState(testCase.name, testCase.activeAfter);

        // Trigger the refresh method with the message
        component.refresh({
          type: TcrMessageType.ROLE,
          text: testCase.message,
        } as TcrMessage);
        tick();

        // Manual change detection
        changeDetectorRef.detectChanges();

        // Verify that the component's role active attribute was updated
        expect(component.role?.active).toEqual(testCase.activeAfter);
      }));
    });
  });

  describe("component toggleRole", () => {
    [
      {
        expectation: "should activate driver role when off",
        name: "driver",
        active: false,
      },
      {
        expectation: "should deactivate driver role when on",
        name: "driver",
        active: true,
      },
      {
        expectation: "should activate navigator role when off",
        name: "navigator",
        active: false,
      },
      {
        expectation: "should deactivate navigator role when on",
        name: "navigator",
        active: true,
      },
    ].forEach((testCase) => {
      it(`${testCase.expectation}`, fakeAsync(() => {
        // Setup initial state
        serviceFake.setRoleState(testCase.name, testCase.active);
        component.name = testCase.name;
        component.role = {
          name: testCase.name,
          description: "role description",
          active: testCase.active,
        };

        // Spy on the service method
        const activateRoleSpy = spyOn(
          serviceFake,
          "activateRole",
        ).and.callThrough();

        // Initialize view
        fixture.detectChanges();
        tick();

        // Trigger the toggleRole call
        component.toggleRole(component.role!);
        tick();

        // Manual change detection
        changeDetectorRef.detectChanges();

        // Verify that the service was called with the correct parameters
        expect(activateRoleSpy).toHaveBeenCalledWith(
          testCase.name,
          !testCase.active,
        );

        // Verify that the component's role active attribute was toggled
        expect(component.role?.active).toEqual(!testCase.active);
      }));
    });
  });
});
