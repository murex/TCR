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

import { TestBed } from "@angular/core/testing";
import { HttpTestingController } from "@angular/common/http/testing";
import {
  configureServiceTestingModule,
  cleanupAngularTest,
  createServiceInInjectionContext,
} from "../../test-helpers/angular-test-helpers";
import { TcrTimerService } from "./tcr-timer.service";
import { TcrTimer } from "../interfaces/tcr-timer";
import { WebsocketService } from "./websocket.service";
import { Subject } from "rxjs";
import { TcrMessage, TcrMessageType } from "../interfaces/tcr-message";
import { DestroyRef } from "@angular/core";

class FakeWebsocketService {
  webSocket$: Subject<TcrMessage> = new Subject<TcrMessage>();
}

// Mock DestroyRef for takeUntilDestroyed
class MockDestroyRef {
  onDestroy(_fn: () => void) {
    // Mock implementation - does nothing for tests
  }
}

describe("TcrTimerService", () => {
  let service: TcrTimerService;
  let httpMock: HttpTestingController;
  let wsServiceFake: WebsocketService;

  beforeEach(() => {
    configureServiceTestingModule(TcrTimerService, [
      { provide: WebsocketService, useClass: FakeWebsocketService },
      { provide: DestroyRef, useClass: MockDestroyRef },
    ]);

    // Create service using injection context helper to handle takeUntilDestroyed
    service = createServiceInInjectionContext<TcrTimerService>(
      TcrTimerService,
      [{ provide: WebsocketService, useClass: FakeWebsocketService }],
    );

    httpMock = TestBed.inject(HttpTestingController);
    wsServiceFake = TestBed.inject(WebsocketService);
  });

  afterEach(() => {
    cleanupAngularTest(httpMock);
  });

  describe("service instance", () => {
    it("should be created", () => {
      expect(service).toBeTruthy();
    });
  });

  describe("getTimer() function", () => {
    it("should return timer info when called", () => {
      const sample: TcrTimer = {
        state: "some-state",
        timeout: "500",
        elapsed: "200",
        remaining: "300",
      };

      let actual: TcrTimer | undefined;
      service.getTimer().subscribe((other) => {
        actual = other;
      });

      const req = httpMock.expectOne(`/api/timer`);
      expect(req.request.method).toBe("GET");
      expect(req.request.responseType).toEqual("json");
      req.flush(sample);
      expect(actual).toBe(sample);
    });

    it("should return undefined when receiving an error response", () => {
      let actual: TcrTimer | undefined;
      service.getTimer().subscribe((other) => {
        actual = other;
      });

      const req = httpMock.expectOne(`/api/timer`);
      expect(req.request.method).toBe("GET");
      req.flush(
        { message: "Some network error" },
        {
          status: 500,
          statusText: "Server Error",
        },
      );
      expect(actual).toBeUndefined();
    });
  });

  describe("websocket message handler", () => {
    it("should forward timer messages", async () => {
      const sampleMessage = { type: TcrMessageType.TIMER } as TcrMessage;

      const messagePromise = new Promise<TcrMessage>((resolve) => {
        service.message$.subscribe((msg) => {
          resolve(msg);
        });
      });

      wsServiceFake.webSocket$.next(sampleMessage);

      const receivedMessage = await messagePromise;
      expect(receivedMessage).toEqual(sampleMessage);
    });

    it("should drop non-timer messages", async () => {
      const timerMessage = { type: TcrMessageType.TIMER } as TcrMessage;
      const infoMessage = { type: TcrMessageType.INFO } as TcrMessage;

      const messagePromise = new Promise<TcrMessage>((resolve) => {
        service.message$.subscribe((msg) => {
          // Should only receive timer messages
          resolve(msg);
        });
      });

      // Send non-timer message first (should be filtered out)
      wsServiceFake.webSocket$.next(infoMessage);
      // Send timer message (should be received)
      wsServiceFake.webSocket$.next(timerMessage);

      const receivedMessage = await messagePromise;
      expect(receivedMessage.type).toEqual(TcrMessageType.TIMER);
    });
  });
});
