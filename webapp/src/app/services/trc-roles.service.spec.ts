import {TestBed} from '@angular/core/testing';

import {TrcRolesService} from './trc-roles.service';

describe('TrcRolesService', () => {
  let service: TrcRolesService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TrcRolesService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
