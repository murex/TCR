import {ComponentFixture, TestBed} from '@angular/core/testing';
import {AppComponent} from './app.component';
import {RouterModule} from "@angular/router";
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('AppComponent', () => {
  let app: AppComponent;
  let fixture: ComponentFixture<AppComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppComponent, HttpClientTestingModule, RouterModule.forRoot([])],
    }).compileComponents();

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
    const testCases = [
      {selector: 'app-header', description: 'a header element'},
      {selector: 'router-outlet', description: 'a router outlet element'},
      {selector: 'app-footer', description: 'a footer element'},
    ];

    testCases.forEach(({selector, description}) => {
      it(`should have ${description}`, () => {
        expect(fixture.nativeElement.querySelector(selector)).toBeTruthy();
      });
    });
  });
});
