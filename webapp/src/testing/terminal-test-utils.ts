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

import { Component } from "@angular/core";
import { NgTerminal } from "ng-terminal";
import { vi } from "vitest";

/**
 * Mock NgTerminal component for testing
 * This mock prevents xterm.js dimension errors in headless environments
 */
@Component({
  selector: "ng-terminal", // eslint-disable-line @angular-eslint/component-selector
  template: '<div class="mock-terminal"></div>',
  standalone: false, // eslint-disable-line @angular-eslint/prefer-standalone
})
export class MockNgTerminalComponent {
  underlying: MockTerminalUnderlying;
  private mockContent = "";

  constructor() {
    this.underlying = new MockTerminalUnderlying();
  }

  write(data: string): void {
    this.mockContent += data;
    // Simulate write callback if needed
  }

  setXtermOptions(_options: unknown): void {
    // Mock implementation
  }

  setRows(_rows: number): void {
    // Mock implementation
  }

  setCols(_cols: number): void {
    // Mock implementation
  }

  setMinWidth(_width: number): void {
    // Mock implementation
  }

  setMinHeight(_height: number): void {
    // Mock implementation
  }

  setDraggable(_draggable: boolean): void {
    // Mock implementation
  }

  getContent(): string {
    return this.mockContent;
  }

  clear(): void {
    this.mockContent = "";
  }
}

/**
 * Mock terminal underlying object that prevents dimension-related errors
 */
export class MockTerminalUnderlying {
  unicode = { activeVersion: "11" };
  private disposed = false;
  private timers: number[] = [];
  private _core: {
    viewport: {
      dimensions: { cols: number; rows: number } | undefined;
      syncScrollArea: ReturnType<typeof vi.fn>;
      _innerRefresh: ReturnType<typeof vi.fn>;
    };
  };

  constructor() {
    // Mock core viewport with safe dimensions
    this._core = {
      viewport: {
        dimensions: { cols: 80, rows: 24 },
        syncScrollArea: vi.fn(),
        _innerRefresh: vi.fn(),
      },
    };
  }

  reset = vi.fn();
  dispose = vi.fn().mockImplementation(() => {
    this.disposed = true;
    // Clear all timers
    this.timers.forEach((timer) => clearTimeout(timer));
    this.timers = [];
    // Clear dimensions to prevent access after disposal
    if (this._core && this._core.viewport) {
      this._core.viewport.dimensions = undefined;
    }
  });
  loadAddon = vi.fn();

  isDisposed(): boolean {
    return this.disposed;
  }

  // Track timers to clean them up
  _registerTimer(timerId: number): void {
    this.timers.push(timerId);
  }
}

/**
 * Creates a properly mocked NgTerminal for testing
 */
export function createMockNgTerminal(): NgTerminal {
  const mock = new MockNgTerminalComponent();
  return mock as unknown as NgTerminal;
}

/**
 * Test utility to safely clean up terminal components
 */
export function cleanupTerminal(terminal?: NgTerminal): void {
  if (terminal?.underlying) {
    try {
      // Stop any pending viewport sync operations
      const underlying = terminal.underlying as unknown as Record<
        string,
        unknown
      >;
      const core = underlying["_core"] as Record<string, unknown> | undefined;
      const viewport = core?.["viewport"] as
        | Record<string, unknown>
        | undefined;

      if (viewport) {
        // Clear any pending timers in viewport
        if (viewport["_refreshAnimationFrame"]) {
          cancelAnimationFrame(viewport["_refreshAnimationFrame"] as number);
          try {
            viewport["_refreshAnimationFrame"] = null;
          } catch {
            // Property might be read-only
          }
        }
        if (viewport["_coreBrowserService"]) {
          try {
            viewport["_coreBrowserService"] = null;
          } catch {
            // Property might be read-only
          }
        }
      }

      // Dispose the terminal
      if (typeof terminal.underlying.dispose === "function") {
        terminal.underlying.dispose();
      }

      // Clear the reference - use Object.defineProperty to bypass read-only
      try {
        Object.defineProperty(terminal, "underlying", {
          value: null,
          writable: true,
          configurable: true,
        });
      } catch {
        // If that fails, try direct assignment
        try {
          (terminal as unknown as Record<string, unknown>)["underlying"] = null;
        } catch {
          // Ignore if property is truly read-only
        }
      }
    } catch (error) {
      // Ignore cleanup errors in tests
      console.warn("Terminal cleanup warning:", error);
    }
  }
}

/**
 * Helper to create a test environment that's safe for terminal components
 */
export function setupTerminalTestEnvironment(): void {
  // Mock ResizeObserver if not available
  if (typeof ResizeObserver === "undefined") {
    (globalThis as { ResizeObserver?: unknown }).ResizeObserver =
      class ResizeObserver {
        observe() {}
        unobserve() {}
        disconnect() {}
      };
  }

  // Mock canvas context for xterm.js
  const mockCanvasContext = {
    fillRect: vi.fn(),
    clearRect: vi.fn(),
    getImageData: vi.fn().mockReturnValue({ data: new Uint8ClampedArray(4) }),
    putImageData: vi.fn(),
    createImageData: vi
      .fn()
      .mockReturnValue({ data: new Uint8ClampedArray(4) }),
    setTransform: vi.fn(),
    drawImage: vi.fn(),
    save: vi.fn(),
    fillText: vi.fn(),
    restore: vi.fn(),
    beginPath: vi.fn(),
    moveTo: vi.fn(),
    lineTo: vi.fn(),
    closePath: vi.fn(),
    stroke: vi.fn(),
    translate: vi.fn(),
    scale: vi.fn(),
    rotate: vi.fn(),
    arc: vi.fn(),
    fill: vi.fn(),
    measureText: vi.fn().mockReturnValue({ width: 10 }),
    transform: vi.fn(),
    rect: vi.fn(),
    clip: vi.fn(),
  };

  if (typeof HTMLCanvasElement !== "undefined") {
    HTMLCanvasElement.prototype.getContext = function () {
      return mockCanvasContext as unknown as CanvasRenderingContext2D;
    } as unknown as typeof HTMLCanvasElement.prototype.getContext;
  }

  // Mock getBoundingClientRect for consistent dimensions
  if (typeof Element !== "undefined") {
    Element.prototype.getBoundingClientRect = vi.fn().mockReturnValue({
      x: 0,
      y: 0,
      width: 800,
      height: 600,
      top: 0,
      right: 800,
      bottom: 600,
      left: 0,
      toJSON: () => {},
    }) as unknown as typeof Element.prototype.getBoundingClientRect;
  }
}

/**
 * Creates a test double for xterm Terminal that prevents dimension errors
 */
export function createSafeXtermMock() {
  return {
    loadAddon: vi.fn(),
    open: vi.fn(),
    write: vi.fn(),
    dispose: vi.fn(),
    reset: vi.fn(),
    focus: vi.fn(),
    blur: vi.fn(),
    select: vi.fn(),
    selectAll: vi.fn(),
    selectLines: vi.fn(),
    clearSelection: vi.fn(),
    getSelection: vi.fn().mockReturnValue(""),
    getSelectionPosition: vi.fn(),
    scrollLines: vi.fn(),
    scrollPages: vi.fn(),
    scrollToTop: vi.fn(),
    scrollToBottom: vi.fn(),
    scrollToLine: vi.fn(),
    clear: vi.fn(),
    unicode: { activeVersion: "11" },
    rows: 24,
    cols: 80,
    element: document.createElement("div"),
    textarea: document.createElement("textarea"),
    _core: {
      viewport: {
        dimensions: { cols: 80, rows: 24 },
        syncScrollArea: vi.fn(),
        _innerRefresh: vi.fn(),
        _refreshAnimationFrame: null,
        _coreBrowserService: null,
      },
      buffer: {
        active: {
          baseY: 0,
          cursorY: 0,
          cursorX: 0,
          length: 24,
        },
      },
    },
    onData: vi.fn(),
    onKey: vi.fn(),
    onLineFeed: vi.fn(),
    onScroll: vi.fn(),
    onSelectionChange: vi.fn(),
    onRender: vi.fn(),
    onResize: vi.fn(),
    onTitleChange: vi.fn(),
    onBell: vi.fn(),
  };
}
