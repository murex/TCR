import {TestBed} from '@angular/core/testing';

import {TcrTimerService} from './tcr-timer.service';

xdescribe('TcrTimerService', () => {
  let service: TcrTimerService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TcrTimerService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
