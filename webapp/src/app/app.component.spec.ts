import {ComponentFixture, TestBed} from '@angular/core/testing';
import {AppComponent} from './app.component';
import {MockComponent, MockDirective} from "ng-mocks";
import {HeaderComponent} from "./components/header/header.component";
import {FooterComponent} from "./components/footer/footer.component";
import {RouterOutlet} from "@angular/router";

describe('AppComponent', () => {
  let app: AppComponent;
  let fixture: ComponentFixture<AppComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppComponent,
        MockComponent(HeaderComponent),
        MockDirective(RouterOutlet),
        MockComponent(FooterComponent),
      ],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(AppComponent);
    app = fixture.componentInstance;
    fixture.detectChanges();
  });

  describe('component instance', () => {
    it('should create the app', () => {
      expect(app).toBeTruthy();
    });

    it(`should have TCR for title`, () => {
      expect(app.title).toEqual('TCR');
    });
  });

  describe('component DOM', () => {
    [
      {selector: 'app-header', description: 'a header element'},
      {selector: 'router-outlet', description: 'a router outlet element'},
      {selector: 'app-footer', description: 'a footer element'},
    ].forEach(({selector, description}) => {
      it(`should have ${description}`, () => {
        expect(fixture.nativeElement.querySelector(selector)).toBeTruthy();
      });
    });
  });
});
