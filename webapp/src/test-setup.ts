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

// Global test configuration
declare const require: {
  context(
    path: string,
    deep?: boolean,
    filter?: RegExp,
  ): {
    <T>(id: string): T;
    keys(): string[];
  };
};

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
  (global as any).ResizeObserver = class ResizeObserver {
    observe() {}
    unobserve() {}
    disconnect() {}
  };
}

// Mock IntersectionObserver if not available
if (typeof IntersectionObserver === "undefined") {
  (global as any).IntersectionObserver = class IntersectionObserver {
    constructor(
      public callback: Function,
      public options?: any,
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

// Global afterEach to clean up any lingering timers or async operations
afterEach(() => {
  // Clear any pending timeouts (limited to reasonable range)
  const dummyTimeoutId = setTimeout(() => {}, 0) as any;
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
  const dummyIntervalId = setInterval(() => {}, 1000) as any;
  clearInterval(dummyIntervalId);
  // Clear a reasonable range of potential interval IDs
  for (
    let i = Math.max(0, Number(dummyIntervalId) - 100);
    i <= Number(dummyIntervalId);
    i++
  ) {
    clearInterval(i);
  }
});

// Set longer timeout for tests that might involve DOM operations
jasmine.DEFAULT_TIMEOUT_INTERVAL = 10000;
