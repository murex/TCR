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

/// <reference types="vitest" />

import { defineConfig } from "vite";
const { angularTransformer } = require("./vite-angular-transformer.js");

export default defineConfig({
  plugins: [angularTransformer()],
  resolve: {
    alias: {
      "@": "/src",
      "@xterm/xterm": "/src/test-helpers/xterm-terminal-mock.ts",
      "@xterm/addon-web-links": "/src/test-helpers/xterm-weblinks-mock.ts",
      "@xterm/addon-unicode11": "/src/test-helpers/xterm-unicode11-mock.ts",
    },
  },
  assetsInclude: ["**/*.html", "**/*.css"],
  optimizeDeps: {
    include: [
      "@angular/core",
      "@angular/common",
      "@angular/common/http",
      "@angular/platform-browser",
      "@angular/platform-browser-dynamic",
      "@angular/core/testing",
      "@angular/common/testing",
      "@angular/common/http/testing",
      "@xterm/xterm",
      "@xterm/addon-web-links",
      "@xterm/addon-unicode11",
    ],
  },
  test: {
    globals: true,
    environment: "jsdom",
    setupFiles: ["src/vitest-setup.ts"],
    isolate: false,
    include: ["src/**/*.spec.ts"],
    exclude: ["node_modules", "dist", "coverage", "**/*.config.*"],
    reporters: ["default", "junit"],
    outputFile: {
      junit: "_test_results/output.xml",
    },
    server: {
      deps: {
        inline: [
          "@xterm/xterm",
          "@xterm/addon-web-links",
          "@xterm/addon-unicode11",
        ],
      },
    },
    coverage: {
      provider: "v8",
      reporter: ["html", "text-summary", "json"],
      reportsDirectory: "coverage",
      exclude: [
        "node_modules/",
        "src/vitest-setup.ts",
        "src/**/*.spec.ts",
        "src/**/*.d.ts",
        "coverage/**",
        "dist/**",
        "**/*.config.*",
        "**/karma.conf.js",
      ],
      thresholds: {
        global: {
          statements: 90,
          branches: 90,
          functions: 90,
          lines: 90,
        },
      },
    },
    pool: "forks",
    poolOptions: {
      forks: {
        singleFork: true,
      },
    },
    testTimeout: 120000,
    hookTimeout: 120000,
    teardownTimeout: 120000,
    // Define global types for better TypeScript support
    typecheck: {
      tsconfig: "tsconfig.spec.json",
    },
  },
  define: {
    "import.meta.vitest": undefined,
  },
  server: {
    fs: {
      allow: [".."],
    },
  },
  ssr: {
    noExternal: [
      "@xterm/xterm",
      "@xterm/addon-web-links",
      "@xterm/addon-unicode11",
    ],
  },
  esbuild: {
    target: "es2020",
  },
});
