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

import {TestBed} from '@angular/core/testing';

import {TcrRolesService} from './trc-roles.service';
import {Subject} from "rxjs";
import {TcrMessage, TcrMessageType} from "../interfaces/tcr-message";
import {
  HttpTestingController,
  provideHttpClientTesting
} from "@angular/common/http/testing";
import {WebsocketService} from "./websocket.service";
import {TcrRole} from "../interfaces/tcr-role";
import {provideHttpClient, withInterceptorsFromDi} from '@angular/common/http';

class FakeWebsocketService {
  webSocket$: Subject<TcrMessage> = new Subject<TcrMessage>();
}

describe('TcrRolesService', () => {
  let service: TcrRolesService;
  let httpMock: HttpTestingController;
  let wsServiceFake: WebsocketService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [],
      providers: [
        TcrRolesService,
        {provide: WebsocketService, useClass: FakeWebsocketService},
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    });

    service = TestBed.inject(TcrRolesService);
    httpMock = TestBed.inject(HttpTestingController);
    wsServiceFake = TestBed.inject(WebsocketService);
  });

  afterEach(() => {
    httpMock.verify();
  });

  describe('service instance', () => {
    it('should be created', () => {
      expect(service).toBeTruthy();
    });
  });

  describe('getRole() function', () => {
    it('should return role info when called', () => {
      const roleName = "some-role";
      const sample: TcrRole = {
        name: roleName,
        description: "some role description",
        active: false,
      };

      let actual: TcrRole | undefined;
      service.getRole(roleName).subscribe(other => {
        actual = other;
      });

      const req = httpMock.expectOne(`/api/roles/${roleName}`);
      expect(req.request.method).toBe('GET');
      expect(req.request.responseType).toEqual('json');
      req.flush(sample);
      expect(actual).toBe(sample);
    });

    it('should return undefined when receiving an error response', () => {
      const roleName = "some-role";
      let actual: TcrRole | undefined;
      service.getRole(roleName).subscribe(other => {
        actual = other;
      });

      const req = httpMock.expectOne(`/api/roles/${roleName}`);
      expect(req.request.method).toBe('GET');
      req.flush({message: 'Bad Request'}, {
        status: 400,
        statusText: 'Bad Request'
      });
      expect(actual).toBeUndefined();
    });

  });

  describe('activateRole() function', () => {
    const testCases = [
      {state: true, action: 'start'},
      {state: false, action: 'stop'}
    ];

    testCases.forEach(({state, action}) => {
      it(`should send ${action} role request when called with state ${state}`, () => {
        const roleName = "some-role";
        const sample: TcrRole = {
          name: roleName,
          description: "some role description",
          active: state,
        };

        let actual: TcrRole | undefined;
        service.activateRole(roleName, state).subscribe(other => {
          actual = other;
        });

        const req = httpMock.expectOne(`/api/roles/${roleName}/${action}`);
        expect(req.request.method).toBe('POST');
        expect(req.request.responseType).toEqual('json');
        req.flush(sample);
        expect(actual).toBe(sample);
      });
    });

    it('should return undefined when receiving an error response', () => {
      const roleName = "some-role";
      let actual: TcrRole | undefined;
      service.activateRole(roleName, true).subscribe(other => {
        actual = other;
      });

      const req = httpMock.expectOne(`/api/roles/${roleName}/start`);
      expect(req.request.method).toBe('POST');
      req.flush({message: 'Bad Request'}, {
        status: 400,
        statusText: 'Bad Request'
      });
      expect(actual).toBeUndefined();
    });

  });

  describe('websocket message handler', () => {
    it('should forward role messages', (done) => {
      const sampleMessage = {type: TcrMessageType.ROLE} as TcrMessage;
      let actual: TcrMessage | undefined;
      service.message$.subscribe((msg) => {
        actual = msg;
        done();
      });
      wsServiceFake.webSocket$.next(sampleMessage);
      expect(actual).toEqual(sampleMessage);
    });

    it('should drop non-role messages', (done) => {
      const sampleMessage = {type: TcrMessageType.INFO} as TcrMessage;
      let actual: TcrMessage | undefined;
      service.message$.subscribe((msg) => {
        actual = msg;
        done();
      });
      wsServiceFake.webSocket$.next(sampleMessage);
      // Wait for the message to be processed by the service before checking the result
      setTimeout(() => done(), 10);
      expect(actual).toBeUndefined();
    });
  });
});
