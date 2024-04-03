import {ComponentFixture, TestBed} from '@angular/core/testing';

import {HomeComponent} from './home.component';
import {NavigationBehaviorOptions, Router, UrlTree} from "@angular/router";
import {By} from "@angular/platform-browser";

class FakeRouter {
  url: string = '';

  navigateByUrl(url: string | UrlTree, _extras?: NavigationBehaviorOptions): Promise<boolean> {
    this.url = url.toString();
    return Promise.resolve(true);
  }
}

describe('HomeComponent', () => {
  let component: HomeComponent;
  let fixture: ComponentFixture<HomeComponent>;
  let router: Router;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HomeComponent],
      providers: [
        {provide: Router, useClass: FakeRouter}
      ]
    }).compileComponents();
  });

  beforeEach(() => {
    router = TestBed.inject(Router);
    fixture = TestBed.createComponent(HomeComponent);
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
      expect(element.textContent).toContain('TCR - Test && Commit || Revert');
    });

    [
      {buttonId: 'console-button', expectedUrl: '/console'},
      {buttonId: 'about-button', expectedUrl: '/about'},
      {buttonId: 'session-button', expectedUrl: '/session'},
    ].forEach(({buttonId, expectedUrl}) => {
      it(`should have a clickable link redirecting to the ${expectedUrl} page`, () => {
        const element = fixture.debugElement.query(
          By.css(`[data-testid="${buttonId}"]`));
        expect(element).toBeTruthy();
        element.triggerEventHandler('click', null);
        expect(router.url).toEqual(expectedUrl);
      });
    });

    it('should alert the user on invalid path', async () => {
      spyOn(window, 'alert');
      router.navigateByUrl = () => Promise.resolve(false);
      await component.navigateTo('/invalid-path');
      expect(window.alert).toHaveBeenCalled();
    });

  });
});
