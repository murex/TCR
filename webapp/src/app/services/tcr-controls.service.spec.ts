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

import {
  injectService,
  configureServiceTestingModule,
  cleanupAngularTest,
} from "../../test-helpers/angular-test-helpers";
import { TcrControlsService } from "./tcr-controls.service";
import { HttpTestingController } from "@angular/common/http/testing";

describe("TcrControlsService", () => {
  let service: TcrControlsService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    configureServiceTestingModule(TcrControlsService);
    service = injectService(TcrControlsService);
    httpMock = injectService(HttpTestingController);
  });

  afterEach(() => {
    cleanupAngularTest(httpMock);
  });

  describe("service instance", () => {
    it("should be created", () => {
      expect(service).toBeTruthy();
    });
  });

  describe("abortCommand() function", () => {
    it(`should send an HTTP POST abort-command request`, () => {
      service.abortCommand().subscribe();

      const req = httpMock.expectOne(`/api/controls/abort-command`);
      expect(req.request.method).toBe("POST");
      req.flush(
        {},
        {
          status: 200,
          statusText: "",
        },
      );
      expect(req.request.responseType).toEqual("json");
    });

    it("should return undefined when receiving an error response", () => {
      let actual: unknown | undefined;
      service.abortCommand().subscribe((other) => {
        actual = other;
      });

      const req = httpMock.expectOne(`/api/controls/abort-command`);
      expect(req.request.method).toBe("POST");
      req.flush(
        { message: "Bad Request" },
        {
          status: 400,
          statusText: "Bad Request",
        },
      );
      expect(actual).toBeUndefined();
    });
  });
});
