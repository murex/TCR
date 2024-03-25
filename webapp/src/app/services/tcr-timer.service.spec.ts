import {HttpClientTestingModule, HttpTestingController} from '@angular/common/http/testing';
import {TestBed} from '@angular/core/testing';
import {TcrTimerService} from './tcr-timer.service';
import {TcrTimer} from '../interfaces/tcr-timer';
import {WebsocketService} from './websocket.service';
import {TcrMessage, TcrMessageType} from "../interfaces/tcr-message";
import {Subject} from "rxjs";

class WebsocketServiceFake {
  webSocket$: Subject<TcrMessage> = new Subject<TcrMessage>();
}

describe('TcrTimerService', () => {
  let service: TcrTimerService;
  let httpMock: HttpTestingController;
  let wsServiceFake: WebsocketService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [
        TcrTimerService,
        {provide: WebsocketService, useClass: WebsocketServiceFake},
      ]
    });

    service = TestBed.inject(TcrTimerService);
    httpMock = TestBed.inject(HttpTestingController);
    wsServiceFake = TestBed.inject(WebsocketService);
  });

  afterEach(() => {
    httpMock.verify();
  });

  describe('service instance', () => {

    it('should be created', () => {
      expect(service).toBeTruthy();
    });
  });

  describe('getTimer() function', () => {

    it('should return timer info when called', () => {
      const sample: TcrTimer = {
        state: "some-state",
        timeout: "500",
        elapsed: "200",
        remaining: "300",
      };

      let actual: TcrTimer | undefined;
      service.getTimer().subscribe(other => {
        actual = other;
      });

      const req = httpMock.expectOne(`/api/timer`);
      expect(req.request.method).toBe('GET');
      expect(req.request.responseType).toEqual('json');
      req.flush(sample);
      expect(actual).toBe(sample);
    });

    it('should return undefined when receiving an error response', () => {
      let actual: TcrTimer | undefined;
      service.getTimer().subscribe(other => {
        actual = other;
      });

      const req = httpMock.expectOne(`/api/timer`);
      expect(req.request.method).toBe('GET');
      req.flush({message: 'Some network error'}, {status: 500, statusText: 'Server Error'});
      expect(actual).toBeUndefined();
    });
  });

  describe('websocket message handler', () => {

    it('should forward timer messages', (done) => {
      const sampleMessage: TcrMessage = {
        type: TcrMessageType.TIMER,
        emphasis: false,
        severity: "",
        text: "",
        timestamp: "",
      };
      let actual: TcrMessage | undefined;
      service.message$.subscribe((msg) => {
        actual = msg;
        done();
      });
      wsServiceFake.webSocket$.next(sampleMessage);
      expect(actual).toEqual(sampleMessage);
    });

    it('should drop non-timer messages', (done) => {
      const sampleMessage: TcrMessage = {
        type: TcrMessageType.INFO,
        emphasis: false,
        severity: "",
        text: "",
        timestamp: "",
      };
      let actual: TcrMessage | undefined;
      service.message$.subscribe((msg) => {
        actual = msg;
        done();
      });
      wsServiceFake.webSocket$.next(sampleMessage);
      // Wait for the message to be processed by the service before checking the result
      setTimeout(() => done(), 10);
      expect(actual).toBeUndefined();
    });
  });
});
