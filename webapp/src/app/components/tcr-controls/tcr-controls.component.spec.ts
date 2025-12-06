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

import { TcrControlsComponent } from "./tcr-controls.component";
import { TcrControlsService } from "../../services/tcr-controls.service";
import { Observable, of } from "rxjs";

class FakeTcrControlsService {
  abortCommand(): Observable<unknown> {
    return of({});
  }
}

describe("TcrControlsComponent", () => {
  let component: TcrControlsComponent;
  let fixture: ComponentFixture<TcrControlsComponent>;
  let serviceFake: TcrControlsService;

  beforeEach(async () => {
    await configureComponentTestingModule(
      TcrControlsComponent,
      [],
      [{ provide: TcrControlsService, useClass: FakeTcrControlsService }],
    );
  });

  beforeEach(() => {
    serviceFake = injectService(TcrControlsService);
    fixture = TestBed.createComponent(TcrControlsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  describe("component instance", () => {
    it("should be created", () => {
      expect(component).toBeTruthy();
    });
  });

  describe("abort command button", () => {
    it("should trigger abort command from controls service", () => {
      const abortCommandFunction = spyOn(
        serviceFake,
        "abortCommand",
      ).and.callThrough();
      // Trigger the abortCommand call
      component.abortCommand();
      // Verify that the service received the request
      expect(abortCommandFunction).toHaveBeenCalledTimes(1);
    });
  });
});
