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

import { TestBed } from "@angular/core/testing";
import { HttpClient } from "@angular/common/http";
import {
  HttpClientTestingModule,
  HttpTestingController,
} from "@angular/common/http/testing";
import { Injectable } from "@angular/core";
import { Observable, catchError, of } from "rxjs";
import { createServiceInInjectionContext } from "../../test-helpers/angular-test-helpers";

// Simple inline test service to isolate the issue
@Injectable({
  providedIn: "root",
})
export class TestHttpService {
  constructor(private http: HttpClient) {}

  getData(): Observable<unknown> {
    return this.http
      .get<unknown>("/api/test")
      .pipe(catchError(() => of(undefined)));
  }
}

describe("Service DI Test", () => {
  let service: TestHttpService;
  let httpClient: HttpClient;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [TestHttpService],
    });

    httpClient = TestBed.inject(HttpClient);
    httpMock = TestBed.inject(HttpTestingController);
    service = createServiceInInjectionContext<TestHttpService>(TestHttpService);
  });

  afterEach(() => {
    httpMock?.verify();
    TestBed.resetTestingModule();
  });

  it("should create HttpClient", () => {
    expect(httpClient).toBeTruthy();
  });

  it("should create HttpTestingController", () => {
    expect(httpMock).toBeTruthy();
  });

  it("should create TestHttpService", () => {
    expect(service).toBeTruthy();
  });

  it("should make HTTP request through service", () => {
    const testData = { test: "data" };

    service.getData().subscribe((response) => {
      expect(response).toEqual(testData);
    });

    const req = httpMock.expectOne("/api/test");
    expect(req.request.method).toBe("GET");
    req.flush(testData);
  });
});
