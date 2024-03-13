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

  it('returns expected session info when getSessionInfo is called', () => {
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

    service.getSessionInfo().subscribe(sessionInfo => {
      expect(sessionInfo).toEqual(sample);
    });

    const req = httpMock.expectOne(`/api/session-info`);
    expect(req.request.method).toBe('GET');
    expect(req.request.responseType).toEqual('json');
    req.flush(sample);
  });

  it('returns undefined when getSessionInfo receives an error response', () => {
    service.getSessionInfo().subscribe(sessionInfo => {
      expect(sessionInfo).toBeUndefined();
    });

    const req = httpMock.expectOne(`/api/session-info`);
    expect(req.request.method).toBe('GET');
    req.flush(null, {status: 500, statusText: 'Server error'})
  });
});
