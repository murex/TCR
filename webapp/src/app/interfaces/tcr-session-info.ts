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

export interface TcrSessionInfo {
  baseDir: string;
  workDir: string;
  language: string;
  toolchain: string;
  vcsName: string;
  vcsSession: string;
  variant: string;
  gitAutoPush: boolean;
  messageSuffix: string;
}

export interface TcrVariant {
  description: string;
  statechartImageFile: string;
}

export const tcrVariants: { [key: string]: TcrVariant } = {
  "original": {
    description: "The Original",
    statechartImageFile: "variant-original.png",
  },
  "btcr": {
    description: "BTCR -- Build && Test && Commit || Revert",
    statechartImageFile: "variant-btcr.png",
  },
  "relaxed": {
    description: "The Relaxed",
    statechartImageFile: "variant-relaxed.png",
  },
  "introspective": {
    description: "The Introspective",
    statechartImageFile: "variant-introspective.png",
  },
};
