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
      providers: [TcrBuildInfoService, provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
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
