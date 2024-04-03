import {ComponentFixture, TestBed} from '@angular/core/testing';
import {HeaderComponent} from './header.component';
import {RouterModule} from "@angular/router";
import {TcrTimerComponent} from "../tcr-timer/tcr-timer.component";
import {MockComponent} from "ng-mocks";

describe('HeaderComponent', () => {
  let component: HeaderComponent;
  let fixture: ComponentFixture<HeaderComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        HeaderComponent,
        MockComponent(TcrTimerComponent),
        RouterModule.forRoot([]),
      ],
      providers: [],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(HeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  describe('component instance', () => {
    it('should be created', () => {
      expect(component).toBeTruthy();
    });
  });

  describe('component DOM', () => {
    [
      {
        description: 'a navigation bar',
        selector: '.navbar'
      },
      {
        description: 'a brand element called TCR',
        selector: '.navbar-brand',
        text: 'TCR'
      },
      {
        description: 'a link to the home page',
        selector: 'a[href="/"]',
        text: 'Home',
      },
      {
        description: 'a link to the session page',
        selector: 'a[href="/session"]',
        text: 'Session',
      },
      {
        description: 'a link to the console page',
        selector: 'a[href="/console"]',
        text: 'Console',
      },
      {
        description: 'a link to the about page',
        selector: 'a[href="/about"]',
        text: 'About',
      },
      {
        description: 'a timer component',
        selector: 'app-tcr-timer',
      },
    ].forEach(({selector, description, text}) => {
      it(`should contain ${description}`, () => {
        const element = fixture.nativeElement.querySelector(selector);
        expect(element).toBeTruthy();
        if (text) {
          expect(element.textContent).toContain(text);
        }
      });
    });
  });
});
