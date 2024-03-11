import {TestBed} from '@angular/core/testing';

import {TcrRolesService} from './trc-roles.service';

xdescribe('TcrRolesService', () => {
  let service: TcrRolesService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TcrRolesService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
