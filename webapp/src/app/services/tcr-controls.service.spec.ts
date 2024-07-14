import {TestBed} from '@angular/core/testing';

import {TcrControlsService} from './tcr-controls.service';
import {
  HttpClientTestingModule,
  HttpTestingController
} from "@angular/common/http/testing";

describe('TcrControlsService', () => {
  let service: TcrControlsService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [
        TcrControlsService,
      ]
    });
    service = TestBed.inject(TcrControlsService);
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

  describe('abortCommand() function', () => {
    it(`should send an HTTP POST abort-command request`, () => {
      service.abortCommand().subscribe();

      const req = httpMock.expectOne(`/api/controls/abort-command`);
      expect(req.request.method).toBe('POST');
      req.flush({}, {
        status: 200,
        statusText: ''
      });
      expect(req.request.responseType).toEqual('json');
    });

    it('should return undefined when receiving an error response', () => {
      let actual: unknown | undefined;
      service.abortCommand().subscribe(other => {
        actual = other;
      });

      const req = httpMock.expectOne(`/api/controls/abort-command`);
      expect(req.request.method).toBe('POST');
      req.flush({message: 'Bad Request'}, {
        status: 400,
        statusText: 'Bad Request'
      });
      expect(actual).toBeUndefined();
    });
  });
});
