import {ComponentFixture, TestBed} from '@angular/core/testing';
import {TcrTraceComponent, toCRLF} from './tcr-trace.component';

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
      expect(component.child).toBeTruthy();
    });

    it('should have ng-terminal child.underlying component', () => {
      expect(component.child.underlying).toBeTruthy();
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
      component.child.write = (input: string) => {
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
      component.child.underlying!.reset = () => {
        cleared = true;
      }
      component.clear();
      expect(cleared).toBeTruthy();
    });
  });

});
