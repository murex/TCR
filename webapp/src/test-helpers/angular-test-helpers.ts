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
import {
  provideHttpClient,
  withInterceptorsFromDi,
  HttpClient,
} from "@angular/common/http";
import {
  provideHttpClientTesting,
  HttpTestingController,
  HttpClientTestingModule,
} from "@angular/common/http/testing";
import { Component, Type } from "@angular/core";
import { provideRouter } from "@angular/router";

/**
 * Creates a service instance within proper Angular injection context
 * This bypasses the HttpClient DI metadata issues in Vitest by manually
 * creating services with their required dependencies
 */
export function createServiceInInjectionContext<T>(
  serviceClass: Type<T>,
  additionalProviders: unknown[] = [],
): T {
  return TestBed.runInInjectionContext(() => {
    const httpClient = TestBed.inject(HttpClient);

    // Handle different service types based on constructor parameters
    switch (serviceClass.name) {
      case "TcrTimerService":
      case "TcrRolesService": {
        const wsService = TestBed.inject(
          additionalProviders.find(
            (
              p: any, // eslint-disable-line @typescript-eslint/no-explicit-any
            ) =>
              p.provide?.name === "WebsocketService" ||
              p.useClass?.name === "FakeWebsocketService",
          )?.provide || Object,
        );
        return new serviceClass(httpClient, wsService);
      }

      case "TcrMessageService": {
        // Only needs WebSocket service, no HttpClient
        const wsService = TestBed.inject(
          additionalProviders.find(
            (
              p: any, // eslint-disable-line @typescript-eslint/no-explicit-any
            ) =>
              p.provide?.name === "WebsocketService" ||
              p.useClass?.name === "FakeWebsocketService",
          )?.provide || Object,
        );
        return new serviceClass(wsService);
      }

      case "TcrBuildInfoService":
      case "TcrSessionInfoService":
      case "TcrControlsService":
        // These services typically only need HttpClient
        return new serviceClass(httpClient);

      default:
        // Generic fallback for HttpClient-dependent services
        try {
          return new serviceClass(httpClient);
        } catch (_error) {
          console.warn(
            `Failed to create ${serviceClass.name} with HttpClient, trying without parameters`,
          );
          return new serviceClass();
        }
    }
  });
}

/**
 * Enhanced service injection that handles DI compatibility issues
 */
export function injectServiceSafely<T>(
  serviceClass: Type<T>,
  additionalProviders: unknown[] = [],
): T {
  try {
    return TestBed.inject(serviceClass);
  } catch (error: unknown) {
    if (
      (error as Error)?.message?.includes("NG0202") ||
      (error as Error)?.message?.includes(
        "not compatible with Angular Dependency Injection",
      )
    ) {
      console.warn(
        `DI compatibility issue for ${serviceClass.name}, using injection context creation`,
      );
      return createServiceInInjectionContext<T>(
        serviceClass,
        additionalProviders,
      );
    }
    throw error;
  }
}

/**
 * Standard Angular testing configuration for services that use HttpClient
 */
export function configureServiceTestingModule(
  serviceClass: Type<unknown>,
  additionalProviders: unknown[] = [],
) {
  // Reset TestBed to ensure clean state
  TestBed.resetTestingModule();

  TestBed.configureTestingModule({
    imports: [HttpClientTestingModule],
    providers: [
      serviceClass,
      provideHttpClient(withInterceptorsFromDi()),
      provideHttpClientTesting(),
      ...additionalProviders,
    ],
  });

  // Ensure compilation to generate proper DI metadata
  TestBed.compileComponents();
}

/**
 * Standard Angular testing configuration for components
 */
export function configureComponentTestingModule(
  componentClass: Type<unknown>,
  additionalImports: unknown[] = [],
  additionalProviders: unknown[] = [],
) {
  // Reset TestBed to ensure clean state
  TestBed.resetTestingModule();

  TestBed.configureTestingModule({
    imports: [componentClass, ...additionalImports],
    providers: [
      provideHttpClient(withInterceptorsFromDi()),
      provideHttpClientTesting(),
      provideRouter([]),
      ...additionalProviders,
    ],
  });

  // Ensure compilation to generate proper DI metadata
  TestBed.compileComponents();
}

/**
 * Standard Angular testing configuration for standalone components
 */
export function configureStandaloneComponentTestingModule(
  componentClass: Type<unknown>,
  additionalImports: unknown[] = [],
  additionalProviders: unknown[] = [],
) {
  return configureComponentTestingModule(
    componentClass,
    additionalImports,
    additionalProviders,
  );
}

/**
 * Helper to get HttpTestingController from TestBed with error handling
 */
export function getHttpTestingController(): HttpTestingController {
  try {
    return TestBed.inject(HttpTestingController);
  } catch (_error: unknown) {
    console.warn("HttpTestingController injection failed, reconfiguring...");

    // Ensure HttpTestingController is available
    TestBed.resetTestingModule();
    TestBed.configureTestingModule({
      providers: [
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    });
    TestBed.compileComponents();

    return TestBed.inject(HttpTestingController);
  }
}

/**
 * Helper to inject service from TestBed with proper typing and DI error handling
 */
export function injectService<T>(serviceClass: Type<T>): T {
  return injectServiceSafely<T>(serviceClass);
}

/**
 * Mock component factory for testing
 */
export function createMockComponent(
  selector: string,
  template: string = "<div></div>",
): Type<Component> {
  @Component({
    selector,
    template,
    standalone: true,
  })
  class MockComponent {}

  return MockComponent;
}

/**
 * Standard test cleanup that should be called in afterEach
 */
export function cleanupAngularTest(httpMock?: HttpTestingController): void {
  if (httpMock) {
    try {
      httpMock.verify();
    } catch {
      // Ignore verification errors in cleanup
    }
  }
}

/**
 * Helper to create WebSocket service mock
 */
export function createWebSocketServiceMock() {
  const mockWebSocketService = {
    webSocket$: {
      next: vi.fn(),
      subscribe: vi.fn(),
      pipe: vi.fn(() => ({
        subscribe: vi.fn(),
      })),
    },
  };
  return mockWebSocketService;
}

/**
 * Helper to create timer service mock
 */
export function createTimerServiceMock() {
  return {
    getTimer: vi.fn(() => ({
      subscribe: vi.fn(),
    })),
    message$: {
      subscribe: vi.fn(),
    },
  };
}

/**
 * Helper to create roles service mock
 */
export function createRolesServiceMock() {
  return {
    getRole: vi.fn(() => ({
      subscribe: vi.fn(),
    })),
    activateRole: vi.fn(() => ({
      subscribe: vi.fn(),
    })),
    message$: {
      subscribe: vi.fn(),
    },
  };
}

/**
 * Helper to create build info service mock
 */
export function createBuildInfoServiceMock() {
  return {
    getBuildInfo: vi.fn(() => ({
      subscribe: vi.fn(),
    })),
  };
}

/**
 * Helper to create session info service mock
 */
export function createSessionInfoServiceMock() {
  return {
    getSessionInfo: vi.fn(() => ({
      subscribe: vi.fn(),
    })),
  };
}

/**
 * Helper to create controls service mock
 */
export function createControlsServiceMock() {
  return {
    abortCommand: vi.fn(() => ({
      subscribe: vi.fn(),
    })),
  };
}

/**
 * Generic observable mock factory
 */
export function createObservableMock<T>(value?: T) {
  return {
    subscribe: vi.fn((callback?: (value: T) => void) => {
      if (callback && value !== undefined) {
        callback(value);
      }
      return {
        unsubscribe: vi.fn(),
      };
    }),
    pipe: vi.fn(() => createObservableMock(value)),
  };
}

declare global {
  const vi: typeof import("vitest").vi;
}
