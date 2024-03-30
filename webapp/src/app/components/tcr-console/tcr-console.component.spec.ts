import {ComponentFixture, TestBed} from '@angular/core/testing';
import {
  formatRoleMessage,
  getRoleAction,
  getRoleName, isRoleStartMessage,
  TcrConsoleComponent
} from './tcr-console.component';
import {Observable} from "rxjs";
import {TcrMessage} from "../../interfaces/tcr-message";
import {TcrMessageService} from "../../services/tcr-message.service";
import {TcrRolesComponent} from "../tcr-roles/tcr-roles.component";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {TcrTraceComponent} from "../tcr-trace/tcr-trace.component";

class TcrMessageServiceFake implements Partial<TcrMessageService> {
  constructor(public message$ = new Observable<TcrMessage>()) {
  }
}

class TcrRolesComponentFake implements Partial<TcrRolesComponent> {
}

class TcrTraceComponentFake implements Partial<TcrTraceComponent> {
}

describe('TcrConsoleComponent', () => {
  let component: TcrConsoleComponent;
  let fixture: ComponentFixture<TcrConsoleComponent>;
  let serviceFake: TcrMessageService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrConsoleComponent, HttpClientTestingModule],
      providers: [
        {provide: TcrMessageService, useClass: TcrMessageServiceFake},
        {provide: TcrRolesComponent, useClass: TcrRolesComponentFake},
        {provide: TcrTraceComponent, useClass: TcrTraceComponentFake},
      ]
    }).compileComponents();

    serviceFake = TestBed.inject(TcrMessageService);
    fixture = TestBed.createComponent(TcrConsoleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  describe('component instance', () => {
    it('should be created', () => {
      expect(component).toBeTruthy();
    });
  });

  describe('component DOM', () => {

    it(`should have a title`, () => {
      const element = fixture.nativeElement.querySelector('h1');
      expect(element).toBeTruthy();
      expect(element.textContent).toContain('TCR Console');
    });

    [
      {
        description: 'a roles component',
        selector: 'app-tcr-roles',
      },
      {
        description: 'a terminal',
        selector: 'ng-terminal'
      },
    ].forEach(testCase => {
      it(`should contain ${testCase.description}`, () => {
        const element = fixture.nativeElement.querySelector(testCase.selector);
        expect(element).toBeTruthy();
      });
    });
  });

  describe('getRoleAction function', () => {
    it('should extract the action from the message', () => {
      const message = "driver:start";
      const result = getRoleAction(message);
      expect(result).toEqual("start");
    });

    it('should return an empty string if the message is empty', () => {
      const message = "";
      const result = getRoleAction(message!);
      expect(result).toEqual("");
    });

    it('should return an empty string if the message is undefined', () => {
      const message = undefined;
      const result = getRoleAction(message!);
      expect(result).toEqual("");
    });
  });

  describe('getRoleName function', () => {
    it('should extract the role name from the message', () => {
      const message = "driver:start";
      const result = getRoleName(message);
      expect(result).toEqual("driver");
    });

    it('should return an empty string if the message is empty', () => {
      const message = "";
      const result = getRoleName(message!);
      expect(result).toEqual("");
    });

    it('should return an empty string if the message is undefined', () => {
      const message = undefined;
      const result = getRoleName(message!);
      expect(result).toEqual("");
    });
  });

  describe('formatRoleMessage function', () => {
    [
      {message: "driver:start", expected: "Starting driver role"},
      {message: "driver:end", expected: "Ending driver role"},
      {message: "navigator:start", expected: "Starting navigator role"},
      {message: "navigator:end", expected: "Ending navigator role"},
      {message: "", expected: ""},
      {message: undefined, expected: ""},
    ].forEach(testCase => {
      it(`should format the message \'${testCase.message}\'`, () => {
        const result = formatRoleMessage(testCase.message!);
        expect(result).toEqual(testCase.expected);
      });
    });
  });

  describe('isRoleStartMessage function', () => {
    [
      {message: "driver:start", expected: true},
      {message: "driver:end", expected: false},
      {message: "navigator:start", expected: true},
      {message: "navigator:end", expected: false},
      {message: "other", expected: false},
    ].forEach(testCase => {
      it(`should return ${testCase.expected} for \'${testCase.message}\' messages`, () => {
        const result = isRoleStartMessage(testCase.message);
        expect(result).toEqual(testCase.expected);
      });
    });
  });
});
