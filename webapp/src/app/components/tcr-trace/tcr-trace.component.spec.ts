import {ComponentFixture, TestBed} from '@angular/core/testing';
import {TcrTraceComponent, toCRLF} from './tcr-trace.component';
import {Subject} from "rxjs";

describe('TcrTraceComponent', () => {
  let component: TcrTraceComponent;
  let fixture: ComponentFixture<TcrTraceComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrTraceComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(TcrTraceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  describe('component instance', () => {
    it('should be created', () => {
      expect(component).toBeTruthy();
    });

    it('should have ng-terminal child component', () => {
      expect(component.ngTerminal).toBeTruthy();
    });

    it('should have ng-terminal child.underlying component', () => {
      expect(component.ngTerminal.underlying).toBeTruthy();
    });

    it('should clear the terminal upon reception of clearTrace observable', () => {
      let cleared = false;
      component.ngTerminal.underlying!.reset = () => {
        cleared = true;
      }

      const clearTrace = new Subject<void>();
      component.clearTrace = clearTrace.asObservable();

      component.ngAfterViewInit();
      clearTrace.next();

      expect(cleared).toBeTruthy();
    });

    it('should print text upon reception of text observable', () => {
      let written = "";
      component.ngTerminal.write = (input: string) => {
        written = input;
      }
      const text = new Subject<string>();

      const input = "Hello World";
      component.text = text.asObservable();

      component.ngAfterViewInit();
      text.next(input);

      expect(written).toEqual(input + "\r\n");
    });

  });

  describe('toCRLF function', () => {
    it('should replace all LF with CRLF in the input string', () => {
      const input = "Hello\nWorld\n";
      const result = toCRLF(input);
      expect(result).toEqual("Hello\r\nWorld\r\n\r\n");
    });

    it('should append CRLF to the input string if it does not end with LF', () => {
      const input = "Hello World";
      const result = toCRLF(input);
      expect(result).toEqual("Hello World\r\n");
    });

    it('should return an empty string if the input string is empty', () => {
      const input = "";
      const result = toCRLF(input);
      expect(result).toEqual("");
    });

    it('should return an empty string if the input string is undefined', () => {
      const input = undefined;
      const result = toCRLF(input!);
      expect(result).toEqual("");
    });
  });

  describe('print function', () => {
    it('should send text to the terminal', () => {
      let written = "";
      component.ngTerminal.write = (input: string) => {
        written = input;
      }
      const input = "Hello World";
      component.print(input);
      expect(written).toEqual(input + "\r\n");
    });
  });

  describe('clear function', () => {
    it('should clear the terminal contents', () => {
      let cleared = false;
      component.ngTerminal.underlying!.reset = () => {
        cleared = true;
      }
      component.clear();
      expect(cleared).toBeTruthy();
    });
  });

});
