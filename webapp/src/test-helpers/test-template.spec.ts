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

/**
 * Test Template - Proper Setup Pattern for Vitest + Angular
 *
 * This file demonstrates the correct patterns for setting up tests
 * in the migrated Vitest environment. Use these patterns to fix failing tests.
 */

import { TestBed } from "@angular/core/testing";
import { HttpTestingController } from "@angular/common/http/testing";
import { ComponentFixture } from "@angular/core/testing";
import {
  configureServiceTestingModule,
  configureComponentTestingModule,
  getHttpTestingController,
  injectService,
  cleanupAngularTest,
  createWebSocketServiceMock,
} from "./angular-test-helpers";

// Example service test pattern
describe("Service Test Template", () => {
  let service: MockService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    // Use helper function to configure TestBed
    configureServiceTestingModule(MockService, [
      // Additional providers if needed
    ]);

    // Inject dependencies
    service = injectService(MockService);
    httpMock = getHttpTestingController();
  });

  afterEach(() => {
    // Clean up HTTP mock
    cleanupAngularTest(httpMock);
  });

  it("should be created", () => {
    expect(service).toBeTruthy();
  });

  it("should make HTTP request", () => {
    const mockData = { test: "data" };

    service.getData().subscribe((data: { test: string }) => {
      expect(data).toEqual(mockData);
    });

    const req = httpMock.expectOne("/api/data");
    expect(req.request.method).toBe("GET");
    req.flush(mockData);
  });
});

// Example component test pattern
describe("Component Test Template", () => {
  let component: MockComponent;
  let fixture: ComponentFixture<MockComponent>;

  beforeEach(async () => {
    // Use helper function to configure TestBed
    configureComponentTestingModule(
      MockComponent,
      [
        // Additional imports
      ],
      [
        // Additional providers
      ],
    );

    await TestBed.compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MockComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  afterEach(() => {
    cleanupAngularTest();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should render content", () => {
    const compiled = fixture.nativeElement;
    expect(compiled.querySelector("h1")).toBeTruthy();
  });
});

// Example test with WebSocket service mock
describe("Service with WebSocket Template", () => {
  let service: MockServiceWithWebSocket;
  let mockWebSocketService: ReturnType<typeof createWebSocketServiceMock>;

  beforeEach(() => {
    mockWebSocketService = createWebSocketServiceMock();

    configureServiceTestingModule(MockServiceWithWebSocket, [
      { provide: "WebSocketService", useValue: mockWebSocketService },
    ]);

    service = injectService(MockServiceWithWebSocket);
  });

  afterEach(() => {
    cleanupAngularTest();
  });

  it("should handle websocket messages", () => {
    const testMessage = { type: "test", data: "value" };

    service.message$.subscribe((message: { type: string; data: string }) => {
      expect(message).toEqual(testMessage);
    });

    // Simulate WebSocket message
    mockWebSocketService.webSocket$.next(testMessage);
  });
});

// Mock classes for examples
class MockService {
  constructor() {}
  getData() {
    return { subscribe: () => {} };
  }
}

class MockComponent {
  constructor() {}
}

class MockServiceWithWebSocket {
  message$ = { subscribe: () => {} };
  constructor() {}
}

/**
 * MIGRATION CHECKLIST FOR FAILING TESTS:
 *
 * 1. Replace TestBed.configureTestingModule with helper functions:
 *    - configureServiceTestingModule() for services
 *    - configureComponentTestingModule() for components
 *
 * 2. Use helper functions for common operations:
 *    - injectService() instead of TestBed.inject()
 *    - getHttpTestingController() for HTTP mocking
 *    - cleanupAngularTest() in afterEach()
 *
 * 3. Ensure proper beforeEach/afterEach structure:
 *    - Configure TestBed in beforeEach
 *    - Clean up in afterEach
 *    - Use async/await for component compilation
 *
 * 4. For WebSocket services:
 *    - Use createWebSocketServiceMock()
 *    - Provide mocks via TestBed providers array
 *
 * 5. Handle component dependencies:
 *    - Mock all injected services
 *    - Provide all required imports
 *    - Use proper component fixture pattern
 *
 * 6. Common fixes for specific errors:
 *    - NG0202: Ensure all constructor dependencies are provided
 *    - TestBed already instantiated: Use proper beforeEach/afterEach
 *    - HTTP mock errors: Use cleanupAngularTest(httpMock) in afterEach
 */
