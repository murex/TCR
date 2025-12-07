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
import { TcrBuildInfoService } from "../../services/tcr-build-info.service";
import { TcrAboutComponent } from "./tcr-about.component";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { TcrBuildInfo } from "../../interfaces/tcr-build-info";
import {
  provideHttpClient,
  withInterceptorsFromDi,
} from "@angular/common/http";
import { FaIconLibrary } from "@fortawesome/angular-fontawesome";
import {
  FONT_AWESOME_TEST_PROVIDERS,
  setupFontAwesomeTestingModule,
} from "../../../testing/font-awesome-test-setup";
import { Injectable } from "@angular/core";

const sample: TcrBuildInfo = {
  version: "1.0.0",
  os: "some-os",
  arch: "some-arch",
  commit: "abc123",
  date: "2024-01-01T00:00:00Z",
  author: "some-author",
};

@Injectable({
  providedIn: "root",
})
class FakeTcrBuildInfoService {
  buildInfo: TcrBuildInfo = sample;

  getBuildInfo(): Observable<TcrBuildInfo> {
    return of(this.buildInfo);
  }
}

describe("TcrAboutComponent", () => {
  let component: TcrAboutComponent;
  let fixture: ComponentFixture<TcrAboutComponent>;

  beforeEach(async () => {
    await configureComponentTestingModule(
      TcrAboutComponent,
      [],
      [
        { provide: TcrBuildInfoService, useClass: FakeTcrBuildInfoService },
        ...FONT_AWESOME_TEST_PROVIDERS,
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    );

    // Register FontAwesome icons
    const library = injectService(FaIconLibrary);
    setupFontAwesomeTestingModule(library);
  });

  beforeEach(() => {
    // Use the multi-strategy approach to handle DI issues gracefully
    const serviceFake = TestBed.inject(
      TcrBuildInfoService,
    ) as unknown as FakeTcrBuildInfoService;

    const dependencies = {
      buildInfoService: serviceFake,
    };

    fixture = createComponentWithStrategies(TcrAboutComponent, dependencies);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should be created", () => {
    expect(component).toBeTruthy();
  });

  it('should have title "About TCR"', () => {
    expect(component.title).toEqual("About TCR");
  });

  it("should fetch TCR build info on init", () => {
    expect(component.buildInfo).toEqual(sample);
  });
});
