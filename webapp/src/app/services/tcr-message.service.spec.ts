import {TestBed} from '@angular/core/testing';
import {WebsocketService} from './websocket.service';
import {TcrMessage, TcrMessageType} from "../interfaces/tcr-message";
import {Subject} from "rxjs";
import {TcrMessageService} from "./tcr-message.service";

class WebsocketServiceFake {
  webSocket$: Subject<TcrMessage> = new Subject<TcrMessage>();
}

describe('TcrMessageService', () => {
  let service: TcrMessageService;
  let wsServiceFake: WebsocketService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [],
      providers: [
        TcrMessageService,
        {provide: WebsocketService, useClass: WebsocketServiceFake},
      ]
    });

    service = TestBed.inject(TcrMessageService);
    wsServiceFake = TestBed.inject(WebsocketService);
  });

  describe('service instance', () => {
    it('should be created', () => {
      expect(service).toBeTruthy();
    });
  });

  describe('websocket message handler', () => {
    Object.values(TcrMessageType).forEach(type => {
      it(`should forward ${type} messages`, (done) => {
        const sampleMessage: TcrMessage = {
          type: type,
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
    });
  });
});
