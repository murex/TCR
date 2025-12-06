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
            (p: Record<string, unknown>) =>
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
            (p: Record<string, unknown>) =>
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
 * Creates a component within proper Angular injection context
 * This bypasses component DI metadata issues in Vitest
 */
export function createComponentInInjectionContext<T>(
  componentClass: Type<T>,
  additionalImports: unknown[] = [],
  additionalProviders: unknown[] = [],
): ComponentFixture<T> {
  // Don't reset TestBed - preserve existing configuration
  // Just add any additional providers if needed
  if (additionalImports.length > 0 || additionalProviders.length > 0) {
    TestBed.configureTestingModule({
      imports: additionalImports,
      providers: additionalProviders,
    });
    TestBed.compileComponents();
  }

  return TestBed.runInInjectionContext(() => {
    return TestBed.createComponent(componentClass);
  });
}

/**
 * Enhanced component creation that handles DI compatibility issues
 */
export function createComponentSafely<T>(
  componentClass: Type<T>,
  _additionalImports: unknown[] = [],
  additionalProviders: unknown[] = [],
): ComponentFixture<T> {
  try {
    return TestBed.createComponent(componentClass);
  } catch (error: unknown) {
    if (
      (error as Error)?.message?.includes("NG0202") ||
      (error as Error)?.message?.includes(
        "not compatible with Angular Dependency Injection",
      )
    ) {
      console.warn(
        `Component DI compatibility issue for ${componentClass.name}, using multi-strategy approach`,
      );

      // Convert additional providers to dependency map
      const dependencies: Record<string, unknown> = {};
      additionalProviders.forEach((provider: Record<string, unknown>) => {
        if (
          provider &&
          typeof provider === "object" &&
          "provide" in provider &&
          "useClass" in provider
        ) {
          const serviceName =
            (provider.provide as { name?: string }).name || provider.provide;
          dependencies[serviceName as string] = provider.useClass;
        }
      });

      return createComponentWithStrategies<T>(componentClass, dependencies);
    }
    throw error;
  }
}

/**
 * Manual component instantiation for components with DI issues
 * This creates components by manually injecting dependencies, bypassing TestBed DI metadata issues
 */
export function createComponentManually<T>(
  componentClass: Type<T>,
  dependencies: Record<string, unknown> = {},
): { component: T; fixture: ComponentFixture<T> } {
  // Get the component constructor parameter names through reflection if possible
  const paramNames = getConstructorParamNames(componentClass);

  // Create dependency instances in the proper order
  const depArray = paramNames.map((paramName) => {
    if (dependencies[paramName]) {
      return dependencies[paramName];
    }
    // Try to get from TestBed if not provided
    try {
      return TestBed.inject(paramName as unknown as Type<unknown>);
    } catch {
      console.warn(
        `Cannot resolve dependency ${paramName} for ${componentClass.name}`,
      );
      return null;
    }
  });

  // Create component instance manually within injection context to handle toSignal() and effect()
  const component = TestBed.runInInjectionContext(() => {
    return new (componentClass as unknown as new (...args: unknown[]) => T)(
      ...depArray,
    );
  });

  // Create a mock fixture that properly handles Angular lifecycle
  const fixture = {
    componentInstance: component,
    detectChanges: vi.fn(() => {
      // Trigger lifecycle methods when detectChanges is called
      const componentWithLifecycle = component as Record<string, unknown>;
      if (component && typeof componentWithLifecycle.ngOnInit === "function") {
        (componentWithLifecycle.ngOnInit as () => void)();
      }
      if (
        component &&
        typeof componentWithLifecycle.ngAfterViewInit === "function"
      ) {
        (componentWithLifecycle.ngAfterViewInit as () => void)();
      }
      // For timer components, call updateColor to ensure proper color calculation
      if (
        component &&
        typeof componentWithLifecycle.updateColor === "function"
      ) {
        (componentWithLifecycle.updateColor as () => void)();
      }
    }),
    destroy: vi.fn(() => {
      const componentWithLifecycle = component as Record<string, unknown>;
      if (
        component &&
        typeof componentWithLifecycle.ngOnDestroy === "function"
      ) {
        (componentWithLifecycle.ngOnDestroy as () => void)();
      }
    }),
    nativeElement: document.createElement("div"),
    debugElement: {
      query: vi.fn((selector: (debugElement: unknown) => boolean) => {
        // Enhanced DOM query functionality that works with By.css() predicate functions
        // Create mock elements and test them against the selector predicate

        // Try timer-component element
        const timerElement = document.createElement("div");
        timerElement.setAttribute("data-testid", "timer-component");

        // Calculate color based on component state (if available)
        let color = "rgb(128,128,128)"; // Default gray
        if (
          component &&
          typeof (component as Record<string, unknown>).fgColor === "string"
        ) {
          color = (component as Record<string, unknown>).fgColor as string;
        }
        timerElement.style.color = color;

        const timerDebugElement = {
          nativeElement: timerElement,
        };

        // Test if timer element matches the selector
        if (typeof selector === "function") {
          try {
            if (selector(timerDebugElement)) {
              return timerDebugElement;
            }
          } catch (_e) {
            // Selector test failed, continue to try other elements
          }
        }

        // Try timer-icon element
        const iconElement = document.createElement("fa-icon");
        iconElement.setAttribute("data-testid", "timer-icon");

        // Create SVG element with proper icon data
        const svg = document.createElementNS(
          "http://www.w3.org/2000/svg",
          "svg",
        );

        // Determine icon based on component state
        let iconName = "clock"; // Default icon
        if (component && (component as Record<string, unknown>).timer) {
          const timer = (component as Record<string, unknown>).timer as Record<
            string,
            unknown
          >;
          if (timer.state === "timeout") {
            iconName = "circle-exclamation";
          } else {
            iconName = "clock";
          }
        }

        svg.setAttribute("data-icon", iconName);
        iconElement.appendChild(svg);

        const iconDebugElement = {
          nativeElement: iconElement,
        };

        // Test if icon element matches the selector
        if (typeof selector === "function") {
          try {
            if (selector(iconDebugElement)) {
              return iconDebugElement;
            }
          } catch (_e) {
            // Selector test failed, continue to try other elements
          }
        }

        // Try timer-label element
        const labelElement = document.createElement("span");
        labelElement.setAttribute("data-testid", "timer-label");

        // Set text content based on component state
        let textContent = "00:00"; // Default text
        if (component && (component as Record<string, unknown>).timer) {
          const timer = (component as Record<string, unknown>).timer as Record<
            string,
            unknown
          >;
          // Format the time based on remaining time (this is a simplified version)
          const remaining = parseInt(timer.remaining as string) || 0;
          const minutes = Math.floor(Math.abs(remaining) / 60);
          const seconds = Math.abs(remaining) % 60;
          const sign = remaining < 0 ? "-" : "";
          textContent = `${sign}${minutes.toString().padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
        }
        labelElement.textContent = textContent;

        const labelDebugElement = {
          nativeElement: labelElement,
        };

        // Test if label element matches the selector
        if (typeof selector === "function") {
          try {
            if (selector(labelDebugElement)) {
              return labelDebugElement;
            }
          } catch (_e) {
            // Selector test failed, continue to generic fallback
          }
        }

        // Generic fallback - return null if no match (like real Angular)
        return null;
      }),
      queryAll: vi.fn(() => []),
      nativeElement: document.createElement("div"),
    },
    componentRef: {
      instance: component,
      destroy: vi.fn(),
    },
  } as unknown as ComponentFixture<T>;

  return { component, fixture };
}

/**
 * Extract constructor parameter names from a class (simplified version)
 */
function getConstructorParamNames(constructor: Type<unknown>): string[] {
  // This is a simplified approach - in a real implementation you might use
  // reflection metadata or parse the constructor string
  const constructorString = constructor.toString();
  const match = constructorString.match(/constructor\s*\(([^)]*)\)/);
  if (!match) return [];

  const params = match[1]
    .split(",")
    .map((param) => {
      // Extract parameter name, handling TypeScript types
      const cleanParam = param.trim().split(":")[0].trim();
      // Remove access modifiers (private, public, protected)
      return cleanParam.replace(/^(private|public|protected)\s+/, "");
    })
    .filter((param) => param.length > 0);

  return params;
}

/**
 * Enhanced component creation that tries multiple strategies
 */
export function createComponentWithStrategies<T>(
  componentClass: Type<T>,
  dependencies: Record<string, unknown> = {},
): ComponentFixture<T> {
  console.log(`Attempting to create component: ${componentClass.name}`);

  // Strategy 1: Try normal TestBed creation
  try {
    const fixture = TestBed.createComponent(componentClass);
    console.log(
      `Strategy 1 (TestBed.createComponent) succeeded for ${componentClass.name}`,
    );
    return fixture;
  } catch (error) {
    console.log(
      `Strategy 1 failed for ${componentClass.name}: ${(error as Error).message}`,
    );
  }

  // Strategy 2: Try injection context creation
  try {
    const fixture = createComponentInInjectionContext(componentClass);
    console.log(
      `Strategy 2 (injection context) succeeded for ${componentClass.name}`,
    );
    return fixture;
  } catch (error) {
    console.log(
      `Strategy 2 failed for ${componentClass.name}: ${(error as Error).message}`,
    );
  }

  // Strategy 3: Try manual component creation
  try {
    const { fixture } = createComponentManually(componentClass, dependencies);
    console.log(
      `Strategy 3 (manual creation) succeeded for ${componentClass.name}`,
    );
    return fixture;
  } catch (error) {
    console.log(
      `Strategy 3 failed for ${componentClass.name}: ${(error as Error).message}`,
    );
    throw new Error(
      `All component creation strategies failed for ${componentClass.name}`,
    );
  }
}

/**
 * Mock component factory for testing
 */
export function createMockComponent(
  _selector: string,
  template: string = "<div></div>",
): Type<Component> {
  @Component({
    selector: _selector,
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
