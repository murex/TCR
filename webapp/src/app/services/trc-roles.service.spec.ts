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

import { HttpTestingController } from "@angular/common/http/testing";
import {
  injectService,
  configureServiceTestingModule,
  cleanupAngularTest,
  createServiceInInjectionContext,
} from "../../test-helpers/angular-test-helpers";
import { TcrRolesService } from "./trc-roles.service";
import { Subject } from "rxjs";
import { TcrMessage, TcrMessageType } from "../interfaces/tcr-message";
import { TcrRole } from "../interfaces/tcr-role";
import { WebsocketService } from "./websocket.service";
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

describe("TcrRolesService", () => {
  let service: TcrRolesService;
  let httpMock: HttpTestingController;
  let wsServiceFake: WebsocketService;

  beforeEach(() => {
    configureServiceTestingModule(TcrRolesService, [
      { provide: WebsocketService, useClass: FakeWebsocketService },
      { provide: DestroyRef, useClass: MockDestroyRef },
    ]);

    // Create service using injection context helper to handle takeUntilDestroyed
    service = createServiceInInjectionContext<TcrRolesService>(
      TcrRolesService,
      [{ provide: WebsocketService, useClass: FakeWebsocketService }],
    );

    httpMock = injectService(HttpTestingController);
    wsServiceFake = injectService(WebsocketService);
  });

  afterEach(() => {
    cleanupAngularTest(httpMock);
  });

  describe("service instance", () => {
    it("should be created", () => {
      expect(service).toBeTruthy();
    });
  });

  describe("getRole() function", () => {
    it("should return role info when called", () => {
      const roleName = "some-role";
      const sample: TcrRole = {
        name: roleName,
        description: "some role description",
        active: false,
      };

      let actual: TcrRole | undefined;
      service.getRole(roleName).subscribe((other) => {
        actual = other;
      });

      const req = httpMock.expectOne(`/api/roles/${roleName}`);
      expect(req.request.method).toBe("GET");
      expect(req.request.responseType).toEqual("json");
      req.flush(sample);
      expect(actual).toBe(sample);
    });

    it("should return undefined when receiving an error response", () => {
      const roleName = "some-role";
      let actual: TcrRole | undefined;
      service.getRole(roleName).subscribe((other) => {
        actual = other;
      });

      const req = httpMock.expectOne(`/api/roles/${roleName}`);
      expect(req.request.method).toBe("GET");
      req.flush(
        { message: "Bad Request" },
        {
          status: 400,
          statusText: "Bad Request",
        },
      );
      expect(actual).toBeUndefined();
    });
  });

  describe("activateRole() function", () => {
    const testCases = [
      { state: true, action: "start" },
      { state: false, action: "stop" },
    ];

    testCases.forEach(({ state, action }) => {
      it(`should send ${action} role request when called with state ${state}`, () => {
        const roleName = "some-role";
        const sample: TcrRole = {
          name: roleName,
          description: "some role description",
          active: state,
        };

        let actual: TcrRole | undefined;
        service.activateRole(roleName, state).subscribe((other) => {
          actual = other;
        });

        const req = httpMock.expectOne(`/api/roles/${roleName}/${action}`);
        expect(req.request.method).toBe("POST");
        expect(req.request.responseType).toEqual("json");
        req.flush(sample);
        expect(actual).toBe(sample);
      });
    });

    it("should return undefined when receiving an error response", () => {
      const roleName = "some-role";
      let actual: TcrRole | undefined;
      service.activateRole(roleName, true).subscribe((other) => {
        actual = other;
      });

      const req = httpMock.expectOne(`/api/roles/${roleName}/start`);
      expect(req.request.method).toBe("POST");
      req.flush(
        { message: "Bad Request" },
        {
          status: 400,
          statusText: "Bad Request",
        },
      );
      expect(actual).toBeUndefined();
    });
  });

  describe("websocket message handler", () => {
    it("should forward role messages", async () => {
      const sampleMessage = { type: TcrMessageType.ROLE } as TcrMessage;

      const messagePromise = new Promise<TcrMessage>((resolve) => {
        service.message$.subscribe((msg) => {
          resolve(msg);
        });
      });

      wsServiceFake.webSocket$.next(sampleMessage);

      const receivedMessage = await messagePromise;
      expect(receivedMessage).toEqual(sampleMessage);
    });

    it("should drop non-role messages", async () => {
      const roleMessage = { type: TcrMessageType.ROLE } as TcrMessage;
      const infoMessage = { type: TcrMessageType.INFO } as TcrMessage;

      const messagePromise = new Promise<TcrMessage>((resolve) => {
        service.message$.subscribe((msg) => {
          // Should only receive role messages
          resolve(msg);
        });
      });

      // Send non-role message first (should be filtered out)
      wsServiceFake.webSocket$.next(infoMessage);
      // Send role message (should be received)
      wsServiceFake.webSocket$.next(roleMessage);

      const receivedMessage = await messagePromise;
      expect(receivedMessage.type).toEqual(TcrMessageType.ROLE);
    });
  });
});
