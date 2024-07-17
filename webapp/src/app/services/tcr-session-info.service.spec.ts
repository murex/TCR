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
import {TcrSessionInfoService} from './tcr-session-info.service';
import {TcrSessionInfo} from '../interfaces/tcr-session-info';
import {provideHttpClient, withInterceptorsFromDi} from '@angular/common/http';

describe('TcrSessionInfoService', () => {
  let service: TcrSessionInfoService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [],
      providers: [TcrSessionInfoService, provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
    });

    service = TestBed.inject(TcrSessionInfoService);
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

  describe('getSessionInfo() function', () => {
    it('should return session info when called', () => {
      const sample: TcrSessionInfo = {
        baseDir: "/my/base/dir",
        commitOnFail: false,
        gitAutoPush: false,
        language: "java",
        messageSuffix: "my-suffix",
        toolchain: "gradle",
        vcsName: "git",
        vcsSession: "my VCS session",
        workDir: "/my/work/dir"
      };

      let actual: TcrSessionInfo | undefined;
      service.getSessionInfo().subscribe(other => {
        actual = other;
      });

      const req = httpMock.expectOne(`/api/session-info`);
      expect(req.request.method).toBe('GET');
      expect(req.request.responseType).toEqual('json');
      req.flush(sample);
      expect(actual).toEqual(sample);
    });

    it('should return undefined when receiving an error response', () => {
      let actual: TcrSessionInfo | undefined;
      service.getSessionInfo().subscribe(other => {
        actual = other;
      });

      const req = httpMock.expectOne(`/api/session-info`);
      expect(req.request.method).toBe('GET');
      req.flush({message: 'Some network error'}, {
        status: 500,
        statusText: 'Server error'
      });
      expect(actual).toBeUndefined();
    });
  });
});
