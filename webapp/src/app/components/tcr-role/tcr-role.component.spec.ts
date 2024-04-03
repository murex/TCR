import {ComponentFixture, TestBed} from '@angular/core/testing';

import {TcrRoleComponent} from './tcr-role.component';
import {Observable, of} from "rxjs";
import {TcrMessage, TcrMessageType} from "../../interfaces/tcr-message";
import {TcrRolesService} from "../../services/trc-roles.service";
import {TcrRole} from "../../interfaces/tcr-role";
import {By} from "@angular/platform-browser";

class FakeTcrRolesService {
  message$ = new Observable<TcrMessage>();

  getRole(): Observable<TcrRole> {
    return of({name: "", description: "", active: false});
  }

  activateRole(name: string, state: boolean): Observable<TcrRole> {
    return of({name: name, description: name + " role", active: state});
  }
}

describe('TcrRoleComponent', () => {
  let component: TcrRoleComponent;
  let fixture: ComponentFixture<TcrRoleComponent>;
  let serviceFake: TcrRolesService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrRoleComponent],
      providers: [
        {provide: TcrRolesService, useClass: FakeTcrRolesService},
      ]
    }).compileComponents();
  });

  beforeEach(() => {
    serviceFake = TestBed.inject(TcrRolesService);
    fixture = TestBed.createComponent(TcrRoleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  describe('component instance', () => {
    it('should be created', () => {
      expect(component).toBeTruthy();
    });
  });

  describe('component initialization', () => {
    [
      {
        name: "driver",
        description: "driver role",
        active: true,
        componentClass: 'driver-on',
        iconClass: 'fa-keyboard-o'
      },
      {
        name: "driver",
        description: "driver role",
        active: false,
        componentClass: 'role-off',
        iconClass: 'fa-keyboard-o'
      },
      {
        name: "navigator",
        description: "navigator role",
        active: true,
        componentClass: 'navigator-on',
        iconClass: 'fa-compass'
      },
      {
        name: "navigator",
        description: "navigator role",
        active: false,
        componentClass: 'role-off',
        iconClass: 'fa-compass'
      }
    ].forEach(testCase => {
      it(`should work with ${testCase.name} role ${testCase.active ? "on" : "off"}`, () => {
        const role: TcrRole = {
          name: testCase.name,
          description: testCase.description,
          active: testCase.active,
        };

        // Have the service fake's getRole method return the role
        serviceFake.getRole = () => of(role);
        fixture = TestBed.createComponent(TcrRoleComponent);
        component = fixture.componentInstance;
        component.name = testCase.name;
        fixture.detectChanges();

        // Verify that the component's role attribute is set correctly
        expect(component.role).toEqual(role);

        // Verify that the component is rendered with the expected CSS class
        const componentElement = fixture.debugElement.query(
          By.css(`[data-testid="role-component"]`));
        expect(componentElement).toBeTruthy();
        expect(componentElement.nativeElement.classList).toContain(testCase.componentClass);

        // Verify that the right role icon is rendered
        const iconElement = fixture.debugElement.query(
          By.css(`[data-testid="role-icon"]`));
        expect(iconElement).toBeTruthy();
        expect(iconElement.nativeElement.classList).toContain(testCase.iconClass);

        // Verify that the role label is rendered
        const labelElement = fixture.debugElement.query(
          By.css(`[data-testid="role-label"]`));
        expect(labelElement).toBeTruthy();
        expect(labelElement.nativeElement.textContent).toEqual(role.description);
      });
    });
  });

  describe('component refresh', () => {
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
    ].forEach(testCase => {
      it(`${testCase.expectation}`, () => {
        // Have the service fake's getRole method return the starting role
        const roleBefore: TcrRole = {
          name: testCase.name,
          description: testCase.description,
          active: testCase.activeBefore,
        };
        serviceFake.getRole = () => of(roleBefore);

        fixture = TestBed.createComponent(TcrRoleComponent);
        component = fixture.componentInstance;
        component.name = testCase.name;

        // Verify that the initial role is set correctly
        fixture.detectChanges();
        expect(component.role?.active).toEqual(testCase.activeBefore);

        // Update the service fake to return the expected new role
        const roleAfter: TcrRole = {
          name: testCase.name,
          description: testCase.description,
          active: testCase.activeAfter,
        };
        serviceFake.getRole = () => of(roleAfter);

        // Trigger the refresh method with the message
        component.refresh({
          type: TcrMessageType.ROLE,
          text: testCase.message
        } as TcrMessage);

        // Verify that the component's role active attribute was updated
        fixture.detectChanges();
        expect(component.role?.active).toEqual(testCase.activeAfter);
      });
    });
  });

  describe('component toggleRole', () => {
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
    ].forEach(testCase => {
      it(`${testCase.expectation}`, () => {
        // Set initial role state
        component.role = {
          name: testCase.name,
          description: "role description",
          active: testCase.active,
        };

        // Verify that the initial role is set correctly
        fixture.detectChanges();
        expect(component.role?.active).toEqual(testCase.active);

        // Trigger the toggleRole call
        component.toggleRole(component.role!);

        // Verify that the component's role active attribute was updated
        fixture.detectChanges();
        expect(component.role?.name).toEqual(testCase.name);
        expect(component.role?.active).toEqual(!testCase.active);
      });
    });
  });
});

