import {TestBed} from '@angular/core/testing';

import {TcrBuildInfoService} from './tcr-build-info.service';

describe('TcrBuildInfoService', () => {
  let service: TcrBuildInfoService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TcrBuildInfoService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
