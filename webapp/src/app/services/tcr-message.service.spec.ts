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

import {
  injectService,
  configureServiceTestingModule,
  createServiceInInjectionContext,
} from "../../test-helpers/angular-test-helpers";
import { WebsocketService } from "./websocket.service";
import { TcrMessage, TcrMessageType } from "../interfaces/tcr-message";
import { Subject } from "rxjs";
import { TcrMessageService } from "./tcr-message.service";
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

describe("TcrMessageService", () => {
  let service: TcrMessageService;
  let wsServiceFake: WebsocketService;

  beforeEach(() => {
    configureServiceTestingModule(TcrMessageService, [
      { provide: WebsocketService, useClass: FakeWebsocketService },
      { provide: DestroyRef, useClass: MockDestroyRef },
    ]);

    // Create service using injection context helper to handle takeUntilDestroyed
    service = createServiceInInjectionContext<TcrMessageService>(
      TcrMessageService,
      [{ provide: WebsocketService, useClass: FakeWebsocketService }],
    );

    wsServiceFake = injectService(WebsocketService);
  });

  describe("service instance", () => {
    it("should be created", () => {
      expect(service).toBeTruthy();
    });
  });

  describe("websocket message handler", () => {
    Object.values(TcrMessageType).forEach((type) => {
      it(`should forward ${type} messages`, async () => {
        const sampleMessage: TcrMessage = {
          type: type,
          emphasis: false,
          severity: "",
          text: "",
          timestamp: "",
        };

        const messagePromise = new Promise<TcrMessage>((resolve) => {
          service.message$.subscribe((msg) => {
            resolve(msg);
          });
        });

        wsServiceFake.webSocket$.next(sampleMessage);

        const receivedMessage = await messagePromise;
        expect(receivedMessage).toEqual(sampleMessage);
      });
    });
  });
});
