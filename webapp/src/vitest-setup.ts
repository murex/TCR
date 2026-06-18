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

import { afterEach, beforeEach, vi } from "vitest";

// COMPREHENSIVE XTerm module mocks - must be hoisted at top level
vi.mock("@xterm/xterm", async () => {
  const mockTerminal = {
    loadAddon: vi.fn(),
    open: vi.fn(),
    write: vi.fn(),
    writeln: vi.fn(),
    clear: vi.fn(),
    reset: vi.fn(),
    dispose: vi.fn(),
    focus: vi.fn(),
    blur: vi.fn(),
    resize: vi.fn(),
    onData: vi.fn(() => ({ dispose: vi.fn() })),
    onResize: vi.fn(() => ({ dispose: vi.fn() })),
    onKey: vi.fn(() => ({ dispose: vi.fn() })),
    getSelection: vi.fn(() => ""),
    select: vi.fn(),
    selectAll: vi.fn(),
    clearSelection: vi.fn(),
    element: null,
    rows: 24,
    cols: 80,
    buffer: {
      active: {
        cursorX: 0,
        cursorY: 0,
        baseY: 0,
        length: 24,
        getLine: vi.fn(() => ({
          translateToString: vi.fn(() => ""),
          isWrapped: false,
        })),
      },
      alternate: {
        cursorX: 0,
        cursorY: 0,
        baseY: 0,
        length: 24,
        getLine: vi.fn(() => ({
          translateToString: vi.fn(() => ""),
          isWrapped: false,
        })),
      },
    },
    parser: {
      registerCsiHandler: vi.fn(),
      registerDcsHandler: vi.fn(),
      registerEscHandler: vi.fn(),
      registerOscHandler: vi.fn(),
    },
    unicode: {
      activeVersion: "11",
      versions: ["6", "11"],
    },
  };

  return {
    default: { Terminal: vi.fn(() => mockTerminal) },
    Terminal: vi.fn(() => mockTerminal),
    __esModule: true,
  };
});

vi.mock("@xterm/addon-web-links", async () => {
  const mockAddon = {
    activate: vi.fn(),
    dispose: vi.fn(),
  };

  return {
    default: { WebLinksAddon: vi.fn(() => mockAddon) },
    WebLinksAddon: vi.fn(() => mockAddon),
    __esModule: true,
  };
});

vi.mock("@xterm/addon-unicode11", async () => {
  const mockAddon = {
    activate: vi.fn(),
    dispose: vi.fn(),
  };

  return {
    default: { Unicode11Addon: vi.fn(() => mockAddon) },
    Unicode11Addon: vi.fn(() => mockAddon),
    __esModule: true,
  };
});

vi.mock("ng-terminal", async () => {
  const mockNgTerminal = {
    underlying: {
      reset: vi.fn(),
      dispose: vi.fn(),
      loadAddon: vi.fn(),
      unicode: { activeVersion: "11" },
    },
    write: vi.fn(),
    clear: vi.fn(),
    setXtermOptions: vi.fn(),
    setRows: vi.fn(),
    setCols: vi.fn(),
    setMinWidth: vi.fn(),
    setMinHeight: vi.fn(),
    setDraggable: vi.fn(),
  };

  return {
    default: {
      NgTerminal: vi.fn(() => mockNgTerminal),
      NgTerminalModule: vi.fn(),
    },
    NgTerminal: vi.fn(() => mockNgTerminal),
    NgTerminalModule: vi.fn(),
    __esModule: true,
  };
});

import "zone.js";
import "zone.js/testing";

import { getTestBed } from "@angular/core/testing";
import {
  BrowserTestingModule,
  platformBrowserTesting,
} from "@angular/platform-browser/testing";
import { FaIconLibrary } from "@fortawesome/angular-fontawesome";
import { registerFontAwesomeIcons } from "./app/shared/font-awesome-icons";

// Initialize the Angular testing environment only once
if (!getTestBed().platform) {
  getTestBed().initTestEnvironment(
    BrowserTestingModule,
    platformBrowserTesting(),
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

// XTerm mocks are now set up at the top of this file

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

// Global setup for each test
beforeEach(() => {
  // Reset TestBed before each test to avoid conflicts
  try {
    getTestBed().resetTestingModule();
  } catch (_error) {
    // If TestBed isn't initialized yet, ignore the error
  }

  // Initialize FontAwesome icons for all tests
  const iconLibrary = new FaIconLibrary();
  registerFontAwesomeIcons(iconLibrary);
});

// Global cleanup after each test
afterEach(() => {
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

