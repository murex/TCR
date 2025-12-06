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

import { ComponentFixture, TestBed } from "@angular/core/testing";
import { configureComponentTestingModule } from "../../../test-helpers/angular-test-helpers";
import { injectService } from "../../../test-helpers/angular-test-helpers";

import { HomeComponent } from "./home.component";
import { NavigationBehaviorOptions, Router, UrlTree } from "@angular/router";
import { By } from "@angular/platform-browser";
import { FaIconLibrary } from "@fortawesome/angular-fontawesome";
import { registerFontAwesomeIcons } from "../../shared/font-awesome-icons";

class FakeRouter {
  url: string = "";

  navigateByUrl(
    url: string | UrlTree,
    _extras?: NavigationBehaviorOptions,
  ): Promise<boolean> {
    this.url = url.toString();
    return Promise.resolve(true);
  }
}

describe("HomeComponent", () => {
  let component: HomeComponent;
  let fixture: ComponentFixture<HomeComponent>;
  let router: Router;

  beforeEach(async () => {
    await configureComponentTestingModule(
      HomeComponent,
      [],
      [{ provide: Router, useClass: FakeRouter }, FaIconLibrary],
    );

    // Register FontAwesome icons
    const library = injectService(FaIconLibrary);
    registerFontAwesomeIcons(library);
  });

  beforeEach(() => {
    router = injectService(Router);
    fixture = TestBed.createComponent(HomeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  describe("component instance", () => {
    it("should be created", () => {
      expect(component).toBeTruthy();
    });
  });

  describe("component DOM", () => {
    it(`should have a title`, () => {
      const element = fixture.nativeElement.querySelector("h1");
      expect(element).toBeTruthy();
      expect(element.textContent).toContain("TCR - Test && Commit || Revert");
    });

    [
      { buttonId: "console-button", expectedUrl: "/console" },
      { buttonId: "about-button", expectedUrl: "/about" },
      { buttonId: "session-button", expectedUrl: "/session" },
    ].forEach(({ buttonId, expectedUrl }) => {
      it(`should have a clickable link redirecting to the ${expectedUrl} page`, () => {
        const element = fixture.debugElement.query(
          By.css(`[data-testid="${buttonId}"]`),
        );
        expect(element).toBeTruthy();
        element.triggerEventHandler("click", null);
        expect(router.url).toEqual(expectedUrl);
      });
    });

    it("should alert the user on invalid path", async () => {
      spyOn(window, "alert");
      router.navigateByUrl = () => Promise.resolve(false);
      await component.navigateTo("/invalid-path");
      expect(window.alert).toHaveBeenCalled();
    });
  });
});
