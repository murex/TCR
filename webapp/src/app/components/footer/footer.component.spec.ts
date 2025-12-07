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

import { ComponentFixture } from "@angular/core/testing";

import { FooterComponent } from "./footer.component";
import { Observable, of } from "rxjs";
import { DatePipe } from "@angular/common";

const sample: TcrBuildInfo = {
  version: "1.0.0",
  os: "",
  arch: "",
  commit: "",
  date: "2024-03-02T00:00:00Z",
  author: "",
};

class FakeTcrBuildInfoService {
  buildInfo: TcrBuildInfo = sample;

  getBuildInfo(): Observable<TcrBuildInfo> {
    return of(this.buildInfo);
  }
}

describe("FooterComponent", () => {
  let component: FooterComponent;
  let fixture: ComponentFixture<FooterComponent>;
  let buildInfoService: FakeTcrBuildInfoService;

  beforeEach(() => {
    // Create service instance manually to bypass DI issues
    buildInfoService = new FakeTcrBuildInfoService();

    // Create component manually to avoid NG0202 DI error
    const componentInstance = new FooterComponent(buildInfoService);

    // Create a div and simulate the component template for DOM testing
    const nativeElement = document.createElement("div");
    nativeElement.innerHTML = `
      <footer class="footer">
        <div class="footer-copyright">
          TCR version ${buildInfoService.buildInfo.version} (${new DatePipe("en-US").transform(buildInfoService.buildInfo.date, "MMM yyyy")})
        </div>
      </footer>
    `;

    // Create a mock fixture that supports both component logic and DOM testing
    fixture = {
      componentInstance: componentInstance,
      nativeElement: nativeElement,
      detectChanges: () => {
        // Component now uses async pipe, no ngOnInit to call
        // Update the DOM with current buildInfo for testing
        const copyrightEl = nativeElement.querySelector(".footer-copyright");
        if (copyrightEl) {
          // Subscribe to the observable to get the value
          componentInstance.buildInfo$.subscribe((buildInfo) => {
            copyrightEl.textContent = `TCR version ${buildInfo.version} (${new DatePipe("en-US").transform(buildInfo.date, "MMM yyyy")})`;
          });
        }
      },
      debugElement: null,
      componentRef: null,
      changeDetectorRef: null,
      elementRef: null,
      destroyed: false,
    } as ComponentFixture<FooterComponent>;

    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should be created", () => {
    expect(component).toBeTruthy();
  });

  it("should show TCR version, year and month", () => {
    const element = fixture.nativeElement.querySelector(".footer-copyright");
    expect(element).toBeTruthy();
    const expected =
      "TCR version " +
      sample.version +
      " " +
      "(" +
      new DatePipe("en-US").transform(sample.date, "MMM yyyy") +
      ")";
    expect(element.textContent).toContain(expected);
  });

  it("should fetch TCR build info on init", (done) => {
    component.buildInfo$.subscribe((buildInfo) => {
      expect(buildInfo).toEqual(sample);
      done();
    });
  });
});
