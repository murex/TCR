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

import "zone.js/testing";
import { getTestBed } from "@angular/core/testing";
import {
  BrowserDynamicTestingModule,
  platformBrowserDynamicTesting,
} from "@angular/platform-browser-dynamic/testing";
import { FaIconLibrary } from "@fortawesome/angular-fontawesome";
import { registerFontAwesomeIcons } from "./app/shared/font-awesome-icons";

// Global test configuration

// Initialize the Angular testing environment
getTestBed().initTestEnvironment(
  BrowserDynamicTestingModule,
  platformBrowserDynamicTesting(),
);

// Global error handler for uncaught promise rejections
window.addEventListener("unhandledrejection", (event) => {
  console.warn("Unhandled promise rejection:", event.reason);
  // Prevent the default browser behavior
  event.preventDefault();
});

// Global error handler for uncaught errors
window.addEventListener("error", (event) => {
  console.warn("Uncaught error:", event.error);
});

// Mock ResizeObserver if not available (common in headless environments)
if (typeof ResizeObserver === "undefined") {
  (globalThis as { ResizeObserver?: unknown }).ResizeObserver =
    class ResizeObserver {
      observe() {}
      unobserve() {}
      disconnect() {}
    };
}

// Mock IntersectionObserver if not available
if (typeof IntersectionObserver === "undefined") {
  (globalThis as { IntersectionObserver?: unknown }).IntersectionObserver =
    class IntersectionObserver {
      constructor(
        public callback: (entries: unknown[], observer: unknown) => void,
        public options?: unknown,
      ) {}
      observe() {}
      unobserve() {}
      disconnect() {}
    };
}

// Mock HTMLElement methods that might not be available in headless Chrome
if (typeof HTMLElement !== "undefined") {
  // Mock scrollIntoView
  if (!HTMLElement.prototype.scrollIntoView) {
    HTMLElement.prototype.scrollIntoView = function () {};
  }

  // Mock focus
  if (!HTMLElement.prototype.focus) {
    HTMLElement.prototype.focus = function () {};
  }

  // Mock getBoundingClientRect for elements that might not have dimensions
  const originalGetBoundingClientRect =
    HTMLElement.prototype.getBoundingClientRect;
  HTMLElement.prototype.getBoundingClientRect = function () {
    const rect = originalGetBoundingClientRect.call(this);
    // Ensure dimensions are always available
    return {
      x: rect.x || 0,
      y: rect.y || 0,
      width: rect.width || 100,
      height: rect.height || 100,
      top: rect.top || 0,
      right: rect.right || 100,
      bottom: rect.bottom || 100,
      left: rect.left || 0,
      toJSON: rect.toJSON || (() => {}),
    };
  };
}

// Global cleanup for terminal/xterm instances
function cleanupXtermInstances() {
  // Clean up any xterm viewport refresh frames
  if (
    typeof window !== "undefined" &&
    (window as unknown as { __xtermInstances?: unknown[] }).__xtermInstances
  ) {
    const instances = (window as unknown as { __xtermInstances: unknown[] })
      .__xtermInstances;
    for (const instance of instances) {
      try {
        if (instance._core?.viewport?._refreshAnimationFrame) {
          cancelAnimationFrame(instance._core.viewport._refreshAnimationFrame);
          instance._core.viewport._refreshAnimationFrame = null;
        }
        if (instance.dispose && typeof instance.dispose === "function") {
          instance.dispose();
        }
      } catch (_e) {
        // Ignore cleanup errors
      }
    }
    (window as unknown as { __xtermInstances: unknown[] }).__xtermInstances =
      [];
  }
}

// Global beforeEach to prepare clean test environment
beforeEach(() => {
  // Initialize xterm instance tracking
  if (typeof window !== "undefined") {
    (window as unknown as { __xtermInstances: unknown[] }).__xtermInstances =
      [];
  }

  // Initialize FontAwesome icons for all tests
  const iconLibrary = new FaIconLibrary();
  registerFontAwesomeIcons(iconLibrary);
});

// Global afterEach to clean up any lingering timers or async operations
afterEach(() => {
  // Clean up any xterm instances first
  cleanupXtermInstances();

  // Cancel any pending animation frames
  if (typeof window !== "undefined" && window.cancelAnimationFrame) {
    // Cancel a reasonable range of animation frame IDs
    for (let i = 1; i <= 1000; i++) {
      try {
        window.cancelAnimationFrame(i);
      } catch (_e) {
        // Ignore errors for invalid IDs
      }
    }
  }

  // Clear any pending timeouts (limited to reasonable range)
  const dummyTimeoutId = setTimeout(() => {}, 0) as unknown as number;
  clearTimeout(dummyTimeoutId);
  // Clear a reasonable range of potential timeout IDs
  for (
    let i = Math.max(0, Number(dummyTimeoutId) - 100);
    i <= Number(dummyTimeoutId);
    i++
  ) {
    clearTimeout(i);
  }

  // Clear any pending intervals (limited to reasonable range)
  const dummyIntervalId = setInterval(() => {}, 1000) as unknown as number;
  clearInterval(dummyIntervalId);
  // Clear a reasonable range of potential interval IDs
  for (
    let i = Math.max(0, Number(dummyIntervalId) - 100);
    i <= Number(dummyIntervalId);
    i++
  ) {
    clearInterval(i);
  }

  // Clean up any pending promises
  if (typeof window !== "undefined" && window.Promise) {
    // Force a microtask flush
    Promise.resolve().then(() => {});
  }
});

// Set longer timeout for tests that might involve DOM operations
jasmine.DEFAULT_TIMEOUT_INTERVAL = 10000;

// Override console.error to suppress expected error messages in tests
const originalConsoleError = console.error;
console.error = (...args: unknown[]) => {
  // Suppress known HTTP error messages from services
  const message = args[0]?.toString() || "";
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
};
