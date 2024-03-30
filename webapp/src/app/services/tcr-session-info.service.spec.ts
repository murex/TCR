import {TestBed} from '@angular/core/testing';
import {HttpClientTestingModule, HttpTestingController} from '@angular/common/http/testing';
import {TcrSessionInfoService} from './tcr-session-info.service';
import {TcrSessionInfo} from '../interfaces/tcr-session-info';

describe('TcrSessionInfoService', () => {
  let service: TcrSessionInfoService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [TcrSessionInfoService]
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
      req.flush({message: 'Some network error'}, {status: 500, statusText: 'Server error'});
      expect(actual).toBeUndefined();
    });
  });
});
