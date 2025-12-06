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
AUTHORS OR COPYRIGHT FOR LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

import { vi } from "vitest";

/**
 * Mock Terminal class for @xterm/xterm module
 */
export class Terminal {
  private _element: HTMLElement | null = null;
  private _rows = 24;
  private _cols = 80;

  constructor(options?: ITerminalOptions) {
    if (options?.rows) this._rows = options.rows;
    if (options?.cols) this._cols = options.cols;
  }

  // Core terminal methods
  open = vi.fn((element: HTMLElement) => {
    this._element = element;
  });

  write = vi.fn((_data: string) => {
    // Mock writing to terminal
  });

  writeln = vi.fn((_data: string) => {
    // Mock writing line to terminal
  });

  clear = vi.fn(() => {
    // Mock clearing terminal
  });

  reset = vi.fn(() => {
    // Mock resetting terminal
  });

  dispose = vi.fn(() => {
    // Mock disposing terminal
    this._element = null;
  });

  focus = vi.fn(() => {
    // Mock focusing terminal
  });

  blur = vi.fn(() => {
    // Mock blurring terminal
  });

  // Resize functionality
  resize = vi.fn((cols: number, rows: number) => {
    this._cols = cols;
    this._rows = rows;
  });

  // Event handling
  onData = vi.fn((_callback: (data: string) => void) => {
    return { dispose: vi.fn() };
  });

  onResize = vi.fn(
    (_callback: (size: { cols: number; rows: number }) => void) => {
      return { dispose: vi.fn() };
    },
  );

  onKey = vi.fn((_callback: (key: unknown) => void) => {
    return { dispose: vi.fn() };
  });

  // Selection
  getSelection = vi.fn(() => "");
  select = vi.fn();
  selectAll = vi.fn();
  clearSelection = vi.fn();

  // Addon support
  loadAddon = vi.fn((_addon: unknown) => {
    // Mock loading addons
  });

  // Properties
  get element(): HTMLElement | null {
    return this._element;
  }

  get rows(): number {
    return this._rows;
  }

  get cols(): number {
    return this._cols;
  }

  get buffer() {
    return {
      active: {
        cursorX: 0,
        cursorY: 0,
        baseY: 0,
        length: this._rows,
        getLine: vi.fn(() => ({
          translateToString: vi.fn(() => ""),
          isWrapped: false,
        })),
      },
      alternate: {
        cursorX: 0,
        cursorY: 0,
        baseY: 0,
        length: this._rows,
        getLine: vi.fn(() => ({
          translateToString: vi.fn(() => ""),
          isWrapped: false,
        })),
      },
    };
  }

  // Parser
  get parser() {
    return {
      registerCsiHandler: vi.fn(),
      registerDcsHandler: vi.fn(),
      registerEscHandler: vi.fn(),
      registerOscHandler: vi.fn(),
    };
  }

  // Unicode handling
  get unicode() {
    return {
      activeVersion: "11",
      versions: ["6", "11"],
    };
  }
}

// Export types for TypeScript support
export interface ITerminalOptions {
  rows?: number;
  cols?: number;
  theme?: unknown;
  fontSize?: number;
  fontFamily?: string;
  cursorBlink?: boolean;
  cursorStyle?: "block" | "underline" | "bar";
  scrollback?: number;
}

// Default export for CommonJS compatibility
export default { Terminal };
