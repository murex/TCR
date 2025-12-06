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

import "zone.js";
import "zone.js/testing";
import { getTestBed } from "@angular/core/testing";
import {
  BrowserDynamicTestingModule,
  platformBrowserDynamicTesting,
} from "@angular/platform-browser-dynamic/testing";
import { FaIconLibrary } from "@fortawesome/angular-fontawesome";
import { registerFontAwesomeIcons } from "./app/shared/font-awesome-icons";
import { vi } from "vitest";
import {
  provideHttpClient,
  withInterceptorsFromDi,
} from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { setupXTermMocks } from "./test-helpers/xterm-mocks";

// Make Jasmine globally available for compatibility
declare global {
  var jasmine: unknown;
}

// Set up global Jasmine object for compatibility with existing tests
if (typeof globalThis.jasmine === "undefined") {
  globalThis.jasmine = {
    createSpy: (_name?: string) => {
      const spy = vi.fn();
      // Add jasmine-specific methods
      spy.and = {
        callFake: (fn: (...args: unknown[]) => unknown) => {
          spy.mockImplementation(fn);
          return spy;
        },
        returnValue: (value: unknown) => {
          spy.mockReturnValue(value);
          return spy;
        },
        callThrough: () => {
          spy.mockImplementation(undefined);
          return spy;
        },
        stub: () => {
          spy.mockImplementation(() => undefined);
          return spy;
        },
      };
      return spy;
    },
    createSpyObj: (baseName: string, methodNames: string[]) => {
      const obj: Record<string, unknown> = {};
      methodNames.forEach((method) => {
        const spy = vi.fn();
        spy.and = {
          callFake: (fn: (...args: unknown[]) => unknown) => {
            spy.mockImplementation(fn);
            return spy;
          },
          returnValue: (value: unknown) => {
            spy.mockReturnValue(value);
            return spy;
          },
          callThrough: () => {
            spy.mockImplementation(undefined);
            return spy;
          },
          stub: () => {
            spy.mockImplementation(() => undefined);
            return spy;
          },
        };
        obj[method] = spy;
      });
      return obj;
    },
    DEFAULT_TIMEOUT_INTERVAL: 10000,
  };
}

// Initialize the Angular testing environment only once
if (!getTestBed().platform) {
  getTestBed().initTestEnvironment(
    BrowserDynamicTestingModule,
    platformBrowserDynamicTesting(),
    {
      teardown: { destroyAfterEach: false },
    },
  );
}

// Ensure proper Angular compilation and DI setup
// Force immediate compilation of the initial test environment
if (getTestBed().platform) {
  try {
    getTestBed().compileComponents();
  } catch {
    // Ignore compilation errors during initial setup
  }
}

// Setup XTerm mocks for testing
setupXTermMocks();

// Mock Web APIs that might not be available in happy-dom
Object.defineProperty(window, "ResizeObserver", {
  writable: true,
  value: vi.fn().mockImplementation(() => ({
    observe: vi.fn(),
    unobserve: vi.fn(),
    disconnect: vi.fn(),
  })),
});

Object.defineProperty(window, "IntersectionObserver", {
  writable: true,
  value: vi.fn().mockImplementation((callback, options) => ({
    observe: vi.fn(),
    unobserve: vi.fn(),
    disconnect: vi.fn(),
    callback,
    options,
  })),
});

// Mock HTMLElement methods
Object.defineProperty(HTMLElement.prototype, "scrollIntoView", {
  writable: true,
  value: vi.fn(),
});

Object.defineProperty(HTMLElement.prototype, "focus", {
  writable: true,
  value: vi.fn(),
});

// Mock getBoundingClientRect with consistent values
const mockGetBoundingClientRect = vi.fn(() => ({
  x: 0,
  y: 0,
  width: 100,
  height: 100,
  top: 0,
  right: 100,
  bottom: 100,
  left: 0,
  toJSON: vi.fn(),
}));

Object.defineProperty(HTMLElement.prototype, "getBoundingClientRect", {
  writable: true,
  value: mockGetBoundingClientRect,
});

// Mock console.error to suppress expected test errors
const originalConsoleError = console.error;
console.error = vi.fn((...args: unknown[]) => {
  const message = args[0]?.toString() || "";

  // Suppress known HTTP error messages from services
  if (
    message.includes("getTimer -") ||
    message.includes("getRole -") ||
    message.includes("getBuildInfo -") ||
    message.includes("activateRole -") ||
    message.includes("abort-command -")
  ) {
    return; // Suppress these expected errors in tests
  }

  // Suppress terminal disposal warnings
  if (
    message.includes("Error disposing xterm:") ||
    message.includes("Error clearing terminal:") ||
    message.includes("Terminal cleanup warning:")
  ) {
    return; // Suppress these expected warnings in tests
  }

  // Call original console.error for other messages
  originalConsoleError.apply(console, args);
});

// Mock xterm instance tracking for cleanup
Object.defineProperty(window, "__xtermInstances", {
  writable: true,
  value: [],
});

// Global cleanup function for xterm instances
function cleanupXtermInstances() {
  if (window.__xtermInstances) {
    for (const instance of window.__xtermInstances) {
      try {
        if (instance?._core?.viewport?._refreshAnimationFrame) {
          cancelAnimationFrame(instance._core.viewport._refreshAnimationFrame);
          instance._core.viewport._refreshAnimationFrame = null;
        }
        if (instance?.dispose && typeof instance.dispose === "function") {
          instance.dispose();
        }
      } catch (_e) {
        // Ignore cleanup errors
      }
    }
    window.__xtermInstances = [];
  }
}

// Global setup for each test
beforeEach(() => {
  // Reset TestBed before each test to avoid conflicts
  try {
    getTestBed().resetTestingModule();
  } catch (_error) {
    // If TestBed isn't initialized yet, ignore the error
  }

  // Reset xterm instances
  window.__xtermInstances = [];

  // Initialize FontAwesome icons for all tests
  const iconLibrary = new FaIconLibrary();
  registerFontAwesomeIcons(iconLibrary);
});

// Global cleanup after each test
afterEach(() => {
  // Clean up xterm instances
  cleanupXtermInstances();

  // Clear all mocks
  vi.clearAllMocks();

  // Reset TestBed after each test
  try {
    getTestBed().resetTestingModule();
  } catch (_error) {
    // If TestBed is already reset, ignore the error
  }

  // Reset mock implementations to default
  mockGetBoundingClientRect.mockReturnValue({
    x: 0,
    y: 0,
    width: 100,
    height: 100,
    top: 0,
    right: 100,
    bottom: 100,
    left: 0,
    toJSON: vi.fn(),
  });
});

// Set default timeout for tests
vi.setConfig({
  testTimeout: 10000,
  hookTimeout: 10000,
});

// Global TestBed configuration helper with enhanced DI support
globalThis.configureTestBed = (config: {
  imports?: unknown[];
  providers?: unknown[];
  declarations?: unknown[];
}) => {
  try {
    getTestBed().resetTestingModule();
  } catch (_error) {
    // Ignore if already reset
  }

  const enhancedProviders = [
    provideHttpClient(withInterceptorsFromDi()),
    provideHttpClientTesting(),
    ...(config.providers || []),
  ];

  const testBedConfig = {
    imports: config.imports || [],
    providers: enhancedProviders,
    declarations: config.declarations || [],
  };

  const testBed = getTestBed().configureTestingModule(testBedConfig);

  // Ensure compilation for proper DI metadata
  testBed.compileComponents();

  return testBed;
};

// Enhanced service injection helper with DI error handling
globalThis.injectServiceWithFallback = <T>(serviceClass: unknown): T => {
  try {
    return getTestBed().inject(serviceClass);
  } catch (error: unknown) {
    if ((error as Error)?.message?.includes("NG0202")) {
      console.warn(
        `Applying DI fallback for ${(serviceClass as { name: string }).name}`,
      );

      // Reset and reconfigure with enhanced providers
      getTestBed().resetTestingModule();
      getTestBed().configureTestingModule({
        providers: [
          serviceClass,
          provideHttpClient(withInterceptorsFromDi()),
          provideHttpClientTesting(),
        ],
      });
      getTestBed().compileComponents();

      return getTestBed().inject(serviceClass);
    }
    throw error;
  }
};

// Extend expect with custom matchers if needed
declare global {
  interface Window {
    __xtermInstances: unknown[];
  }

  // Enhanced global test helpers
  function configureTestBed(config: {
    imports?: unknown[];
    providers?: unknown[];
    declarations?: unknown[];
  }): { compileComponents(): void };

  function injectServiceWithFallback<T>(serviceClass: {
    new (...args: unknown[]): T;
    name: string;
  }): T;
}
