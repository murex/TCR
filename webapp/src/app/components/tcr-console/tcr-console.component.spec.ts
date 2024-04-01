import {ComponentFixture, TestBed} from '@angular/core/testing';
import {
  formatRoleMessage,
  getRoleAction,
  getRoleName,
  isRoleStartMessage,
  TcrConsoleComponent
} from './tcr-console.component';
import {Observable} from "rxjs";
import {TcrMessage} from "../../interfaces/tcr-message";
import {TcrMessageService} from "../../services/tcr-message.service";
import {TcrRolesComponent} from "../tcr-roles/tcr-roles.component";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {TcrTraceComponent} from "../tcr-trace/tcr-trace.component";
import {
  bgDarkGray,
  cyan,
  green,
  lightCyan,
  lightYellow,
  red,
  yellow
} from "ansicolor";

class TcrMessageServiceFake implements Partial<TcrMessageService> {
  message$ = new Observable<TcrMessage>()
}

class TcrRolesComponentFake {
}

class TcrTraceComponentFake {
}

describe('TcrConsoleComponent', () => {
  let component: TcrConsoleComponent;
  let fixture: ComponentFixture<TcrConsoleComponent>;
  // let serviceFake: TcrMessageService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrConsoleComponent, HttpClientTestingModule],
      providers: [
        {provide: TcrMessageService, useClass: TcrMessageServiceFake},
        {provide: TcrRolesComponent, useClass: TcrRolesComponentFake},
        {provide: TcrTraceComponent, useClass: TcrTraceComponentFake},
      ]
    }).compileComponents();

    TestBed.inject(TcrMessageService);
    TestBed.inject(TcrRolesComponent);
    TestBed.inject(TcrTraceComponent);
    fixture = TestBed.createComponent(TcrConsoleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  afterEach(() => {
    fixture.destroy();
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
        description: 'a trace component',
        selector: 'app-tcr-trace',
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
      it(`should format the message '${testCase.message}'`, () => {
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
      it(`should return ${testCase.expected} for '${testCase.message}' messages`, () => {
        const result = isRoleStartMessage(testCase.message);
        expect(result).toEqual(testCase.expected);
      });
    });
  });

  describe('printSimple function', () => {
    it('should print a simple text', () => {
      const text = "some simple text";

      let actual: string | undefined;
      component.text.subscribe((value) => {
        actual = value;
      });

      component.printSimple(text);
      expect(actual).toEqual(text);
    });
  });

  describe('printInfo function', () => {
    it('should format and print an info text', () => {
      const text = "some info text";

      let actual: string | undefined;
      component.text.subscribe((value) => {
        actual = value;
      });

      component.printInfo(text);
      expect(actual).toEqual(cyan(text));
    });
  });

  describe('printTitle function', () => {
    it('should format and print a title text', () => {
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

  describe('printRole function', () => {
    it('should format and print a role text', () => {
      const text = "driver:start";

      let actual: string | undefined;
      component.text.subscribe((value) => {
        actual = value;
      });

      component.printRole(text);
      const lineSep = yellow("─".repeat(80));
      expect(actual).toEqual(lineSep
        + "\n" + lightYellow("Starting driver role")
        + "\n" + lineSep);
    });
  });

  describe('printSuccess function', () => {
    it('should format and print a success text', () => {
      const text = "some success text";

      let actual: string | undefined;
      component.text.subscribe((value) => {
        actual = value;
      });

      component.printSuccess(text);
      expect(actual).toEqual("🟢- " + green(text));
    });
  });

  describe('printWarning function', () => {
    it('should format and print a warning text', () => {
      const text = "some warning text";

      let actual: string | undefined;
      component.text.subscribe((value) => {
        actual = value;
      });

      component.printWarning(text);
      expect(actual).toEqual("🔶- " + yellow(text));
    });
  });

  describe('printError function', () => {
    it('should format and print an error text', () => {
      const text = "some error text";

      let actual: string | undefined;
      component.text.subscribe((value) => {
        actual = value;
      });

      component.printError(text);
      expect(actual).toEqual("🟥- " + red(text));
    });
  });

  describe('printUnhandled function', () => {
    it('should format and print an unhandled message', () => {
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

  xdescribe('clear function', () => {
    it('should update clearTrace subject (disabled due to side effect on xterm.js)', () => {

      let actual: boolean = false;
      component.clearTrace.subscribe(() => {
        actual = true;
      });

      component.clear();
      expect(actual).toBeTruthy();
    });
  });
});

