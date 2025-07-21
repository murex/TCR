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

/**
 * Mock NgTerminal component for testing
 * This mock prevents xterm.js dimension errors in headless environments
 */
@Component({
  selector: "app-ng-terminal", // eslint-disable-line @angular-eslint/component-selector
  template: '<div class="mock-terminal"></div>',
  standalone: true, // eslint-disable-line @angular-eslint/prefer-standalone
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
  private _core: {
    viewport: {
      dimensions: { cols: number; rows: number };
      syncScrollArea: jasmine.Spy;
      _innerRefresh: jasmine.Spy;
    };
  };

  constructor() {
    // Mock core viewport with safe dimensions
    this._core = {
      viewport: {
        dimensions: { cols: 80, rows: 24 },
        syncScrollArea: jasmine.createSpy("syncScrollArea"),
        _innerRefresh: jasmine.createSpy("_innerRefresh"),
      },
    };
  }

  reset = jasmine.createSpy("reset");
  dispose = jasmine.createSpy("dispose").and.callFake(() => {
    this.disposed = true;
  });
  loadAddon = jasmine.createSpy("loadAddon");

  isDisposed(): boolean {
    return this.disposed;
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
      if (typeof terminal.underlying.dispose === "function") {
        terminal.underlying.dispose();
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
    fillRect: jasmine.createSpy("fillRect"),
    clearRect: jasmine.createSpy("clearRect"),
    getImageData: jasmine
      .createSpy("getImageData")
      .and.returnValue({ data: new Uint8ClampedArray(4) }),
    putImageData: jasmine.createSpy("putImageData"),
    createImageData: jasmine
      .createSpy("createImageData")
      .and.returnValue({ data: new Uint8ClampedArray(4) }),
    setTransform: jasmine.createSpy("setTransform"),
    drawImage: jasmine.createSpy("drawImage"),
    save: jasmine.createSpy("save"),
    fillText: jasmine.createSpy("fillText"),
    restore: jasmine.createSpy("restore"),
    beginPath: jasmine.createSpy("beginPath"),
    moveTo: jasmine.createSpy("moveTo"),
    lineTo: jasmine.createSpy("lineTo"),
    closePath: jasmine.createSpy("closePath"),
    stroke: jasmine.createSpy("stroke"),
    translate: jasmine.createSpy("translate"),
    scale: jasmine.createSpy("scale"),
    rotate: jasmine.createSpy("rotate"),
    arc: jasmine.createSpy("arc"),
    fill: jasmine.createSpy("fill"),
    measureText: jasmine
      .createSpy("measureText")
      .and.returnValue({ width: 10 }),
    transform: jasmine.createSpy("transform"),
    rect: jasmine.createSpy("rect"),
    clip: jasmine.createSpy("clip"),
  };

  if (typeof HTMLCanvasElement !== "undefined") {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    HTMLCanvasElement.prototype.getContext = function () {
      return mockCanvasContext as any; // eslint-disable-line @typescript-eslint/no-explicit-any
    } as any; // eslint-disable-line @typescript-eslint/no-explicit-any
  }

  // Mock getBoundingClientRect for consistent dimensions
  if (typeof Element !== "undefined") {
    Element.prototype.getBoundingClientRect = jasmine
      .createSpy("getBoundingClientRect")
      .and.returnValue({
        x: 0,
        y: 0,
        width: 800,
        height: 600,
        top: 0,
        right: 800,
        bottom: 600,
        left: 0,
        toJSON: () => {},
      });
  }
}

/**
 * Creates a test double for xterm Terminal that prevents dimension errors
 */
export function createSafeXtermMock() {
  return {
    loadAddon: jasmine.createSpy("loadAddon"),
    open: jasmine.createSpy("open"),
    write: jasmine.createSpy("write"),
    dispose: jasmine.createSpy("dispose"),
    reset: jasmine.createSpy("reset"),
    focus: jasmine.createSpy("focus"),
    blur: jasmine.createSpy("blur"),
    select: jasmine.createSpy("select"),
    selectAll: jasmine.createSpy("selectAll"),
    selectLines: jasmine.createSpy("selectLines"),
    clearSelection: jasmine.createSpy("clearSelection"),
    getSelection: jasmine.createSpy("getSelection").and.returnValue(""),
    getSelectionPosition: jasmine.createSpy("getSelectionPosition"),
    scrollLines: jasmine.createSpy("scrollLines"),
    scrollPages: jasmine.createSpy("scrollPages"),
    scrollToTop: jasmine.createSpy("scrollToTop"),
    scrollToBottom: jasmine.createSpy("scrollToBottom"),
    scrollToLine: jasmine.createSpy("scrollToLine"),
    clear: jasmine.createSpy("clear"),
    unicode: { activeVersion: "11" },
    rows: 24,
    cols: 80,
    element: document.createElement("div"),
    textarea: document.createElement("textarea"),
    _core: {
      viewport: {
        dimensions: { cols: 80, rows: 24 },
        syncScrollArea: jasmine.createSpy("syncScrollArea"),
        _innerRefresh: jasmine.createSpy("_innerRefresh"),
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
    onData: jasmine.createSpy("onData"),
    onKey: jasmine.createSpy("onKey"),
    onLineFeed: jasmine.createSpy("onLineFeed"),
    onScroll: jasmine.createSpy("onScroll"),
    onSelectionChange: jasmine.createSpy("onSelectionChange"),
    onRender: jasmine.createSpy("onRender"),
    onResize: jasmine.createSpy("onResize"),
    onTitleChange: jasmine.createSpy("onTitleChange"),
    onBell: jasmine.createSpy("onBell"),
  };
}
