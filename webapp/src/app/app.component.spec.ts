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
import { AppComponent } from "./app.component";
import { Component } from "@angular/core";
import { RouterModule } from "@angular/router";

// Mock components for testing
@Component({
  selector: "app-header",
  template: '<div class="mock-header"></div>',
  standalone: true,
})
class MockHeaderComponent {}

@Component({
  selector: "app-footer",
  template: '<div class="mock-footer"></div>',
  standalone: true,
})
class MockFooterComponent {}

describe("AppComponent", () => {
  let app: AppComponent;
  let fixture: ComponentFixture<AppComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppComponent,
        MockHeaderComponent,
        MockFooterComponent,
        RouterModule.forRoot([]),
      ],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(AppComponent);
    app = fixture.componentInstance;
    fixture.detectChanges();
  });

  describe("component instance", () => {
    it("should create the app", () => {
      expect(app).toBeTruthy();
    });

    it(`should have TCR for title`, () => {
      expect(app.title).toEqual("TCR");
    });
  });

  describe("component DOM", () => {
    [
      { selector: "app-header", description: "a header element" },
      { selector: "router-outlet", description: "a router outlet element" },
      { selector: "app-footer", description: "a footer element" },
    ].forEach(({ selector, description }) => {
      it(`should have ${description}`, () => {
        expect(fixture.nativeElement.querySelector(selector)).toBeTruthy();
      });
    });
  });
});
