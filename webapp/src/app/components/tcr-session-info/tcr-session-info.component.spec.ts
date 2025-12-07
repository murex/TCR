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

import { Observable, of } from "rxjs";
import {
  configureComponentTestingModule,
  injectService,
  createComponentWithStrategies,
} from "../../../test-helpers/angular-test-helpers";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { TcrSessionInfo } from "../../interfaces/tcr-session-info";
import { TcrSessionInfoService } from "../../services/tcr-session-info.service";
import { TcrSessionInfoComponent } from "./tcr-session-info.component";
import { FaIconLibrary } from "@fortawesome/angular-fontawesome";
import {
  FONT_AWESOME_TEST_PROVIDERS,
  setupFontAwesomeTestingModule,
} from "../../../testing/font-awesome-test-setup";
import { Injectable } from "@angular/core";

const sample: TcrSessionInfo = {
  baseDir: "/my/base/dir",
  variant: "relaxed",
  gitAutoPush: false,
  language: "java",
  messageSuffix: "my-suffix",
  toolchain: "gradle",
  vcsName: "git",
  vcsSession: "my VCS session",
  workDir: "/my/work/dir",
};

@Injectable({
  providedIn: "root",
})
class FakeTcrSessionInfoService {
  sessionInfo: TcrSessionInfo = sample;

  getSessionInfo(): Observable<TcrSessionInfo> {
    return of(this.sessionInfo);
  }
}

describe("TcrSessionInfoComponent", () => {
  let component: TcrSessionInfoComponent;
  let fixture: ComponentFixture<TcrSessionInfoComponent>;

  beforeEach(async () => {
    await configureComponentTestingModule(
      TcrSessionInfoComponent,
      [],
      [
        { provide: TcrSessionInfoService, useClass: FakeTcrSessionInfoService },
        ...FONT_AWESOME_TEST_PROVIDERS,
      ],
    );

    // Register FontAwesome icons
    const library = injectService(FaIconLibrary);
    setupFontAwesomeTestingModule(library);
  });

  beforeEach(() => {
    // Use the multi-strategy approach to handle DI issues gracefully
    const serviceFake = TestBed.inject(
      TcrSessionInfoService,
    ) as unknown as FakeTcrSessionInfoService;

    const dependencies = {
      sessionInfoService: serviceFake,
    };

    fixture = createComponentWithStrategies(
      TcrSessionInfoComponent,
      dependencies,
    );
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  describe("component instance", () => {
    it("should be created", () => {
      expect(component).toBeTruthy();
    });

    it('should have title "TCR Session Information"', () => {
      expect(component.title).toEqual("TCR Session Information");
    });
  });

  describe("component initialization", () => {
    it("should fetch TCR session info on init", async () => {
      const sessionInfo = await new Promise<TcrSessionInfo>((resolve) => {
        component.sessionInfo$.subscribe((sessionInfo) => {
          resolve(sessionInfo);
        });
      });
      expect(sessionInfo).toEqual(sample);
    });
  });
});
