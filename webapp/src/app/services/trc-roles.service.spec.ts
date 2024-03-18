import {TestBed} from '@angular/core/testing';

import {TcrRolesService} from './trc-roles.service';
import {Subject} from "rxjs";
import {TcrMessage} from "../interfaces/tcr-message";
import {HttpClientTestingModule, HttpTestingController} from "@angular/common/http/testing";
import {WebsocketService} from "./websocket.service";
import {TcrRole} from "../interfaces/tcr-role";

class WebsocketServiceFake {
  webSocket$: Subject<TcrMessage> = new Subject<TcrMessage>();
}

describe('TcrRolesService', () => {
  let service: TcrRolesService;
  let httpMock: HttpTestingController;
  let wsServiceFake: WebsocketService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [
        TcrRolesService,
        {provide: WebsocketService, useClass: WebsocketServiceFake},
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
      req.flush({message: 'Bad Request'}, {status: 400, statusText: 'Bad Request'});
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
      req.flush({message: 'Bad Request'}, {status: 400, statusText: 'Bad Request'});
      expect(actual).toBeUndefined();
    });

  });

  describe('websocket message handler', () => {

    it('should forward role messages', (done) => {
      const sampleMessage: TcrMessage = {
        type: "role",
        emphasis: false,
        severity: "",
        text: "",
        timestamp: "",
      };
      let actual: TcrMessage | undefined;
      service.webSocket$.subscribe((msg) => {
        actual = msg;
        done();
      });
      wsServiceFake.webSocket$.next(sampleMessage);
      expect(actual).toEqual(sampleMessage);
    });

    it('should drop non-role messages', (done) => {
      const sampleMessage: TcrMessage = {
        type: "other",
        emphasis: false,
        severity: "",
        text: "",
        timestamp: "",
      };
      let actual: TcrMessage | undefined;
      service.webSocket$.subscribe((msg) => {
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
