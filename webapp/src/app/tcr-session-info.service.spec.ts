import { TestBed } from '@angular/core/testing';

import { TcrSessionInfoService } from './tcr-session-info.service';

describe('TcrSessionInfoService', () => {
  let service: TcrSessionInfoService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TcrSessionInfoService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
