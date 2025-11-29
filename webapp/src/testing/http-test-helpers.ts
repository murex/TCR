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

import { Injectable, OnDestroy } from "@angular/core";
import { Observable, of, Subject } from "rxjs";
import { TcrTimer, TcrTimerState } from "../app/interfaces/tcr-timer";
import { TcrRole } from "../app/interfaces/tcr-role";
import { TcrBuildInfo } from "../app/interfaces/tcr-build-info";
import { TcrMessage } from "../app/interfaces/tcr-message";

/**
 * Mock TcrTimerService for testing
 */
@Injectable()
export class MockTcrTimerService {
  message$ = new Subject<TcrMessage>();

  getTimer(): Observable<TcrTimer> {
    return of({
      state: TcrTimerState.OFF,
      timeout: "0",
      elapsed: "0",
      remaining: "0",
    });
  }
}

/**
 * Mock TcrRolesService for testing
 */
@Injectable()
export class MockTcrRolesService {
  message$ = new Subject<TcrMessage>();
  private roleStates: Map<string, boolean> = new Map();

  getRole(name: string): Observable<TcrRole> {
    const active = this.roleStates.get(name) || false;
    return of({ name: name, description: `${name} role`, active: active });
  }

  activateRole(name: string, state: boolean): Observable<TcrRole> {
    this.roleStates.set(name, state);
    return of({ name: name, description: `${name} role`, active: state });
  }

  setRoleState(name: string, state: boolean): void {
    this.roleStates.set(name, state);
  }
}

/**
 * Mock TcrBuildInfoService for testing
 */
@Injectable()
export class MockTcrBuildInfoService {
  getBuildInfo(): Observable<TcrBuildInfo> {
    return of({
      version: "test-version",
      buildTime: new Date().toISOString(),
      os: "test-os",
      architecture: "test-arch",
    });
  }
}

/**
 * Mock TcrControlsService for testing
 */
@Injectable()
export class MockTcrControlsService {
  abortCommand(): Observable<unknown> {
    return of({});
  }
}

/**
 * Mock TcrMessageService for testing
 */
@Injectable()
export class MockTcrMessageService {
  message$ = new Subject<TcrMessage>();

  sendMessage(message: TcrMessage): void {
    this.message$.next(message);
  }

  clear(): void {
    // Mock implementation
  }
}

/**
 * Mock WebsocketService for testing
 */
@Injectable()
export class MockWebsocketService implements OnDestroy {
  webSocket$ = new Subject<TcrMessage>();

  ngOnDestroy(): void {
    this.webSocket$.complete();
  }
}

/**
 * Provider configuration for all mock HTTP services
 */
export const MOCK_HTTP_PROVIDERS = [
  { provide: "TcrTimerService", useClass: MockTcrTimerService },
  { provide: "TcrRolesService", useClass: MockTcrRolesService },
  { provide: "TcrBuildInfoService", useClass: MockTcrBuildInfoService },
  { provide: "TcrControlsService", useClass: MockTcrControlsService },
  { provide: "TcrMessageService", useClass: MockTcrMessageService },
  { provide: "WebsocketService", useClass: MockWebsocketService },
];

/**
 * Helper function to configure TestBed with mock HTTP services
 */
export function provideMockHttpServices() {
  return [
    MockTcrTimerService,
    MockTcrRolesService,
    MockTcrBuildInfoService,
    MockTcrControlsService,
    MockTcrMessageService,
    MockWebsocketService,
  ];
}

/**
 * Helper to create a mock HTTP error response
 */
export function createMockHttpError(status: number, statusText: string) {
  return {
    status,
    statusText,
    error: { message: `${status} ${statusText}` },
  };
}

/**
 * Helper function to flush all pending HTTP requests in tests
 */
export function flushPendingRequests(httpTestingController: {
  match: (fn: () => boolean) => unknown[];
}): void {
  const requests = httpTestingController.match(() => true);
  requests.forEach((req: { url: string; flush: (data: unknown) => void }) => {
    if (req.url.includes("/api/timer")) {
      req.flush({
        state: TcrTimerState.OFF,
        timeout: "0",
        elapsed: "0",
        remaining: "0",
      });
    } else if (req.url.includes("/api/roles/")) {
      const roleName = req.url.split("/").pop()?.split("?")[0];
      req.flush({
        name: roleName,
        description: `${roleName} role`,
        active: false,
      });
    } else if (req.url.includes("/api/build-info")) {
      req.flush({
        version: "test-version",
        buildTime: new Date().toISOString(),
        os: "test-os",
        architecture: "test-arch",
      });
    } else if (req.url.includes("/api/controls/")) {
      req.flush({});
    } else {
      // For any other API call, return empty success
      req.flush({});
    }
  });
}

/**
 * Mock HTTP Testing Module configuration
 */
export const mockHttpTestingModule = {
  imports: [],
  providers: provideMockHttpServices(),
};
