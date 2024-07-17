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
import {Subject} from 'rxjs';
import {WebsocketService} from './websocket.service';
import {TcrMessage, TcrMessageType} from "../interfaces/tcr-message";

// Mocking the websocket
let fakeSocket: Subject<TcrMessage>;
const fakeSocketCtor = jasmine
  .createSpy('WEBSOCKET_CTOR')
  .and.callFake(() => fakeSocket);

describe('WebsocketService', () => {
  const sampleMessage: TcrMessage = {
    emphasis: false,
    type: TcrMessageType.INFO,
    severity: "0",
    text: "some info message",
    timestamp: "2024-01-01T00:00:00Z",
  };
  let service: WebsocketService;

  beforeEach(() => {
    TestBed.runInInjectionContext(() => {
      // Make a new socket so we don't get lingering values leaking across tests
      fakeSocket = new Subject<TcrMessage>();
      // Spy on it, so we don't have to subscribe to verify it was called
      spyOn(fakeSocket, 'next').and.callThrough();
      // Reset your spies
      fakeSocketCtor.calls.reset();
      // Make the service using the fake ctor
      service = new WebsocketService(fakeSocketCtor);
    });
  });

  describe('service instance', () => {
    it('should be created', () => {
      expect(service).toBeTruthy();
    });

    it('should attempt a websocket connection on create', () => {
      const expectedUrl = 'ws://' + window.location.host + '/ws';
      expect(fakeSocketCtor).toHaveBeenCalledOnceWith(expectedUrl);
    });

    it('should be able to forward received TCR messages', (done) => {
      let actual: TcrMessage | undefined;
      service.webSocket$.subscribe((msg) => {
        actual = msg;
        done();
      });
      fakeSocket.next(sampleMessage);
      expect(actual).toBe(sampleMessage);
    });

    it('should handle websocket errors', () => {
      const sampleError = new Error('WebSocket error');
      let actual: Error | undefined;
      service.webSocket$.asObservable().subscribe({
        error: (err) => actual = err,
      });
      fakeSocket.error(sampleError);
      expect(actual).toEqual(sampleError);
    });

    it('should close the websocket on destroy', () => {
      spyOn(fakeSocket, 'complete');
      service.ngOnDestroy();
      expect(fakeSocket.complete).toHaveBeenCalled();
    });
  });
});
