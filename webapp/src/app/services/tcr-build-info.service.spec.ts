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

import {TestBed} from '@angular/core/testing';
import {
  HttpTestingController,
  provideHttpClientTesting
} from '@angular/common/http/testing';
import {TcrBuildInfoService} from './tcr-build-info.service';
import {TcrBuildInfo} from '../interfaces/tcr-build-info';
import {provideHttpClient, withInterceptorsFromDi} from '@angular/common/http';

describe('TcrBuildInfoService', () => {
  let service: TcrBuildInfoService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [],
      providers: [
        TcrBuildInfoService,
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    });

    service = TestBed.inject(TcrBuildInfoService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  describe('service instance', () => {

    it('should be created', () => {
      expect(service).toBeTruthy();
    });

  });

  describe('getBuildInfo() function', () => {

    it('should return build info when called', () => {
      const sample: TcrBuildInfo = {
        version: "1.0.0",
        os: "some-os",
        arch: "some-arch",
        commit: "abc123",
        date: "2024-01-01T00:00:00Z",
        author: "some-author",
      };

      let actual: TcrBuildInfo | undefined;
      service.getBuildInfo().subscribe(other => {
        actual = other;
      });

      const req = httpMock.expectOne(`/api/build-info`);
      expect(req.request.method).toBe('GET');
      expect(req.request.responseType).toEqual('json');
      req.flush(sample);
      expect(actual).toEqual(sample);
    });

    it('should return undefined when receiving an error response', () => {
      let actual: TcrBuildInfo | undefined;
      service.getBuildInfo().subscribe(other => {
        actual = other;
      });

      const req = httpMock.expectOne(`/api/build-info`);
      expect(req.request.method).toBe('GET');
      req.flush({message: 'Some network error'}, {
        status: 500,
        statusText: 'Server Error'
      });
      expect(actual).toBeUndefined();
    });

  });

});
