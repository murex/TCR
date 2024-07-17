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
import {WebsocketService} from './websocket.service';
import {TcrMessage, TcrMessageType} from "../interfaces/tcr-message";
import {Subject} from "rxjs";
import {TcrMessageService} from "./tcr-message.service";

class FakeWebsocketService {
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
        {provide: WebsocketService, useClass: FakeWebsocketService},
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
