import {TestBed} from '@angular/core/testing';
import {HttpClientTestingModule, HttpTestingController} from '@angular/common/http/testing';
import {TcrBuildInfoService} from './tcr-build-info.service';
import {TcrBuildInfo} from '../interfaces/tcr-build-info';

describe('TcrBuildInfoService', () => {
  let service: TcrBuildInfoService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [TcrBuildInfoService]
    });

    service = TestBed.inject(TcrBuildInfoService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should return expected build info when getBuildInfo is called', () => {
    const sample: TcrBuildInfo = {
      version: "1.0.0",
      os: "some-os",
      arch: "some-arch",
      commit: "abc123",
      date: "2022-01-01T00:00:00Z",
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

  it('should return undefined when getBuildInfo receives an error response', () => {
    let actual: TcrBuildInfo | undefined;
    service.getBuildInfo().subscribe(other => {
      actual = other;
    });

    const req = httpMock.expectOne(`/api/build-info`);
    expect(req.request.method).toBe('GET');
    req.flush({message: 'Some network error'}, {status: 500, statusText: 'Server Error'});
    expect(actual).toBeUndefined();
  });
});
